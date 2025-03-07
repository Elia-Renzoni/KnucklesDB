# KnucklesDB
Key Value Store as a Failure Detector

KnucklesDB is a key-value store that serves as a failure detector for distributed environments, such as microservice architectures, horizontal scaling, and unstructured peer-to-peer systems, like cluster-based distributed systems. <br>

Through the paginator algorithm, it is possible for a single replica to manage nearly 3,000 nodes. The failure detector model is based on the Clock Pagination Algorithm.

## Key Features
* KnucklesDB is an eventually consistent database.
* KnucklesDB is based on a leaderless architecture.
* KnucklesDB uses version vectors to resolve write-write conflicts and provide anti-entropy.
* KnucklesDB uses a WAL and does not store any items on disk.
* KnucklesDB's crash recovery model is not based on full checkpointing, but rather on reading the WAL.
* KnucklesDB cluster membership is based on the SWIM protocol.

## Eventual Consistency
KnucklesDB is a distributed, eventually consistent database, meaning it supports the weakest data consistency model. It can therefore be considered a distributed AP system that, in the event of a network partition, allows update operations on the database state at the cost of returning stale reads. Eventual consistency is achieved through the use of the gossip protocol, which is naturally eventually consistent; additionally, KnucklesDB also utilizes version vectors to handle write-write conflicts.