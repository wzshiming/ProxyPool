package checker

import (
	"fmt"
	"net"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/wzshiming/requests"
)

type CheckerIPCN struct{}

func (c *CheckerIPCN) Name() string {
	return "ip.cn"
}

func (c *CheckerIPCN) Scheme() string {
	return "https"
}

// Check returns ip from ip.cn
func (c *CheckerIPCN) Check(cli *requests.Client, proxy *url.URL) (net.IP, error) {
	const urlCN = "https://ip.cn/"
	resp, err := cli.SetProxyURL(proxy).NewRequest().Get(urlCN)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.RawBody())
	if err != nil {
		return nil, err
	}

	text := doc.Find("#result").Text()
	text = ipReg.FindString(text)
	ip := net.ParseIP(text)

	if ip == nil {
		return nil, fmt.Errorf("error")
	}
	return ip, nil
}
