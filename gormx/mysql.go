package gormx

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/PirateDreamer/going/zlog.go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
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
	ZapLogger     *zap.Logger
	SlowThreshold time.Duration
}

func NewGormLogger() GormLogger {
	return GormLogger{
		ZapLogger:     zlog.Logger,
		SlowThreshold: 200 * time.Millisecond,
	}
}

func (l GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return GormLogger{
		ZapLogger:     l.ZapLogger,
		SlowThreshold: l.SlowThreshold,
	}
}

func (l GormLogger) Info(ctx context.Context, str string, args ...interface{}) {
	l.logger().Sugar().Debugf(str, args...)
}

func (l GormLogger) Warn(ctx context.Context, str string, args ...interface{}) {
	l.logger().Sugar().Warnf(str, args...)
}

func (l GormLogger) Error(ctx context.Context, str string, args ...interface{}) {
	l.logger().Sugar().Errorf(str, args...)
}

func (l GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {

	// 获取运行时间
	elapsed := time.Since(begin)
	// 获取 SQL 请求和返回条数
	sql, rows := fc()

	// 通用字段
	logFields := []zap.Field{
		zap.String("sql", sql),
		zap.String("time", fmt.Sprintf("%dµs", elapsed/time.Microsecond)),
		zap.Int64("rows", rows),
	}

	// Gorm 错误
	if err != nil {
		// 记录未找到的错误使用 warning 等级
		if errors.Is(err, gorm.ErrRecordNotFound) {
			l.logger().Warn("Database ErrRecordNotFound", logFields...)
		} else {
			// 其他错误使用 error 等级
			logFields = append(logFields, zap.Error(err))
			l.logger().Error("Database Error", logFields...)
		}
	}

	// 慢查询日志
	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		l.logger().Warn("Database Slow Log", logFields...)
	}

	// 记录所有 SQL 请求
	l.logger().Debug("Database Query", logFields...)
}

// logger 内用的辅助方法，确保 Zap 内置信息 Caller 的准确性（如 paginator/paginator.go:148）
func (l GormLogger) logger() *zap.Logger {

	// 跳过 gorm 内置的调用
	var (
		gormPackage    = filepath.Join("gorm.io", "gorm")
		zapgormPackage = filepath.Join("moul.io", "zapgorm2")
	)

	// 减去一次封装，以及一次在 logger 初始化里添加 zap.AddCallerSkip(1)
	clone := l.ZapLogger.WithOptions(zap.AddCallerSkip(-2))

	for i := 2; i < 15; i++ {
		_, file, _, ok := runtime.Caller(i)
		switch {
		case !ok:
		case strings.HasSuffix(file, "_test.go"):
		case strings.Contains(file, gormPackage):
		case strings.Contains(file, zapgormPackage):
		default:
			// 返回一个附带跳过行号的新的 zap logger
			return clone.WithOptions(zap.AddCallerSkip(i))
		}
	}
	return l.ZapLogger
}
