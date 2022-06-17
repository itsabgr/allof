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
var CmdPubFlagIP = CmdPub.String("ip", "", "ip version")
var CmdPubNodes = CmdPub.Args("peers", "")

//
var CmdSub = Cmd.SubCommand("sub", "subscribe messages")
var CmdSubFlagAddr = CmdSub.String("a", "", "addr")
var CmdSubFlagTopic = CmdSub.String("t", "", "topic")
var CmdSubFlagIP = CmdSub.String("ip", "", "ip version")
var CmdSubNodes = CmdSub.Args("peers", "")

func main() {
	handy.Throw(Cmd.ParseArgs(os.Args...))
	switch {
	case CmdPub.Parsed():
		node, err := allof.Listen("udp"+(*CmdPubFlagIP), ":0", []byte(*CmdPubFlagTopic))
		handy.Throw(err)
		for _, peer := range *CmdPubNodes {
			addrPort := netip.MustParseAddrPort(peer)
			node.Add(addrPort.Addr().AsSlice(), int(addrPort.Port()))
		}
		handy.Throw(node.Send(200, []byte(*CmdPubFlagMsg)))
	case CmdSub.Parsed():
		node, err := allof.Listen("udp"+(*CmdSubFlagIP), *CmdSubFlagAddr, []byte(*CmdSubFlagTopic))
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
