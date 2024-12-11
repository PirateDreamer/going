package gormx

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"gorm.io/gorm/schema"
)

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
