FROM golang:alpine AS builder

WORKDIR $GOPATH/src/habr-searcher

COPY . .

RUN go build -o /go/bin/habr-searcher cmd/main.go



FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /go/bin/habr-searcher /go/bin/habr-searcher

ENTRYPOINT ["/go/bin/habr-searcher"]

