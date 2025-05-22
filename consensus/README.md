# Eventual Consistency

KnucklesDB is an Eventually Consistent Distributed Database. Eventual Consistency is reached via Gossip Protocol and Version Vectors. 
Version Vectors are very useful in case of write-write conflits.
The Gossip Protocol use a data pipeline to store 5 entries in a row. In this way is possible to reduce the data insertion on the network.