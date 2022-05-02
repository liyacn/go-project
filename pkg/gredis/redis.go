package gredis

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
	"log"
	"net"
	"project/pkg/logger"
)

type Config struct {
	Address    string
	Username   string
	Password   string
	DB         int
	PoolSize   int
	MinIdle    int
	MaxIdle    int
	ServerName string
	Cert       string
	Key        string
	Ca         string
	TraceLog   bool
}

func NewClient(cfg *Config) *redis.Client {
	cli := redis.NewClient(&redis.Options{
		Addr:                  cfg.Address,     // default: "localhost:6379"
		Username:              cfg.Username,    //
		Password:              cfg.Password,    //
		DB:                    cfg.DB,          //
		MaxRetries:            0,               // default: 3
		MinRetryBackoff:       0,               // default: 8 milliseconds
		MaxRetryBackoff:       0,               // default: 512 milliseconds
		DialTimeout:           0,               // default: 5 seconds
		DialerRetries:         0,               // default: 5
		DialerRetryTimeout:    0,               // default: 100 milliseconds
		ReadTimeout:           0,               // default: 3 seconds
		WriteTimeout:          0,               // default: same as ReadTimeout
		ReadBufferSize:        0,               // default: 32KiB
		WriteBufferSize:       0,               // default: 32KiB
		PoolSize:              cfg.PoolSize,    // default: 10 * runtime.GOMAXPROCS(0)
		PoolTimeout:           0,               // default: ReadTimeout + 1 second
		MinIdleConns:          cfg.MinIdle,     //
		MaxIdleConns:          cfg.MaxIdle,     //
		MaxActiveConns:        0,               //
		ConnMaxIdleTime:       0,               // default: 30 minutes
		ConnMaxLifetime:       0,               //
		TLSConfig:             cfg.tlsConfig(), //
		FailingTimeoutSeconds: 0,               // default: 15
		MaintNotificationsConfig: &maintnotifications.Config{
			Mode: maintnotifications.ModeDisabled, // default: ModeAuto
		},
	})
	if err := cli.Ping(context.Background()).Err(); err != nil {
		log.Fatal(err)
	}
	if cfg.TraceLog {
		cli.AddHook(debugHook{})
	}
	return cli
}

type debugHook struct{}

func (debugHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}
func (debugHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		logger.FromContext(ctx).Debug("redis", cmd.String())
		return next(ctx, cmd)
	}
}
func (debugHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		l := logger.FromContext(ctx)
		for _, cmd := range cmds {
			l.Debug("redis", cmd.String())
		}
		return next(ctx, cmds)
	}
}

func (cfg *Config) tlsConfig() *tls.Config {
	if cfg.ServerName == "" && (cfg.Cert == "" || cfg.Key == "") && cfg.Ca == "" {
		return nil
	}
	tlsConfig := &tls.Config{ServerName: cfg.ServerName}
	if cfg.Cert != "" && cfg.Key != "" {
		certificate, err := tls.X509KeyPair([]byte(cfg.Cert), []byte(cfg.Key))
		if err != nil {
			log.Fatal(err)
		}
		tlsConfig.Certificates = []tls.Certificate{certificate}
	}
	if cfg.Ca != "" {
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM([]byte(cfg.Ca)) {
			log.Fatal("failed to parse root certificate")
		}
		tlsConfig.RootCAs = pool
	}
	return tlsConfig
}
