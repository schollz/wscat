package main

import (
	"bufio"
	"flag"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var flagAddr = flag.String("addr", "localhost", "address")
var flagPort = flag.String("port", "5555", "port")
var flagSubProtocol = flag.String("sub", "bus.sp.nanomsg.org", "sub protocol")
var flagSend = flag.String("send", "<wscat>", "format sent piped input")

// wscat --send 'last_command=[[<wscat>]];<wscat>'
func main() {
	flag.Parse()
	u := url.URL{Scheme: "ws", Host: *flagAddr + ":" + *flagPort, Path: "/"}
	var cstDialer = websocket.Dialer{
		Subprotocols:     []string{*flagSubProtocol},
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		HandshakeTimeout: 3 * time.Second,
	}

	norns, _, err := cstDialer.Dial(u.String(), nil)
	if err != nil {
		os.Exit(1)
	}
	defer norns.Close()
	stdin := bufio.NewScanner(os.Stdin)
	s := ""
	for stdin.Scan() {
		s += stdin.Text()
	}
	norns.WriteMessage(websocket.TextMessage, []byte(strings.Replace(*flagSend, "<wscat>", s, 1)+"\n"))
}
