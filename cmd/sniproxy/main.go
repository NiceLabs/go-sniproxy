package main

import (
	"flag"
	"log"
	"net"

	"golang.org/x/net/proxy"

	"github.com/NiceLabs/go-sniproxy"
)

func main() {
	listenHTTPS := flag.String("https", "localhost:https", "Listen HTTPS address")
	proxyAddress := flag.String("socks5-proxy", "localhost:7890", "Proxy Server address")
	flag.Parse()
	dial, err := proxy.SOCKS5("tcp", *proxyAddress, nil, proxy.Direct)
	if err != nil {
		log.Panicln(err)
	}
	if err = ListenHTTPS(*listenHTTPS, dial.Dial); err != nil {
		log.Panicln(err)
	}
}

func ListenHTTPS(address string, dial sniproxy.Dial) (err error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go sniproxy.ForwardTLS(conn, dial)
	}
}
