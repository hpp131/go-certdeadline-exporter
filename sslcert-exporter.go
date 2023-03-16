package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strings"
	"time"
)

type sslStateCollector struct {
	expireDay []*prometheus.Desc
	domains   []string
}

func NewSslStateCollector(s []string) prometheus.Collector {
	var descSlice []*prometheus.Desc
	for i := 0; i < len(s); i++ {
		descSlice = append(descSlice, prometheus.NewDesc(
			"deadline",
			"the deadline of specific domain",
			[]string{"domain_name"},
			nil))
	}
	return &sslStateCollector{
		expireDay: descSlice,
		domains:   s,
	}
}

func (s *sslStateCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, value := range s.expireDay {
		ch <- value
	}
}

/*
use map to store "domain" "deadline"
*/
func (s *sslStateCollector) Collect(ch chan<- prometheus.Metric) {
	for key, value := range checkHttps(s.domains) {
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				"deadline",
				"the deadline of specific domain",
				[]string{"domain_name"},
				nil),
			prometheus.GaugeValue, value, key)
	}
}

func checkHttps(s []string) (t map[string]float64) {
	/*
		get domain string and split into []string
	*/
	deadlineMap := make(map[string]float64)
	for _, domain := range s {
		checkUrl := "https://" + domain
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		resp, err := client.Get(checkUrl)
		fmt.Printf("the panic is happend")
		resp.Body.Close() ///// panic
		if err != nil {
			fmt.Errorf("requeset is failed", err)
		}
		certinfo := resp.TLS.PeerCertificates[0]
		fmt.Println("过期时间", certinfo.NotAfter)
		deadlineMap[domain] = certinfo.NotAfter.Sub(time.Now()).Hours() / 24
	}
	return deadlineMap
}

func main() {
	domainString := flag.String("domains", "", "a set of domain separated by comma")
	flag.Parse()
	domainSlice := strings.Split(*domainString, ",")
	// use custom registry without build-in go-metrics
	reg := prometheus.NewRegistry()
	reg.Register(NewSslStateCollector(domainSlice))
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	http.ListenAndServe(":9999", nil)
}
