package gredis

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/redis/go-redis/v9"
	"log"
)

type Config struct {
	Address       string
	Username      string
	Password      string
	DB            int
	PoolSize      int
	MinIdle       int
	MaxIdle       int
	Cert, Key, Ca string
}

func NewClient(cfg *Config) *redis.Client {
	var tlsConfig *tls.Config
	if cfg.Cert != "" && cfg.Key != "" && cfg.Ca != "" {
		certificate, err := tls.X509KeyPair([]byte(cfg.Cert), []byte(cfg.Key))
		if err != nil {
			log.Fatal(err)
		}
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM([]byte(cfg.Ca)) {
			log.Fatal("failed to parse root certificate")
		}
		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{certificate},
			RootCAs:      pool,
		}
	}
	cli := redis.NewClient(&redis.Options{
		Addr:         cfg.Address,
		Username:     cfg.Username,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdle,
		MaxIdleConns: cfg.MaxIdle,
		TLSConfig:    tlsConfig,
	})
	if err := cli.Ping(context.Background()).Err(); err != nil {
		log.Fatal(err)
	}
	return cli
}
