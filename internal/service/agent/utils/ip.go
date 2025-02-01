package utils

import "net"

// GetLocalIP get your ip address
func GetLocalIP(address string) (string, error) {
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
