package utils

import "net"

type IPResolver interface {
	GetLocalIP(address string) (string, error)
}

type DefaultIPResolver struct{}

func NewDefaultIPResolver() *DefaultIPResolver {
	return &DefaultIPResolver{}
}

// GetLocalIP get your ip address
func (resolver *DefaultIPResolver) GetLocalIP(address string) (string, error) {
	addr, err := net.ResolveTCPAddr("tcp", address)
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
