/*
*	This file contains the implementation of the membership protocol.
*	through which nodes can join the cluster
*
**/
package swim

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"slices"
	"gopkg.in/yaml.v3"
)

type ClusterManager struct {
	// this field contains a list of nodes
	// that have joined the cluster.
	clusterMetadata []*Node
	marshaler       ProtocolMarshaer
	// TODO -> gossip field...
}

type SeedNodeMetadata struct {
	SeedNodeAddress    string `yaml:"seed_address"`
	SeedNodeListenPort int    `yaml:"seed_listen_port"`
}

func NewClusterManager() *ClusterManager {
	return &ClusterManager{
		clusterMetadata: make([]*Node, 0),
	}
}

/*
*	@brief ...
*	@param IP address.
*	@param listen port.
**/
func (c *ClusterManager) JoinRequest(host, port string) {
	var ackResult AckMessage
	seedInfo := c.getSeedNodeHostPort()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// by using an infinite loop we create a stubbon link
	for {
		conn, err := net.Dial("tcp", seedInfo)
		if err != nil {
			fmt.Printf("%v", err)
		}
		data, _ := c.marshaler.MarshalJoinMessage(host, port)
		conn.Write(data)

		reply := make([]byte, 2044)

		select {
		case <-ctx.Done():
			conn.Close()
			continue
		default:
			count, _ := conn.Read(reply)
			json.Unmarshal(reply[:count], &ackResult)
		}
		break
	}
}

func (c *ClusterManager) getSeedNodeHostPort() string {
	var seedNodeInfo SeedNodeMetadata
	yamlSeedNodeFile, err := os.Open("join.yaml")
	if err != nil {
		return ""
	}
	defer yamlSeedNodeFile.Close()

	fileData := make([]byte, 1000)
	count, errRead := yamlSeedNodeFile.Read(fileData)
	if errRead != nil {
		return ""
	}

	yaml.Unmarshal(fileData[:count], &seedNodeInfo)

	joined := net.JoinHostPort(seedNodeInfo.SeedNodeAddress, string(seedNodeInfo.SeedNodeListenPort))
	return joined
}

/*
*	@brief this method will be called by the seed server to add new nodes to the cluster.
**/
func (c *ClusterManager) JoinCluster(address net.IP, port int) {
	n := NewNode(address, port, STATUS_ALIVE)
	// Idempotency is achieved through the loop above
	for index, value := range c.clusterMetadata {
		if value.nodeAddress.String() == address.String() && value.nodeListenPort == port {
			c.clusterMetadata = slices.Delete(c.clusterMetadata, index, index)
			break
		}
	}
	c.clusterMetadata = append(c.clusterMetadata, n)

	// TODO -> start gossip cycle
}

/*
*	@brief this method check if the caller is the seed node
*	@param IP address of the caller
*	@param Listen Port of the caller
*	@return result of the operation
*	@return errors occured
**/
func (c *ClusterManager) IsSeed(address string, port int) (bool, error) {
	var (
		yamlSeedNodeFile *os.File
		err              error
		seedNodeInfo     SeedNodeMetadata
	)

	yamlSeedNodeFile, err = os.Open("join.yaml")
	defer yamlSeedNodeFile.Close()

	if err != nil {
		return false, err
	}

	fileData := make([]byte, 1000)
	count, errRead := yamlSeedNodeFile.Read(fileData)
	if errRead != nil {
		return false, errRead
	}

	err = yaml.Unmarshal(fileData[:count], &seedNodeInfo)
	if err != nil {
		return false, err
	}

	fmt.Printf("%s - %d", seedNodeInfo.SeedNodeAddress, seedNodeInfo.SeedNodeListenPort)

	if seedNodeInfo.SeedNodeAddress == address && seedNodeInfo.SeedNodeListenPort == port {
		return true, nil
	}

	return false, nil
}
