package logging

import (
	"bufio"
	"fmt"
	"net"
	"testing"
	"time"
)

func init() {
	go server()
}

func server() {
	ln, err := net.Listen("tcp", ":30000")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	b := bufio.NewReader(conn)
	for {
		line, _, err := b.ReadLine()
		if err != nil {
			panic(err)
		}
		fmt.Println(string(line))
	}
}

func TestNetHandler(t *testing.T) {
	h, err := NewNetHandler("tcp", "127.0.0.1:30000", 10*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	AddHandler("net", h)
	Debug("%d, %s", 1, "OK")
	time.Sleep(100 * time.Millisecond)
}
