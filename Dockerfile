FROM docker.io/golang
WORKDIR  /root
COPY . .
RUN go build
CMD ["./sslcert-exporter  -domains=www.xx.com,www.yyy.com"]