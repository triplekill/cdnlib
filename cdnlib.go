package cdnlib

import (
	_ "embed"
	"encoding/json"
	"net/netip"

	"go4.org/netipx"
)

type MatchEngine struct {
	ipv4Set *netipx.IPSet
	ipv6Set *netipx.IPSet
}

func IsCDN(ip string) bool {
	parsedIP := netip.MustParseAddr(ip)
	if parsedIP.Is4() {
		return matchEngine.ipv4Set.Contains(parsedIP)
	}
	return matchEngine.ipv6Set.Contains(parsedIP)
}

//go:embed source_data.json
var data string

type Data struct {
	Ipv4 []string `json:"ipv4"`
	Ipv6 []string `json:"ipv6"`
}

// 需要一个全局变量来存储解析后的数据
var parsedData Data
var matchEngine *MatchEngine

func init() {
	err := json.Unmarshal([]byte(data), &parsedData)
	if err != nil {
		panic(err)
	}
	ipv4Builder := new(netipx.IPSetBuilder)
	ipv6Builder := new(netipx.IPSetBuilder)

	for _, cidr := range parsedData.Ipv4 {
		prefix, err := netip.ParsePrefix(cidr)
		if err != nil {
			panic(err)
		}
		ipv4Builder.AddPrefix(prefix)
	}
	for _, cidr := range parsedData.Ipv6 {
		prefix, err := netip.ParsePrefix(cidr)
		if err != nil {
			panic(err)
		}
		ipv6Builder.AddPrefix(prefix)
	}
	ipv4Set, err := ipv4Builder.IPSet()
	if err != nil {
		panic(err)
	}
	ipv6Set, err := ipv6Builder.IPSet()
	if err != nil {
		panic(err)
	}
	matchEngine = &MatchEngine{
		ipv4Set: ipv4Set,
		ipv6Set: ipv6Set,
	}
}
