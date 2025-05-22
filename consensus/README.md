# Eventual Consistency
KnucklesDB is an eventually consistent distributed database.
Eventual consistency is achieved through the use of a gossip protocol and version vectors.
Version vectors are particularly useful for resolving write-write conflicts, as they allow the system to track and compare different versions of data across nodes.

The gossip protocol employs a data pipeline that batches up to five entries per message. This batching mechanism reduces the frequency and volume of data transmissions across the network, thereby improving overall efficiency.