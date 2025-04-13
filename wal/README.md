# WAL
KnucklesDB implements a Write-Ahead Log (WAL) to ensure durability. By using the WAL, it is possible to restore the memory content simply by replaying the log.

Each Set operation triggered by a client is stored in a buffer. A consumer goroutine periodically collects up to five entries, if available, and writes them to the WAL.

The writing process uses logical offsets to improve performance. Despite this, after a crash, the WAL is read sequentially to restore the memory state.

Before writing a new entry to the WAL, the algorithm checks—based solely on the hash value—whether the entry is already present in the file. This is done by performing a lookup in a secondary hash map that stores the offsets of all entries. If the entry exists, its offset is reused to optimize the write operation. Otherwise, a new offset is calculated and stored in the hash map for future reference.

![wal](https://github.com/user-attachments/assets/0675e26e-b0c2-4203-814e-dfdbdb0bb837)
