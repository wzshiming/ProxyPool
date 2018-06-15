package collector

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/wzshiming/requests"
)

type CollectorGoubanjia struct{}

func (c *CollectorGoubanjia) Name() string {
	return "goubanjia.com"
}

// Collect returns proxys from www.goubanjia.com
func (c *CollectorGoubanjia) Collect(cli *requests.Client) (result []*url.URL, err error) {
	const urlGoubanjia = "http://www.goubanjia.com/"
	resp, err := cli.NewRequest().Get(urlGoubanjia)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.RawBody())
	if err != nil {
		return nil, err
	}

	doc.Find("tbody > tr").Each(func(i int, s *goquery.Selection) {
		s = s.Children()
		addrs := []string{}
		s.Eq(0).Children().Each(func(i int, s *goquery.Selection) {
			//	ffmt.Mark(s.AttrOr("style", ""))
			if !strings.Contains(s.AttrOr("style", ""), "none") {
				addrs = append(addrs, s.Text())
			}
		})
		addrs = append(addrs[:len(addrs)-1], ":", addrs[len(addrs)-1])
		addr := strings.Join(addrs, "")
		scheme := s.Eq(2).Text()
		proxyStr := fmt.Sprintf("%s://%s", scheme, addr)
		proxy, err := url.Parse(proxyStr)
		if err == nil {
			result = append(result, proxy)
		}
	})
	return
}
