# SWIM
 The protocol relies on periodic pings, executed in a round-robin fashion by each node in the cluster. If a node fails to respond to a ping, the sender initiates a "help request" to K randomly selected nodes. These nodes then attempt to ping the suspected node and report the outcome back to the original sender.

If the target node responds successfully, it is no longer considered suspicious. Otherwise, it is marked for removal from the cluster.

To propagate cluster state information, each node periodically performs gossip rounds.