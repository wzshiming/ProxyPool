package proxypool

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"sync"

	ffmt "gopkg.in/ffmt.v1"
)

var ErrNoProxy = fmt.Errorf("There's no proxy IP")

type pools struct {
	Pools   map[string]*pool
	Inuse   *pool
	Discard *pool
}

func newPools() *pools {
	return &pools{
		Pools:   map[string]*pool{},
		Inuse:   &pool{},
		Discard: &pool{},
	}
}

func (p *pools) Init() {
	ps := []string{}

	p.Inuse.pool.Range(func(key, value interface{}) bool {
		r, _ := key.(string)
		if r != "" {
			ps = append(ps, r)
		}
		return true
	})
	p.Inuse = &pool{}
	for _, v := range ps {
		u, _ := url.Parse(v)
		p.Put(u)
	}

}

func (p *pools) ExistsDiscard(u *url.URL) bool {
	return p.Discard.exists(u.String())
}

func (p *pools) PutDiscard(u *url.URL) {
	raw := u.String()
	p.Discard.put(raw)
	if p.Inuse.exists(raw) {
		p.Inuse.delete(raw)
	}
}

func (p *pools) Put(u *url.URL) {
	if p.Pools[u.Scheme] == nil {
		p.Pools[u.Scheme] = &pool{}
	}
	raw := u.String()
	p.Pools[u.Scheme].put(raw)
	if p.Inuse.exists(raw) {
		p.Inuse.delete(raw)
	} else {
		ffmt.Mark("[COLLECT PROXY] ", u)
	}
}

func (p *pools) ProxyFunc(u *http.Request) (*url.URL, error) {
	pool := p.Pools[u.URL.Scheme]
	if pool == nil {
		return nil, nil
	}

	proxy := pool.get()
	if proxy == "" {
		return nil, nil
	}
	p.Inuse.put(proxy)
	return url.Parse(proxy)

}

type pool struct {
	pool sync.Map
}

func (p *pool) MarshalJSON() ([]byte, error) {
	ps := []string{}
	p.pool.Range(func(key, value interface{}) bool {
		r, _ := key.(string)
		if r != "" {
			ps = append(ps, r)
		}
		return true
	})
	sort.Strings(ps)
	return json.Marshal(ps)
}

func (p *pool) UnmarshalJSON(d []byte) error {
	ps := []string{}
	err := json.Unmarshal(d, &ps)
	if err != nil {
		return err
	}
	for _, v := range ps {
		p.put(v)
	}
	return nil
}

func (p *pool) exists(raw string) bool {
	_, ok := p.pool.Load(raw)
	return ok
}

func (p *pool) delete(raw string) {
	p.pool.Delete(raw)
}

func (p *pool) put(raw string) {
	p.pool.Store(raw, 1)
}

func (p *pool) get() string {
	r := ""
	p.pool.Range(func(key, value interface{}) bool {
		r, _ = key.(string)
		return false
	})
	if r != "" {
		p.pool.Delete(r)
	}
	return r
}
