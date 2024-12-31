# KnucklesDB
Key Value Store as a Failure Detector

KnucklesDB is a key-value store that serves as a failure detector for distributed environments, such as microservice architectures, horizontal scaling, and unstructured peer-to-peer systems, like cluster-based distributed systems. <br>

Through the paginator algorithm, it is possible for a single replica to manage nearly 3,000 nodes. The failure detector model is based on a tree structure, where heartbeats from nodes are stored with an associated logical clock based on Lamport's clock behavior. Every N seconds, Knuckles schedules the paginator algorithm, which removes all pages associated with a clock value that is 10 less than the root clock

## Key Features
* KnucklesDB is an eventually consistent database.
* KnucklesDB is based on a leaderless architecture.
* KnucklesDB uses version vectors to resolve write-write conflicts and provide anti-entropy.
* KnucklesDB uses a WAL and does not store any items on disk.
* KnucklesDB's crash recovery model is not based on full checkpointing, but rather on reading the WAL.
* KnucklesDB cluster membership is based on the SWIM protocol.


