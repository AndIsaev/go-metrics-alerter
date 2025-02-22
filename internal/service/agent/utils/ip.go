package utils

import (
	"net"
)

type IPResolver interface {
	GetLocalIP() (string, error)
}

type DefaultIPResolver struct {
	Address string
}

func NewDefaultIPResolver(address string) *DefaultIPResolver {
	return &DefaultIPResolver{Address: address}
}

// GetLocalIP get your ip address
func (r *DefaultIPResolver) GetLocalIP() (string, error) {
	addr, err := net.ResolveTCPAddr("tcp", r.Address)
	if err != nil {
		return "", err
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.TCPAddr)

	return localAddr.IP.String(), nil
}
