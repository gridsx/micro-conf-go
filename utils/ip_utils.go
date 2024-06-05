package utils

import (
	"fmt"
	"math/rand"
	"net"
	"strings"

	"github.com/winjeg/go-commons/log"
)

var (
	ip     = getHostIp()
	logger = log.GetLogger(nil)
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

func ChooseAddr(metaServers string) string {
	addrArr := strings.Split(metaServers, ",")
	if len(addrArr) < 1 {
		logger.Panic("ChooseAddr - meta server address incorrect")
	}
	idx := rand.Intn(len(addrArr))
	return addrArr[idx]
}
