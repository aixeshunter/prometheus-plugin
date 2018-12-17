FROM golang:1.11.2-alpine as builder

WORKDIR /go/src/github.com/aixeshunter/prometheus-plugin
ADD . .
RUN CGO_ENABLED=0 go build -o prometheus-plugin cmd/main.go

FROM busybox:1.29.3
COPY --from=builder /go/src/github.com/aixeshunter/prometheus-plugin/prometheus-plugin /
CMD ["/prometheus-plugin"]
