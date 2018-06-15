package collector

import (
	"fmt"
	"net/url"
	"regexp"

	"github.com/wzshiming/requests"
)

type CollectorIP66 struct{}

func (c *CollectorIP66) Name() string {
	return "66ip.cn"
}

// Collect returns proxys from www.66ip.cn
func (c *CollectorIP66) Collect(cli *requests.Client) (result []*url.URL, err error) {
	const urlIP66 = "http://www.66ip.cn/mo.php?tqsl=100"
	var ipReg = regexp.MustCompile(`\s\S+<br`)
	resp, err := cli.NewRequest().Get(urlIP66)
	if err != nil {
		return nil, err
	}

	for _, v := range ipReg.FindAllString(string(resp.Body()), -1) {
		if v == "" {
			continue
		}
		v = v[1 : len(v)-3]
		v = fmt.Sprintf("http://%s", v)
		proxy, err := url.Parse(v)
		if err == nil {
			result = append(result, proxy)
		}
	}
	return
}
