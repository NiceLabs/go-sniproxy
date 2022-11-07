package sniproxy

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"net/http"
	"strconv"
	"sync"
)

type Dial func(network, address string) (net.Conn, error)

func ForwardHTTP(conn net.Conn, dial Dial) {
	defer conn.Close()
	var peeked bytes.Buffer
	request, err := http.ReadRequest(bufio.NewReader(io.TeeReader(conn, &peeked)))
	if err != nil {
		return
	}
	port := conn.LocalAddr().(*net.TCPAddr).Port
	host := net.JoinHostPort(request.Host, strconv.Itoa(port))
	remote, err := dial("tcp", host)
	if err != nil {
		return
	}
	_, _ = remote.Write(peeked.Bytes())
	var wg sync.WaitGroup
	wg.Add(2)
	go copyConn(conn, remote, &wg)
	go copyConn(remote, conn, &wg)
	wg.Wait()
	return
}

func ForwardTLS(conn net.Conn, dial Dial) {
	defer conn.Close()
	var peeked bytes.Buffer
	serverName := readClientHello(io.TeeReader(conn, &peeked))
	if serverName == "" {
		return
	}
	port := conn.LocalAddr().(*net.TCPAddr).Port
	host := net.JoinHostPort(serverName, strconv.Itoa(port))
	remote, err := dial("tcp", host)
	if err != nil {
		err = conn.Close()
		return
	}
	_, _ = remote.Write(peeked.Bytes())
	var wg sync.WaitGroup
	wg.Add(2)
	go copyConn(conn, remote, &wg)
	go copyConn(remote, conn, &wg)
	wg.Wait()
	return
}
