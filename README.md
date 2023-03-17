### usage  
- Start by commandline  
Between numbers of domain,using comma to split. Start method are as follows:    
go build  sslcert-exporter  && ./sslcert-exporter  -domains=www.xxx.com,www.xx.com  

- Start by container  
Change Dockerfile "CMD" field, replace target domains to which you care

### Implement
- The custom interface named "Collector" contains 2 methods: describe and collect
- This implement don't need usinig roundrobin to get metrics data periodically. When we interview the metircs api(e.g http://localhost:9999/metrics),the collect method is called. And in collect method, the  "checkHttps" method is called automatically.