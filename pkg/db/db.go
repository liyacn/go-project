package db

import (
	"context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glog "gorm.io/gorm/logger"
	"log"
	"project/pkg/logger"
	"time"
)

type Config struct {
	Address     string
	Username    string
	Password    string
	Database    string
	MaxOpen     int
	MaxIdle     int
	MaxLifeTime time.Duration
	MaxIdleTime time.Duration
	TraceLog    bool
}

func NewMysql(cfg *Config) *gorm.DB {
	dsn := cfg.Username + ":" + cfg.Password + "@tcp(" + cfg.Address + ")/" + cfg.Database +
		"?charset=utf8mb4&parseTime=true&loc=Local"
	return newDB(mysql.Open(dsn), cfg)
}

//func NewPostgres(cfg *Config) *gorm.DB {
//	host, port, _ := net.SplitHostPort(cfg.Address)
//	dsn := "host=" + host + " port=" + port + " user=" + cfg.Username + " password=" + cfg.Password +
//		" dbname=" + cfg.Database + " sslmode=disable TimeZone=Asia/Shanghai"
//	return newDB(postgres.Open(dsn), cfg)
//}
//
//func NewSqlserver(cfg *Config) *gorm.DB {
//	dsn := "sqlserver://" + cfg.Username + ":" + cfg.Password +
//		"@" + cfg.Address + "?database=" + cfg.Database
//	return newDB(sqlserver.Open(dsn), cfg)
//}

func newDB(dial gorm.Dialector, cfg *Config) *gorm.DB {
	opt := &gorm.Config{Logger: glog.Discard.LogMode(glog.Silent)}
	if cfg.TraceLog {
		opt.Logger = &gormLog{opt.Logger}
	}
	orm, err := gorm.Open(dial, opt)
	if err != nil {
		log.Fatal(err)
	}
	sqlDB, _ := orm.DB()
	if cfg.MaxOpen > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpen)
	}
	if cfg.MaxIdle > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdle)
	}
	if cfg.MaxLifeTime > 0 {
		sqlDB.SetConnMaxLifetime(cfg.MaxLifeTime)
	}
	if cfg.MaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(cfg.MaxIdleTime)
	}
	return orm
}

func Close(orm *gorm.DB) error {
	sqlDB, _ := orm.DB()
	return sqlDB.Close()
}

type gormLog struct {
	glog.Interface
}

func (*gormLog) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rows int64), err error) {
	sql, rows := fc()
	logger.FromContext(ctx).Debug("gorm", sql, rows, err, time.Since(begin))
}
