package utils

import (
	"crypto/md5"
	"fmt"
	"github.com/google/uuid"
	"net"
	"strconv"
	"time"
)

func UUID() string {
	return uuid.New().String()
}

func Time(t time.Time) string {
	t.UTC()
	return t.Format("2006/01/02 15:04:05")
}

// Retry 自动重试
func Retry(retry int, f func() error) error {
	var err error
	err = f()
	if err == nil {
		return nil
	}
	for re := 0; re < retry; re++ {
		err = f()
		if err == nil {
			break
		}
	}
	return err
}

func ParseInt32(s string) (int32, error) {
	n, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(n), nil
}

func GetLocalIPs() (ipv4s []string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && !ipnet.IP.IsPrivate() {
			if ipnet.IP.To4() != nil {
				ipv4s = append(ipv4s, ipnet.IP.String())
			}
		}
	}
	return ipv4s
}

// GetPodIP 获取本地出口地址
func GetPodIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "127.0.0.1", err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	if localAddr.IP.To4() != nil && !localAddr.IP.IsLoopback() {
		return localAddr.IP.String(), nil
	}
	return "127.0.0.1", nil
}

// GetLocalIP 获取本机IP
func GetLocalIP() (string, error) {
	adds, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1", err
	}
	for _, address := range adds {
		// 检查ip地址判断是否回环地址
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.IsGlobalUnicast() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}
	return "127.0.0.1", nil
}

func Hash(data []byte) string {
	return fmt.Sprintf("%x", md5.Sum(data))
}
