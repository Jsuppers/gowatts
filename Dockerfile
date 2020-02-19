FROM golang@sha256:cee6f4b901543e8e3f20da3a4f7caac6ea643fd5a46201c3c2387183a332d989 as builder

WORKDIR /go/gowatts
COPY . /go/gowatts
RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates && apk add make && apk add git
RUN cat Makefile && go mod download && make build

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/gowatts/resources /resources
COPY --from=builder /go/gowatts/bin/gowatts service

ENV GIN_MODE release 

CMD ["./service"]
