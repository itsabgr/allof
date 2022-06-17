package allof

import (
	"context"
	"fmt"
	"golang.org/x/net/ipv4"
	"net"
	"net/netip"
	"sync"
)

type Node struct {
	peers []*net.UDPAddr
	topic []byte
	mutex sync.RWMutex
	conn  *ipv4.PacketConn
}

func Listen(network, addr string, topic []byte) (*Node, error) {
	switch network {
	case "udp", "udp4", "udp6":
	default:
		panic(fmt.Errorf("unsupported network %q", network))
	}
	pConn, err := net.ListenPacket(network, addr)
	if err != nil {
		return nil, err
	}

	n := &Node{
		conn:  ipv4.NewPacketConn(pConn),
		peers: make([]*net.UDPAddr, 0, 100),
		topic: topic,
	}
	return n, nil
}
func (node *Node) Addr() netip.AddrPort {
	addr := node.conn.PacketConn.LocalAddr()
	return netip.MustParseAddrPort(addr.String())
}
func (node *Node) Ping() error {
	return node.send(Encode(0, Encode(200, []byte{}, node.topic), []byte("ping")))
}
func (node *Node) Send(ttl byte, msg []byte) error {
	return node.send(Encode(ttl, msg, node.topic))
}
func (node *Node) send(data []byte) error {
	node.mutex.RLock()
	MsgS := make([]ipv4.Message, len(node.peers))
	for i, addr := range node.peers {
		MsgS[i].Addr = addr
		MsgS[i].Buffers = net.Buffers{data}
	}
	node.mutex.RUnlock()
	if len(MsgS) == 0 {
		return nil
	}
	_, err := node.conn.WriteBatch(MsgS, 0)
	return err
}

func (node *Node) Recv(ctx context.Context) ([]byte, error) {
	b := make([]byte, 512)
	for {
		err := ctx.Err()
		if err != nil {
			return nil, err
		}
		n, _, from, err := node.conn.ReadFrom(b)
		if err != nil {
			return nil, err
		}
		pack := b[:n]
		ttl, data, err := Decode(pack, node.topic)
		if err == nil {
			if ttl > 0 {
				_ = node.Send(ttl-1, data)
			}
			return data, nil
		}
		_, topic, err := Decode(pack, []byte("ping"))
		if err != nil {
			continue
		}
		ttl, addr, err := Decode(topic, node.topic)
		if err != nil {
			continue
		}
		if len(addr) == 0 {
			addr = []byte(from.String())
		}
		addrPort, err := netip.ParseAddrPort(string(addr))
		if err != nil {
			continue
		}
		_ = node.send(Encode(0, Encode(ttl-1, addr, node.topic), []byte("ping")))
		node.Add(addrPort.Addr().AsSlice(), int(addrPort.Port()))
	}
}
func (node *Node) Add(ip net.IP, port int) {
	if ip.Equal(node.Addr().Addr().AsSlice()) {
		return
	}
	node.mutex.Lock()
	defer node.mutex.Unlock()
	for i, peer := range node.peers {
		if peer.IP.Equal(ip) && peer.Port == port {
			node.peers[i].IP = ip
			node.peers[i].Port = port
			return
		}
	}
	addr := &net.UDPAddr{Port: port, IP: ip}
	if len(node.peers) >= cap(node.peers) {
		node.peers = append(node.peers[1:], addr)
	} else {
		node.peers = append(node.peers, addr)
	}
}

func (node *Node) Close() error {
	return node.conn.PacketConn.Close()
}
