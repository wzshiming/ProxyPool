package collector

import (
	"fmt"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/wzshiming/requests"
)

type CollectorKuaidaili struct{}

func (c *CollectorKuaidaili) Name() string {
	return "kuaidaili.com"
}

// Collect returns proxys from www.kuaidaili.com
func (c *CollectorKuaidaili) Collect(cli *requests.Client) (result []*url.URL, err error) {
	const urlKuaidaili = "https://www.kuaidaili.com/free/inha/"
	resp, err := cli.NewRequest().Get(urlKuaidaili)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.RawBody())
	if err != nil {
		return nil, err
	}

	doc.Find("tbody > tr").Each(func(i int, s *goquery.Selection) {
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
