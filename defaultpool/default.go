package defaultpool

import (
	"github.com/wzshiming/proxypool"
	"github.com/wzshiming/proxypool/checker"
	"github.com/wzshiming/proxypool/collector"
)

var Default = proxypool.NewDispatcher()

func init() {
	Default.AddChecker(&checker.CheckerIP138{})
	Default.AddChecker(&checker.CheckerIPCN{})
	Default.AddCollector(&collector.CollectorData5u{})
	Default.AddCollector(&collector.CollectorGoubanjia{})
	Default.AddCollector(&collector.CollectorIP66{})
	Default.AddCollector(&collector.CollectorIP181{})
	Default.AddCollector(&collector.CollectorKuaidaili{})
	Default.AddCollector(&collector.CollectorXicidaili{})
}
