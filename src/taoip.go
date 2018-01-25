package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
)

const (
	VERSION = "1.0"
)

const (
	TAOBAO_IP_QUERY_URL      = "http://ip.taobao.com/service/getIpInfo2.php"
	TAOBAO_HTTP_CONTENT_TYPE = "application/x-www-form-urlencoded"
)

type TaobaoIPQueryResp struct {
	Code int
	Data TaobaoIP
}

type TaobaoIP struct {
	Country_id  string
	County_id   string
	Isp         string
	Area        string
	Area_id     string
	City_id     string
	QueryFromIp string `json:"ip"`
	Region_id   string
	Region      string
	City        string
	County      string
	Isp_id      string
	Country     string
}

func usage(code int) {
	fmt.Printf(
		`taoip %s 
taoip: query IPv4 description from taobao
Usage: taoip HostName
`, VERSION)
	os.Exit(code)
}

func http_post(url string, bodyType string, cmd_body string) (resp_body []byte, err error) {
	resp, err := http.Post(url, bodyType, strings.NewReader(cmd_body))
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	defer resp.Body.Close()
	resp_body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	//log.Println(string(resp_body))
	return resp_body, nil
}

func queryIPFromTaobao(ip string) (ip_desc *TaobaoIP, err error) {

	post_body := fmt.Sprintf("ip=%s", ip)
	resp, err := http_post(TAOBAO_IP_QUERY_URL, TAOBAO_HTTP_CONTENT_TYPE, post_body)
	if err != nil {
		fmt.Println("get device failed", err.Error())
		return nil, err
	}

	var queryResp TaobaoIPQueryResp
	err = json.Unmarshal(resp, &queryResp)
	if err != nil {
		fmt.Printf("query failed, code:%d\n", queryResp.Code)
		return nil, err
	}

	//log.Println(string(resp))

	return &queryResp.Data, nil
}

func main() {
	arg_num := len(os.Args)
	if arg_num != 2 {
		usage(0)
	}

	ip_str := os.Args[1]
	ns, err := net.LookupHost(ip_str)
	if err != nil {
		fmt.Printf("taoip: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("%s:\n", ip_str)
	for _, ip_addr := range ns {
		// skip IPv6
		ip := net.ParseIP(ip_addr)
		if ip == nil || ip.To4() == nil {
			continue
		}

		ip_desc, err := queryIPFromTaobao(ip_addr)
		if err != nil {
			os.Exit(1)
		}

		fmt.Printf("\t%s [%s--%s--%s--%s--%s]\n",
			ip_addr,
			ip_desc.Country,
			ip_desc.Region,
			ip_desc.City,
			ip_desc.County,
			ip_desc.Isp)
	}

}
