package pkg

import (
	"bytes"
	"fmt"
	"slices"
	"strconv"
)

type ConsistentHashRing struct {
	// Number of Labels per entry (vnodes per node)
	Labels int
	// Slice of Keys
	Keys [][]byte
	// Slice of Node
	Nodes []Node
}

func NewConsistentHashRing(labels int) *ConsistentHashRing {
	return &ConsistentHashRing{
		Labels: labels,
	}
}

// String approximates a pretty-print of the contents of the ConsistentHashRing
// XXX: Maybe do a not-pretty-print variant as a basis for a PrettyPrint function.
func (c *ConsistentHashRing) String() string {
	var a, b string
	for _, v := range c.Nodes {
		a += "\t\t\t" + v.String()
	}
	for _, v := range c.Keys {
		b += "\t\t\t" + string(v) + "\n"
	}
	return fmt.Sprintf("{\n\tLabels: %d\n\tKeys: {\n%s\t}\n\tNodes: {\n%s\t}\n}", c.Labels, b, a)
}

type Node struct {
	Hash []byte
	Host string
	IP   string
	SID  string
}

func NewNode(host string, ip string, sid string) Node {
	return Node{
		Hash: []byte{},
		Host: host,
		IP:   ip,
		SID:  sid,
	}
}

func (n Node) String() string {
	result := make([]byte, len(n.Hash))
	buff := bytes.NewBuffer(result)
	for _, b := range n.Hash {
		fmt.Fprintf(buff, "0x%02x ", b)
	}
	return fmt.Sprintf("Hash: {%s} Host: %s IP: %s SID: %s\n", buff.String(), n.Host, n.IP, n.SID)
}

// Add adds a node given its name.
// The given nodeName is hashed among the number of labels.
func (c *ConsistentHashRing) Add(nodeName string, node Node) {
	for i := 0; i < c.Labels; i++ {
		hash := CalculateHash(nodeName + strconv.Itoa(i))
		node.Hash = hash
		c.Nodes = append(c.Nodes, node)
		c.Keys = SortedInsertByte(c.Keys, hash)
	}
}

// Get returns a node given a key.
// The node replica with a hash value nearest but not
// less than that of the given name is returned. If the hash
// of the given name is greater than the greatest hash,
// return the lowest hashed node.
// If the Hashring is empty or any other case happens,
// return an empty Node type.
func (c *ConsistentHashRing) Get(keyname string) *Node {
	// if empty, return empty
	if len(c.Nodes) == 0 {
		return &Node{}
	}
	hash := CalculateHash(keyname)
	idx, ok := binarySearchBytes(c.Keys, hash, 0, len(c.Keys)-1)
	if !ok {
		//return &Node{}
		idx = 0
	}
	if idx == len(c.Keys) {
		idx = 0
	}
	var x []byte
	if len(c.Keys) > 0 {
		x = c.Keys[idx]
	} else {
		return &Node{}
	}
	for k, v := range c.Nodes {
		if bytes.Equal(v.Hash, x) {
			return &c.Nodes[k]
		}
	}
	return &Node{}
}

// Delete deletes a node given its name.
func (c *ConsistentHashRing) Delete(nodeName string) {
	for i := 0; i < c.Labels; i++ {
		hash := CalculateHash(nodeName + strconv.Itoa(i))
		// delete from c.Nodes where hash matches
		for k, v := range c.Nodes {
			if bytes.Equal(v.Hash, hash) {
				c.Nodes = slices.Delete(c.Nodes, k, k+1)
			}
		}
		// delete from c.Keys
		for k, val := range c.Keys {
			if bytes.Equal(val, hash) {
				c.Keys = slices.Delete(c.Keys, k, k+1)
			}
		}
	}
}
