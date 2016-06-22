//copy from Planb and change some thing to avoid the global values uses
package iplib

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type IpCollection struct {
	Start    uint32
	End      uint32
	CityCode uint32
}

type IpCollections []IpCollection

func Load(ipFile string) (IpCollections, error) {
	ipCollections := make(IpCollections, 0, 160000)
	f, err := os.Open(ipFile)
	if nil == err {
		buff := bufio.NewReader(f)
		for {
			line, _, err := buff.ReadLine()
			if err != nil {
				if io.EOF != err {
					return nil, errors.New("ipfiles load fails while buff.ReadString or meet file ending")
				} else {
					break
				}

			}
			if strings.Count(string(line), ",") < 2 {
				return nil, errors.New("ip file format error")
			}
			ss := strings.SplitN(string(line), ",", -1)
			start := ip2long(ss[0])
			end := ip2long(ss[1])
			citycode, err := strconv.ParseUint(ss[2], 0, 64)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("get city code failed err:%s", err))
			}
			ipCollections = append(ipCollections, IpCollection{uint32(start), uint32(end), uint32(citycode)})
		}
	} else {
		return nil, err
	}
	return ipCollections, nil
}

func ip2long(ipstr string) uint32 {
	ip := net.ParseIP(ipstr)
	if ip == nil {
		log.Println(fmt.Sprintf("ip:%s, parse ip failed", ipstr))
		return 0
	}
	ip = ip.To4()
	if ip == nil {
		log.Println(fmt.Sprintf("ip:%s addr is not ip4!", ipstr))
		return 0
	}
	return binary.BigEndian.Uint32(ip)
}

func (ips IpCollections) binarySearchAux(key uint32, min int, max int) IpCollection {
	if max < min {
		return IpCollection{} //TODO error
	} else {
		mid := (min + max) / 2

		if ips[mid].Start > key {
			return ips.binarySearchAux(key, min, mid-1)
		} else if ips[mid].End < key {
			return ips.binarySearchAux(key, mid+1, max)
		} else {
			return ips[mid]
		}
	}
}

//TODO: to be graceful
func (ips IpCollections) Ip2CityCode(ipStr string) uint32 {
	ipInt := ip2long(ipStr)
	ip := ips.binarySearchAux(ipInt, 0, len(ips))
	return ip.CityCode
}
