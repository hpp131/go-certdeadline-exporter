package main

import (
	"crypto/tls"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

type sslStateCollector struct {
	expireDay *prometheus.Desc
}

func NewSslStateCollector(domain string) prometheus.Collector {
	return &sslStateCollector{
		expireDay: prometheus.NewDesc(
			"deadline",
			"the deadline of specific domain",
			nil,
			prometheus.Labels{"domain_name": domain})}
}

func (s *sslStateCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- s.expireDay
}

func (s *sslStateCollector) Collect(ch chan<- prometheus.Metric) {

	ch <- prometheus.MustNewConstMetric(s.expireDay, prometheus.GaugeValue, checkHttps())
}

func checkHttps() (t float64) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	checkUrl := "https://www.baidu.com"
	resp, err := client.Get(checkUrl)
	fmt.Printf("the panic is happend")
	defer resp.Body.Close() ///// panic
	if err != nil {
		fmt.Errorf("requeset is failed", err)
	}
	certinfo := resp.TLS.PeerCertificates[0]
	fmt.Println("过期时间", certinfo.NotAfter)
	return certinfo.NotAfter.Sub(time.Now()).Hours() / 24

}

func main() {
	reg := prometheus.NewRegistry()
	reg.Register(NewSslStateCollector("www.baidu.com"))
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	http.ListenAndServe(":9999", nil)
}
