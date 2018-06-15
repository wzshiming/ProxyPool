package checker

import (
	"fmt"
	"net"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/wzshiming/requests"
)

type CheckerIP138 struct{}

func (c *CheckerIP138) Name() string {
	return "ip138.com"
}

func (c *CheckerIP138) Scheme() string {
	return "http"
}

// Check returns ip from 2018.ip138.com
func (c *CheckerIP138) Check(cli *requests.Client, proxy *url.URL) (net.IP, error) {
	const urlIP138 = "http://2018.ip138.com/ic.asp"
	resp, err := cli.SetProxyURL(proxy).NewRequest().Get(urlIP138)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.RawBody())
	if err != nil {
		return nil, err
	}

	text := doc.Find("center").Text()
	text = ipReg.FindString(text)
	ip := net.ParseIP(text)

	if ip == nil {
		return nil, fmt.Errorf("error")
	}
	return ip, nil
}
