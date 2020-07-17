package maxminddb

import (
	"net"
)

// Internal structure used to keep track of nodes we still need to visit.
type netNode struct {
	ip      net.IP
	bit     uint
	pointer uint
}

// Networks represents a set of subnets that we are iterating over.
type Networks struct {
	reader   *Reader
	nodes    []netNode // Nodes we still have to visit.
	lastNode netNode
	err      error
}

var allIPv4 = &net.IPNet{IP: make(net.IP, 4), Mask: net.CIDRMask(0, 32)}
var allIPv6 = &net.IPNet{IP: make(net.IP, 16), Mask: net.CIDRMask(0, 128)}

// Networks returns an iterator that can be used to traverse all networks in
// the database.
//
// Please note that a MaxMind DB may map IPv4 networks into several locations
// in an IPv6 database. This iterator will iterate over all of these
// locations separately.
func (r *Reader) Networks() *Networks {
	var networks *Networks
	if r.Metadata.IPVersion == 6 {
		networks = r.NetworksWithin(allIPv6)
	} else {
		networks = r.NetworksWithin(allIPv4)
	}

	return networks
}

// NetworksWithin returns an iterator that can be used to traverse all networks
// in the database which are contained in a given network.
//
// Please note that a MaxMind DB may map IPv4 networks into several locations
// in an IPv6 database. This iterator will iterate over all of these locations
// separately.
//
// If the provided network is contained within a network in the database, the
// iterator will iterate over exactly one network, the containing network.
func (r *Reader) NetworksWithin(network *net.IPNet) *Networks {
	ip := network.IP
	prefixLength, _ := network.Mask.Size()

	if r.Metadata.IPVersion == 6 && len(ip) == net.IPv4len {
		ip = net.IP.To16(ip)
		prefixLength += 96
	}

	pointer, bit := r.traverseTree(ip, 0, uint(prefixLength))
	return &Networks{
		reader: r,
		nodes: []netNode{
			{
				ip:      ip,
				bit:     uint(bit),
				pointer: pointer,
			},
		},
	}
}

// Next prepares the next network for reading with the Network method. It
// returns true if there is another network to be processed and false if there
// are no more networks or if there is an error.
func (n *Networks) Next() bool {
	for len(n.nodes) > 0 {
		node := n.nodes[len(n.nodes)-1]
		n.nodes = n.nodes[:len(n.nodes)-1]

		for node.pointer != n.reader.Metadata.NodeCount {
			if node.pointer > n.reader.Metadata.NodeCount {
				n.lastNode = node
				return true
			}
			ipRight := make(net.IP, len(node.ip))
			copy(ipRight, node.ip)
			if len(ipRight) <= int(node.bit>>3) {
				n.err = newInvalidDatabaseError(
					"invalid search tree at %v/%v", ipRight, node.bit)
				return false
			}
			ipRight[node.bit>>3] |= 1 << (7 - (node.bit % 8))

			offset := node.pointer * n.reader.nodeOffsetMult
			rightPointer := n.reader.nodeReader.readRight(offset)

			node.bit++
			n.nodes = append(n.nodes, netNode{
				pointer: rightPointer,
				ip:      ipRight,
				bit:     node.bit,
			})

			node.pointer = n.reader.nodeReader.readLeft(offset)
		}
	}

	return false
}

// Network returns the current network or an error if there is a problem
// decoding the data for the network. It takes a pointer to a result value to
// decode the network's data into.
func (n *Networks) Network(result interface{}) (*net.IPNet, error) {
	if err := n.reader.retrieveData(n.lastNode.pointer, result); err != nil {
		return nil, err
	}

	return &net.IPNet{
		IP:   n.lastNode.ip,
		Mask: net.CIDRMask(int(n.lastNode.bit), len(n.lastNode.ip)*8),
	}, nil
}

// Err returns an error, if any, that was encountered during iteration.
func (n *Networks) Err() error {
	return n.err
}
