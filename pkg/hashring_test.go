package pkg_test

import (
	"github.com/nikkicoon/consistenthashring/pkg"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConsistentHashRing_Get(t *testing.T) {
	t.Parallel()
	cr := pkg.NewConsistentHashRing(5)
	n := pkg.NewNode("foobar.de", "127.0.0.1", "0")
	cr.Add("node0", n)
	hashesNode0 := [][]byte{
		{0x25, 0xc3, 0xcf, 0xcb, 0x8b, 0x52, 0x0d, 0xbc, 0x5f, 0x85, 0xe9, 0xfb, 0x16, 0xe8, 0x6b, 0xaa, 0x13, 0x9f, 0x7a, 0x99},
		{0xaa, 0x99, 0x2f, 0xdc, 0x24, 0x46, 0xe7, 0x27, 0xc3, 0x30, 0xa3, 0xcc, 0x4b, 0x9c, 0x2f, 0x46, 0x7e, 0x43, 0x8d, 0x6e},
		{0x68, 0x4e, 0x61, 0xda, 0xe6, 0xe5, 0x5c, 0xfc, 0x61, 0x1a, 0x2b, 0xbd, 0xa6, 0x99, 0x4a, 0x84, 0xfd, 0xab, 0x91, 0x4e},
		{0x43, 0xc9, 0x45, 0xef, 0x85, 0x95, 0x9c, 0x47, 0x12, 0x1d, 0x74, 0x31, 0x0c, 0xfa, 0x9f, 0xa9, 0x34, 0xa7, 0x30, 0xb3},
		{0x28, 0x4b, 0xff, 0x8f, 0x21, 0xb8, 0x41, 0x0c, 0x4e, 0xb8, 0x84, 0x49, 0x26, 0xe6, 0x63, 0x60, 0x87, 0x99, 0xd0, 0xc2}}
	q := cr.Get("node0")
	if q != nil {
		assert.Equal(t, "foobar.de", q.Host)
		assert.Equal(t, 5, len(cr.Nodes))
		for i := 0; i < len(cr.Nodes); i++ {
			assert.Equal(t, hashesNode0[i], cr.Nodes[i].Hash)
		}
	} else {
		panic("some error with get('node0')")
	}
	cr.Delete("node0")
	l := cr.Get("node9000")
	assert.IsType(t, &pkg.Node{}, l)
	assert.Equal(t, pkg.Node{}, *l)
}

func TestConsistentHashRing_Add(t *testing.T) {
	cr := pkg.NewConsistentHashRing(2000)
	cr.Add("node0", pkg.NewNode("host0.domain.tld", "123.123.123.123", "0"))
	cr.Add("node1", pkg.NewNode("host1.domain.tld", "123.123.123.123", "0"))
	cr.Add("node2", pkg.NewNode("host2.domain.tld", "123.123.123.123", "0"))
	cr.Add("node3", pkg.NewNode("host3.domain.tld", "123.123.123.123", "0"))
	cr.Add("node4", pkg.NewNode("host4.domain.tld", "123.123.123.123", "0"))
	cr.Add("node5", pkg.NewNode("host5.domain.tld", "123.123.123.123", "0"))
	for _, res := range cr.Nodes {
		assert.NotZero(t, res.Hash)
		assert.NotNil(t, res.Hash)
	}
	assert.NotZero(t, len(cr.Nodes))
	assert.NotZero(t, len(cr.Keys))
}
