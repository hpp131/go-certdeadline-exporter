package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	//"github.com/prometheus/client_golang/api/prometheus/"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strings"
	"time"
)

type sslStateCollector struct {
	expireDay *prometheus.Desc
}

func NewSslStateCollector() prometheus.Collector {
	return &sslStateCollector{
		expireDay: prometheus.NewDesc(
			"deadline",
			"the deadline of specific domain",
			[]string{"domain_name"},
			nil)}
}

func (s *sslStateCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- s.expireDay
}

/*
use map to store domain: deadline
*/
func (s *sslStateCollector) Collect(ch chan<- prometheus.Metric) {
	for key, value := range checkHttps() {
		ch <- prometheus.MustNewConstMetric(s.expireDay, prometheus.GaugeValue, value, key)
	}
	//ch <- prometheus.MustNewConstMetric(s.expireDay, prometheus.GaugeValue, checkHttps(), "www.baidu.com")
}

func checkHttps() (t map[string]float64) {
	/*
		get domain string and split into []string
	*/
	domainString := *flag.String("domains", "", "a set of domain separated by comma")
	domainSlice := strings.Split(domainString, ",")
	var deadlineMap map[string]float64
	for _, value := range domainSlice {
		checkUrl := "https://" + value
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		resp, err := client.Get(checkUrl)
		fmt.Printf("the panic is happend")
		defer resp.Body.Close() ///// panic
		if err != nil {
			fmt.Errorf("requeset is failed", err)
		}
		certinfo := resp.TLS.PeerCertificates[0]
		fmt.Println("过期时间", certinfo.NotAfter)
		//deadlineSlice = append(deadlineSlice, certinfo.NotAfter.Sub(time.Now()).Hours()/24)
		deadlineMap[value] = certinfo.NotAfter.Sub(time.Now()).Hours() / 24
	}

	return deadlineMap

}

func main() {
	reg := prometheus.NewRegistry()
	reg.Register(NewSslStateCollector())
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	http.ListenAndServe(":9999", nil)
}
