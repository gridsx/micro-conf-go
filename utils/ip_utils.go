package utils

import (
	"fmt"
	"net"
	"strings"
)

var (
	ip = getHostIp()
)

func getHostIp() string {
	conn, err := net.Dial("udp", "1.2.4.8:53")
	defer conn.Close()
	if err != nil {
		fmt.Println("get current host ip err: ", err)
		return ""
	}
	addr := conn.LocalAddr().(*net.UDPAddr)
	ip := strings.Split(addr.String(), ":")[0]
	return ip
}

func GetIP() string {
	return ip
}
