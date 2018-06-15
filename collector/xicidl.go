package collector

import (
	"fmt"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/wzshiming/requests"
)

type CollectorXicidaili struct{}

func (c *CollectorXicidaili) Name() string {
	return "xicidaili.com"
}

// Collect returns proxys from www.xicidaili.com
func (c *CollectorXicidaili) Collect(cli *requests.Client) (result []*url.URL, err error) {
	const urlXicidaili = "http://www.xicidaili.com/nn/"
	resp, err := cli.NewRequest().Get(urlXicidaili)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.RawBody())
	if err != nil {
		return nil, err
	}

	doc.Find("tbody > tr").Each(func(i int, s *goquery.Selection) {
		s = s.Children()
		ip := s.Eq(1).Text()
		port := s.Eq(2).Text()
		scheme := s.Eq(5).Text()
		proxyStr := fmt.Sprintf("%s://%s:%s", scheme, ip, port)
		proxy, err := url.Parse(proxyStr)
		if err == nil {
			result = append(result, proxy)
		}
	})
	return
}
