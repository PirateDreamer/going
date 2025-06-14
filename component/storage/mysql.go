package storage

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/PirateDreamer/going/zlog"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func InitMysql(config *viper.Viper) *gorm.DB {
	schema.RegisterSerializer("timestamp", TimestampSerializer{})

	host := config.GetString("mysql.host")
	if host == "" {
		return nil
	}

	port := config.GetInt("mysql.port")
	user := config.GetString("mysql.user")
	pass := config.GetString("mysql.pass")
	database := config.GetString("mysql.database")
	dsn := fmt.Sprint(user, ":", pass, "@tcp(", host, ":", port, ")/", database, "?charset=utf8mb4&parseTime=True&loc=Local")

	gormConfig := &gorm.Config{}
	if config.GetBool("mysql.enable_log") {
		gormConfig.Logger = NewGormLogger()
	}

	mysqlDB, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		panic(err)
	}

	sqlDB, err := mysqlDB.DB()
	if err != nil {
		panic(err)
	}

	if idleConns := config.GetInt("mysql.idle_conns"); idleConns > 0 {
		sqlDB.SetMaxIdleConns(idleConns)
	}

	if maxOpenConns := config.GetInt("mysql.max_open_conns"); maxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(maxOpenConns)
	}

	if maxLiftSeconds := config.GetInt64("mysql.max_life_seconds"); maxLiftSeconds > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(maxLiftSeconds * int64(time.Second)))
	}

	return mysqlDB
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

// 用于序列化 datetime,int,var -> int64/int32
type TimestampSerializer struct{}

// Scan 实现 sql.Scanner 接口，用于从数据库读取
func (TimestampSerializer) Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) (err error) {
	// model 定义的类型必须是 int64, int32, int
	kind := field.IndirectFieldType.Kind()
	if kind != reflect.Int64 && kind != reflect.Int32 && kind != reflect.Int {
		return fmt.Errorf("invalid field type %#v for Serializer, only int, int32, int64 supported", kind)
	}

	if dbValue == nil {
		return
	}

	var result int64

	// 判断 db 中的值类型，解析成时间戳
	switch v := dbValue.(type) {
	case time.Time:
		result = v.Local().UnixMilli()
	case string:
		var t time.Time
		if t, err = time.Parse("2006-01-02 15:04:05", v); err != nil {
			return
		}
		result = t.Local().UnixMilli()
	case int64, int32, int:
		result = reflect.ValueOf(v).Int()
	default:
		return fmt.Errorf("failed to unmarshal timestamp value: %#v, only DATETIME INT VARCHAR supported", dbValue)
	}

	if result < 0 {
		result = 0
	}

	// 转换成 model 定义的类型
	if kind == reflect.Int32 {
		field.Set(ctx, dst, int32(result/1000))
	} else {
		field.Set(ctx, dst, result)
	}

	return
}

// Value 实现 driver.Valuer 接口，用于写入数据库
func (TimestampSerializer) Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (value interface{}, err error) {
	rv := reflect.ValueOf(fieldValue)
	v := reflect.Indirect(rv).Int()

	if rv.IsZero() {
		value = time.Now()
		return
	}

	switch len(strconv.FormatInt(v, 10)) {
	case 10:
		value = time.Unix(v, 0)
	case 13:
		value = time.UnixMilli(v)
	default:
		err = fmt.Errorf("invalid timestamp value: %d", v)
	}

	return
}
