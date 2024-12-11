package gormx

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/PirateDreamer/going/zlog"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Mysql *gorm.DB

func InitMysql() {
	host := viper.GetString("mysql.host")
	if host == "" {
		return
	}

	port := viper.GetInt("mysql.port")
	user := viper.GetString("mysql.user")
	pass := viper.GetString("mysql.pass")
	database := viper.GetString("mysql.database")
	dsn := fmt.Sprint(user, ":", pass, "@tcp(", host, ":", port, ")/", database, "?charset=utf8mb4&parseTime=True&loc=Local")

	gormConfig := &gorm.Config{}
	if viper.GetBool("mysql.enable_log") {
		gormConfig.Logger = NewGormLogger()
	}

	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	if idleConns := viper.GetInt("mysql.idle_conns"); idleConns > 0 {
		sqlDB.SetMaxIdleConns(idleConns)
	}

	if maxOpenConns := viper.GetInt("mysql.max_open_conns"); maxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(maxOpenConns)
	}

	if maxLiftSeconds := viper.GetInt64("mysql.max_life_seconds"); maxLiftSeconds > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(maxLiftSeconds * int64(time.Second)))
	}

	Mysql = db

}

type GormLogger struct {
	SlowThreshold time.Duration
}

func NewGormLogger() GormLogger {
	return GormLogger{
		SlowThreshold: 200 * time.Millisecond,
	}
}

func (l GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return GormLogger{
		SlowThreshold: l.SlowThreshold,
	}
}

func (l GormLogger) Info(ctx context.Context, str string, args ...interface{}) {
	zlog.LogInfo(ctx, str, args...)
}

func (l GormLogger) Warn(ctx context.Context, str string, args ...interface{}) {
	zlog.LogWarn(ctx, str, args...)
}

func (l GormLogger) Error(ctx context.Context, str string, args ...interface{}) {
	zlog.LogError(ctx, str, args...)
}

func (l GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {

	// 获取运行时间
	elapsed := time.Since(begin)
	// 获取 SQL 请求和返回条数
	sql, rows := fc()

	sqlExecInfo := fmt.Sprintf("SQL: %s Time: %dµs Rows: %d ", sql, elapsed/time.Microsecond, rows)

	// Gorm 错误
	if err != nil {
		// 记录未找到的错误使用 warning 等级
		if errors.Is(err, gorm.ErrRecordNotFound) {
			zlog.LogWarn(ctx, sqlExecInfo+"WARNING: Database ErrRecordNotFound")
		} else {
			// 其他错误使用 error 等级
			zlog.LogError(ctx, sqlExecInfo+"ERROR: Database Error %s", err.Error())
		}
		return
	}

	// 慢查询日志
	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		zlog.LogWarn(ctx, sqlExecInfo+"WARNING: Database Slow Log")
	}

	// 记录所有 SQL 请求
	zlog.LogInfo(ctx, sqlExecInfo+"INFO: Database Query")
}
