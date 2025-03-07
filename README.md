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

## Leaderless Architecture
KnucklesDB was designed and implemented with a leaderless architecture, meaning it consists of so-called flat groups—process groups in which replicas do not have well-defined roles relative to one another. Leader-based architectures are essential for ensuring a stronger data consistency model, such as Linearizability, Serializability, or even Sequential Consistency. However, since KnucklesDB is an eventually consistent database, using a primary-backup architecture would not have made sense. Such an approach would have also introduced additional overhead due to the need for distributed consensus protocols that assume fail-noisy conditions regarding the leader.

Thanks to its architecture, KnucklesDB can be scaled horizontally with great ease and simplicity. This is particularly useful in unstructured P2P networks, where database failures could result in the loss of the entire network.

## Clock Algorithm
KnucklesDB uses a paging algorithm that is an implementation of the clock algorithm with a second chance. This algorithm is scheduled every three seconds. In the first phase, it sets the flag related to each entry in memory to false. In the second phase, it checks the flags one last time—entries with a false flag will be removed from memory.

This algorithm is particularly useful because it allows clients to set a heartbeat time that can be relatively high, such as four seconds.