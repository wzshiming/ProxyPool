package collector

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/wzshiming/requests"
)

type CollectorIP181 struct{}

func (c *CollectorIP181) Name() string {
	return "ip181.com"
}

// Collect returns proxys from www.ip181.com
func (c *CollectorIP181) Collect(cli *requests.Client) (result []*url.URL, err error) {
	const urlIP181 = "http://www.ip181.com/"
	resp, err := cli.NewRequest().Get(urlIP181)
	if err != nil {
		return nil, err
	}

	mod := struct {
		ErrorCode string `json:"ERRORCODE"`
		Result    []struct {
			Pos  string `json:"position"`
			Port string `json:"port"`
			IP   string `json:"ip"`
		} `json:"RESULT"`
	}{}

	err = json.Unmarshal(resp.Body(), &mod)
	if err != nil {
		return nil, err
	}

	for _, v := range mod.Result {
		proxyStr := fmt.Sprintf("http://%s:%s", v.IP, v.Port)
		proxy, err := url.Parse(proxyStr)
		if err == nil {
			result = append(result, proxy)
		}
	}
	return
}
