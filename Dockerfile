FROM golang:1.25-alpine AS builder
ARG srv
WORKDIR /build
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct \
    GOSUMDB=sum.golang.google.cn
COPY . .
RUN apk add --no-cache tzdata && \
    go mod tidy && \
    go build -tags sonic -o app ./${srv}

FROM scratch
ARG srv
WORKDIR /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=builder /build/app .
COPY ${srv}/conf.yaml .
COPY ${srv}/docs docs
ENTRYPOINT ["./app"]
