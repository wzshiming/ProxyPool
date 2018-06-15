package proxypool

import (
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/wzshiming/fork"

	"github.com/wzshiming/requests"
	"github.com/wzshiming/task"
	ffmt "gopkg.in/ffmt.v1"
)

const invalid = "invalid"

// Dispatcher is dispatcher of the proxy pool
type Dispatcher struct {
	task       *task.Task
	cliCollect *requests.Client
	cliCheck   *requests.Client
	collectors []Collector
	checkers   map[string][]Checker
	ch         chan *url.URL
	//pools      *pools
	checkFork *fork.Fork
	Storage
}

func NewDispatcher() *Dispatcher {
	d := &Dispatcher{
		task:       task.NewTask(100),
		cliCollect: requests.NewClient().SetTimeout(time.Second * 10),
		cliCheck:   requests.NewClient().SetTimeout(time.Second * 5),
		checkers:   map[string][]Checker{},
		ch:         make(chan *url.URL, 1000),
		checkFork:  fork.NewFork(10),
	}
	d.cliCollect.SetProxyFunc(d.ProxyFunc)
	return d
}

func (d *Dispatcher) ProxyFunc(r *http.Request) (*url.URL, error) {
	u, err := d.Storage.GetRandom(r.URL.Scheme)
	if err != nil {
		return nil, err
	}
	if u == "" {
		return nil, nil
	}
	d.Storage.IncUsed(r.URL.Scheme, u)
	return url.Parse(u)
}

func (d *Dispatcher) MustProxyFunc(r *http.Request) (*url.URL, error) {
	for {
		u, _ := d.ProxyFunc(r)
		if u != nil {
			d.ch <- u
			return u, nil
		}
		time.Sleep(time.Second)
	}
}

func (d *Dispatcher) SetStorage(s Storage) {
	d.Storage = s
}

func (d *Dispatcher) AddCollector(c Collector) {
	d.collectors = append(d.collectors, c)
}

func (d *Dispatcher) AddChecker(c Checker) {
	scheme := c.Scheme()
	d.checkers[scheme] = append(d.checkers[scheme], c)
}

func (d *Dispatcher) Run() {
	go d.runCheckers()
	for k, _ := range d.checkers {
		d.runRecheckers(k)
	}

	for _, v := range d.collectors {
		d.runCollectors(v)
	}
	d.task.Join()
}

func (d *Dispatcher) runCheckers() {
	index := 0
	for v := range d.ch {
		index++
		checks := d.checkers[v.Scheme]
		check := checks[index%len(checks)]
		d.check(check, v)
	}
}

func (d *Dispatcher) check(check Checker, v *url.URL) {
	d.checkFork.Push(func() {
		_, err := check.Check(d.cliCheck, v)
		if err != nil {
			d.Storage.IncFailure(v.Scheme, v.String())
			return
		}
		d.Storage.IncCheck(v.Scheme, v.String())
	})
}

func (d *Dispatcher) runRecheckers(pre string) {
	interval := time.Second * 30
	d.task.AddPeriodic(task.PeriodicIntervalCount(time.Now().Add(1-interval), interval, -1), func() {
		tmp, _ := d.Storage.GetAll(pre)
		for _, v := range tmp {
			u, _ := url.Parse(v)
			if u != nil {
				d.ch <- u
			}
		}
	})
}

func (d *Dispatcher) runCollectors(coll Collector) {
	interval := time.Second * 60
	d.task.AddPeriodic(task.PeriodicIntervalCount(time.Now().Add(1-interval), interval, -1), func() {
		proxys, err := coll.Collect(d.cliCollect)
		if err != nil {
			ffmt.Mark(coll.Name(), err)
			return
		}
		for _, v := range proxys {
			if ok, _ := d.Storage.IsExist(v.Scheme, v.String()); ok {
				continue
			}
			d.ch <- v
		}
	})
}

// Checker is to access the third party url to query the IP to check whether the proxy is valid
type Checker interface {
	// Name the name of Checker
	Name() string
	// Scheme the scheme of Checker
	Scheme() string
	// Check is check the proxy
	Check(*requests.Client, *url.URL) (net.IP, error)
}

// Collector is collection proxy
type Collector interface {
	// Name the name of Collector
	Name() string
	// Collect is Collect the proxy
	Collect(*requests.Client) ([]*url.URL, error)
}

type Storage interface {
	IsExist(prefix, raw string) (bool, error)
	IncUsed(prefix, raw string) error
	IncCheck(prefix, raw string) error
	IncFailure(prefix, raw string) error
	GetRandom(prefix string) (string, error)
	GetAll(prefix string) ([]string, error)
}
