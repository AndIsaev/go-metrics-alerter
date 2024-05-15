package main

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

var flagRunAddr string

type agentAddress struct {
	host string
	port int
}

func (a *agentAddress) String() string {
	return a.host + ":" + strconv.Itoa(a.port)
}

func (a *agentAddress) Set(s string) error {
	hp := strings.Split(s, ":")

	if len(hp) != 2 {
		return errors.New("need address in a form host:port")
	}
	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}
	a.host = hp[0]
	a.port = port
	return nil
}

func parseFlags() {
	addr := new(agentAddress)

	_ = flag.Value(addr)
	flag.Var(addr, "a", "Net address host:port")
	flag.Parse()

	if addr.port == 0 {
		flagRunAddr = "http://localhost:8080"
		return
	}
	flagRunAddr = fmt.Sprintf("http://%v:%v", addr.host, addr.port)

}
