package collector

import (
	"fmt"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/wzshiming/requests"
)

type CollectorData5u struct{}

func (c *CollectorData5u) Name() string {
	return "data5u.com"
}

// Collect returns proxys from www.data5u.com
func (c *CollectorData5u) Collect(cli *requests.Client) (result []*url.URL, err error) {
	const urlData5u = "http://www.data5u.com/free/index.shtml"
	resp, err := cli.NewRequest().Get(urlData5u)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.RawBody())
	if err != nil {
		return nil, err
	}

	doc.Find("ul > li > ul").Each(func(i int, s *goquery.Selection) {
		s = s.Children()
		ip := s.Eq(0).Text()
		port := s.Eq(1).Text()
		scheme := s.Eq(3).Text()
		proxyStr := fmt.Sprintf("%s://%s:%s", scheme, ip, port)
		proxy, err := url.Parse(proxyStr)
		if err == nil {
			result = append(result, proxy)
		}
	})

	return
}
