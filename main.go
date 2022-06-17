package main

import (
	"context"
	"fmt"
	"github.com/itsabgr/allof/pkg/allof"
	"github.com/itsabgr/go-handy"
	"github.com/posener/cmd"
	"net/netip"
	"os"
)

var Cmd = cmd.New()

//
var CmdPub = Cmd.SubCommand("pub", "publish message")
var CmdPubFlagMsg = CmdPub.String("m", "", "message")
var CmdPubFlagTopic = CmdPub.String("t", "", "topic")
var CmdPubNodes = CmdPub.Args("peers", "")

//
var CmdSub = Cmd.SubCommand("sub", "subscribe messages")
var CmdSubFlagAddr = CmdSub.String("a", "0.0.0.0:0", "addr")
var CmdSubFlagTopic = CmdSub.String("t", "", "topic")
var CmdSubNodes = CmdSub.Args("peers", "")

func main() {
	handy.Throw(Cmd.ParseArgs(os.Args...))
	switch {
	case CmdPub.Parsed():
		node, err := allof.Listen("udp", "0.0.0.0:0", []byte(*CmdPubFlagTopic))
		handy.Throw(err)
		fmt.Println(node.Addr())
		for _, peer := range *CmdPubNodes {
			addrPort := netip.MustParseAddrPort(peer)
			node.Add(addrPort.Addr().AsSlice(), int(addrPort.Port()))
		}
		handy.Throw(node.Send(200, []byte(*CmdPubFlagMsg)))
	case CmdSub.Parsed():
		node, err := allof.Listen("udp", *CmdSubFlagAddr, []byte(*CmdSubFlagTopic))
		handy.Throw(err)
		fmt.Println(node.Addr())
		for _, peer := range *CmdSubNodes {
			addrPort := netip.MustParseAddrPort(peer)
			node.Add(addrPort.Addr().AsSlice(), int(addrPort.Port()))
		}
		for {
			data, err := node.Recv(context.Background())
			handy.Throw(err)
			fmt.Println(string(data))
		}
	}
}
