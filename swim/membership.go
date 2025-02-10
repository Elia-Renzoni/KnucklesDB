/*
*	This file contains the implementation of the membership protocol.
*	through which nodes can join the cluster
*
**/
package swim

import (
	"fmt"
	"net"
	"os"

	"gopkg.in/yaml.v3"
)

type ClusterManager struct {
	// this field contains a list of nodes
	// that have joined the cluster.
	clusterMetadata []*Node
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
*	@brief this method will be called by the new nodes to join the cluster.
**/
func (c *ClusterManager) JoinCluster(address net.IP, port int) {
	n := NewNode(address, port, STATUS_ALIVE)
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
