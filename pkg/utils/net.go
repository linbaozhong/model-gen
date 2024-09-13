// @Title 网络地址转换
// @Description
// @Author 蔺保仲 2020/04/20
// @Update 蔺保仲 2020/04/20
package utils

import (
	"errors"
	"math/big"
	"net"
	"strconv"
	"strings"

	"github.com/kataras/iris/v12"
)

var (
	localIP    = ""   //本机公网ip
	intranetIP net.IP //本机内网ip
)

func IP6ToInt(ip net.IP) *big.Int {
	ipv6int := big.NewInt(0)
	return ipv6int.SetBytes(ip.To16())
}
func IpToInt(ip string) int64 {
	net_ip := net.ParseIP(ip)
	if len(net_ip) == net.IPv6len {
		return IP6ToInt(net_ip).Int64()
	}
	return INetAtoN(net_ip)
}

//ip地址 string 转 int
func INetAtoN(ip net.IP) int64 {
	bits := strings.Split(ip.String(), ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum int64

	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)

	return sum
}

//ip地址 int 转 string
func INetNtoA(ip int64) net.IP {
	var bytes [4]byte
	bytes[0] = byte(ip & 0xFF)
	bytes[1] = byte((ip >> 8) & 0xFF)
	bytes[2] = byte((ip >> 16) & 0xFF)
	bytes[3] = byte((ip >> 24) & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}

// @desc	获取客户端ip 返回x.xx.xxx.xxxx
// @author	xu.sun
func GetClientIPStr(ctx iris.Context) string {
	ipStr := ctx.RemoteAddr()
	if ipStr == "::1" {
		ipStr = "127.0.0.1"
	}
	return ipStr
}

// @desc	获取客户端ip 返回2130706433
// author	xu.sun
func GetClientIPInt64(ctx iris.Context) int64 {
	ipStr := GetClientIPStr(ctx)
	return INetAtoN(net.ParseIP(ipStr))
}

//GetLocalIP 获得本机外网ip
func GetLocalIP() string {
	if localIP != "" {
		return localIP
	}
	conn, err := net.Dial("udp", "baidu.com:80")
	if err != nil {
		return ""
	}
	defer conn.Close()
	localIP = conn.LocalAddr().(*net.UDPAddr).IP.String()
	return localIP
}

//GetIntranetIP 获取本机内网ip
func GetIntranetIP() (net.IP, error) {
	if intranetIP != nil {
		return intranetIP, nil
	}
	ifaces, e := net.Interfaces()
	if e != nil {
		return nil, e
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network?")
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP

	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() || ip.To4() == nil {
		return nil
	}

	return ip
}
