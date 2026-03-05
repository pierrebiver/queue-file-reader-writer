# queue-file-reader-writer

## How to run it
1. Runs the Queue service with `go run ./queue-server`
2. Run the worker with `go run ./worker`

It should write the output in `files/output.txt`
## How to run the test
`go test ./...`

## Documentation of design choices and alternatives

## Queue service
The queue service is a standalone TCP server that holds an in-memory FIFO queue. It is intentionally dumb — it knows nothing about files or workers, only how to store and return lines. It runs as its own process so the producer and consumer are fully decoupled and only share state through the network, as they would in a real distributed system.
## Worker
The worker runs as a single process with two goroutines running concurrently. The producer goroutine reads the input file line by line and pushes each line to the queue service over TCP. The consumer goroutine independently pops lines from the queue service and writes them to the output file. They share no state directly — the queue service is the only point of coordination between them.

### Commands
Each protocol command is implemented as a struct satisfying a common Command interface. This means the server has no knowledge of how any individual command works, it just iterates the registry, finds a match, and delegates. Adding a new command in the future requires no changes to the server, only a new struct and a line in the registry. This follows the open/closed principle: open for extension, closed for modification. Each command is also independently testable and the server dispatch logic never needs to change.
### Queue
The initial implementation used a []string protected by a sync.Mutex. This worked, but it required the consumer to poll, sending POP, getting EMPTY back, sleeping briefly, and retrying until something arrived. That is wasteful both in TCP round trips and in the sleep latency it adds.
Switching to a buffered chan string solved this cleanly. Pop now blocks directly on the channel receive until an item is available, so the consumer gets a response the moment the producer pushes something — no polling, no sleep, no EMPTY response to handle.

The channel also provides backpressure for free — if the producer fills the buffer, Push blocks until the consumer catches up, without any extra code.
The tradeoff is that capacity must be decided upfront.

### Registry
The Registry is a simple slice of Command values that acts as the single source of truth for the protocol. The server, and any future tooling (documentation generators, protocol validators), can derive everything they need from it.
The tradeoff is that registry ordering matters — the first matching command wins. This is not a problem today since PUSH, POP and EOF are unambiguous, but it is worth being aware of if new commands are added with overlapping prefixes.

### EOF command
The producer signals completion by pushing the literal string "EOF" into the queue via a dedicated EOFCommand. The consumer stops when it pops this value rather than writing it to the output file.
The tradeoff is that "EOF" is now a reserved value, a file containing a line with only the text EOF would cause the consumer to stop early. A more robust approach would be a dedicated out-of-band signal (a separate TCP command that closes the queue), but for the scope of this assignment the string is simpler and the limitation is easy to document.

## Scaling strategy
The current design is intentionally simple and suitable for single-file transfers. The main scaling constraint is the queue server, which is in-memory and single-process. Horizontal scaling would require partitioning the queue across multiple server instances and adding persistence for fault tolerance. For high-throughput scenarios, batching multiple line of text per TCP call would significantly reduce round-trip overhead. The queue server also does not have any limit of connections, to resolve that we could use a worker pool of connection using Go channel, only allowing a new one if the pool is not full already.
