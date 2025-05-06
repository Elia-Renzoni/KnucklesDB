/*
*	This file contains the implementation of the membership protocol.
*	through which nodes can join the cluster
*
**/
package swim

import (
	"context"
	"encoding/json"
	"knucklesdb/wal"
	"math/rand"
	"net"
	"os"
	"slices"
	"strconv"
	"time"
	"fmt"
	"sync"

	"gopkg.in/yaml.v3"
)

type ClusterManager struct {
	// this field contains a list of nodes
	// that have joined the cluster.
	cluster        *Cluster
	syncronizer    *sync.WaitGroup
	marshaler      ProtocolMarshaer
	logger         *wal.ErrorsLogger
	gossipSpreader *Dissemination
}

type SeedNodeMetadata struct {
	SeedNodeAddress    string `yaml:"seed_address"`
	SeedNodeListenPort int    `yaml:"seed_listen_port"`
}

func NewClusterManager(syncronizer *sync.WaitGroup, cluster *Cluster, logger *wal.ErrorsLogger, gossipProtocol *Dissemination) *ClusterManager {
	return &ClusterManager{
		cluster:        cluster,
		syncronizer: syncronizer,
		logger:         logger,
		gossipSpreader: gossipProtocol,
	}
}

/*
*	@brief Method Called by the new nodes to enter the cluster
*	@param IP address.
*	@param listen port.
**/
func (c *ClusterManager) JoinRequest(host, port string) {
	// wait for the listener
	c.syncronizer.Wait()

	var ackResult AckMessage
	seedInfo := c.getSeedNodeHostPort()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// by using an infinite loop we create a stubborn link
	for {
		conn, err := net.Dial("tcp", seedInfo)
		if err != nil {
			c.logger.ReportError(err)
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

	joined := net.JoinHostPort(seedNodeInfo.SeedNodeAddress, strconv.Itoa(seedNodeInfo.SeedNodeListenPort))
	return joined
}

/*
*	@brief this method will be called by the seed server to add new nodes to the cluster.
**/
func (c *ClusterManager) JoinCluster(address string, port int) {
	n := NewNode(address, port, STATUS_ALIVE)
	// Idempotency is achieved through the loop above
	for index, value := range c.cluster.clusterMetadata {
		if value.nodeAddress == address && value.nodeListenPort == port {
			c.cluster.clusterMetadata = slices.Delete(c.cluster.clusterMetadata, index, index+1)
			break
		}
	}
	c.cluster.clusterMetadata = append(c.cluster.clusterMetadata, n)

	fanoutNodeList := c.SetFanoutList()
	fmt.Println(fanoutNodeList)
	go c.gossipSpreader.SpreadMembershipList(c.cluster.clusterMetadata, fanoutNodeList, true)
}

func (c *ClusterManager) CheckIfFanoutIsPossible(senderAddr string, myAddr string) bool {
	var (
		result bool = true
		counter int = 0
	) 

	for _, replica := range c.cluster.clusterMetadata {
		castedPort := strconv.Itoa(replica.nodeListenPort)
		replicaJoinedAddr := net.JoinHostPort(replica.nodeAddress, castedPort)
		if replicaJoinedAddr == senderAddr || replicaJoinedAddr == myAddr || replicaJoinedAddr == "127.0.0.1:5050"{
			counter += 1
		}
	}

	fmt.Println(counter)

	fmt.Println(len(c.cluster.clusterMetadata))
	if counter == len(c.cluster.clusterMetadata) {
		result = false
	}

	return result
}

func (c *ClusterManager) DeleteNodeFromCluster(address, port string) {
	castedPort, _ := strconv.Atoi(port)
	for index := range c.cluster.clusterMetadata {
		if c.cluster.clusterMetadata[index].nodeAddress == address && c.cluster.clusterMetadata[index].nodeListenPort == castedPort {
			c.cluster.clusterMetadata = slices.Delete(c.cluster.clusterMetadata, index, index+1)
			break
		}
	}
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

	if seedNodeInfo.SeedNodeAddress == address && seedNodeInfo.SeedNodeListenPort == port {
		return true, nil
	}

	return false, nil
}

func (c *ClusterManager) SetFanout() (fanoutFactor int) {

	if len(c.cluster.clusterMetadata) <= 7 {
		fanoutFactor = 1
	} else {
		fanoutFactor = 3
	}

	return
}

func (c *ClusterManager) SetFanoutList() []string {
	var fanoutNodeList []string = make([]string, 0)

	for i := 0; i < c.SetFanout(); i++ {
		var selectedNode = rand.Intn(len(c.cluster.clusterMetadata))
		port := strconv.Itoa(c.cluster.clusterMetadata[selectedNode].nodeListenPort)
		fanoutNodeList = append(fanoutNodeList, net.JoinHostPort(c.cluster.clusterMetadata[selectedNode].nodeAddress, port))
	}
	
	for indexNode := range fanoutNodeList {
		if c.isRedundant(fanoutNodeList[indexNode], fanoutNodeList) {
			fanoutNodeList = slices.Delete(fanoutNodeList, indexNode, indexNode)
		}
	}

	return fanoutNodeList
}

func (c *ClusterManager) isRedundant(nodeToSearch string, fanoutNodeList []string) bool {
	var redundancyCounter int

	for _, node := range fanoutNodeList {
		if node == nodeToSearch {
			redundancyCounter += 1
		}
	}

	if redundancyCounter >= 2 {
		return true
	}

	return false
}
