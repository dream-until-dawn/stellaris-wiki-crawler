package utils

import (
	"reflect"
	"strconv"
	"time"
)

func AbsFloat(f float64) float64 {
	if f < 0 {
		return -f
	} else {
		return f
	}
}

func StringToFloat64(s string) float64 {
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return val
}

func StringToInt64(s string) int64 {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return val
}

func TimestampToString(timestamp int64) string {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	return time.Unix(timestamp, 0).In(loc).Format("2006-01-02")
}

// GetYesterdayZeroMillisCST 返回东八区昨天0点的毫秒级时间戳
func GetYesterdayZeroMillisCST() int64 {
	// 东八区（UTC+8）
	loc := time.FixedZone("CST", 8*3600)

	// 当前时间（东八区）
	now := time.Now().In(loc)

	// 今天 0 点
	todayZero := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		loc,
	)

	// 昨天 0 点
	yesterdayZero := todayZero.AddDate(0, 0, -1)

	// 毫秒级时间戳
	return yesterdayZero.UnixNano() / int64(time.Millisecond)
}

func FlattenToFloatMap(v any) map[string]float64 {
	result := make(map[string]float64)
	flatten(reflect.ValueOf(v), "", result)
	return result
}

func flatten(val reflect.Value, prefix string, out map[string]float64) {
	// 处理指针
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return
		}
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Struct:
		t := val.Type()
		for i := 0; i < val.NumField(); i++ {
			fieldVal := val.Field(i)
			fieldType := t.Field(i)

			// 非导出字段跳过
			if fieldType.PkgPath != "" {
				continue
			}

			key := fieldType.Name
			if prefix != "" {
				key = prefix + "." + key
			}

			flatten(fieldVal, key, out)
		}

	case reflect.Float32, reflect.Float64:
		out[prefix] = val.Convert(reflect.TypeOf(float64(0))).Float()

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		out[prefix] = float64(val.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		out[prefix] = float64(val.Uint())

	// 其他类型（string、bool、slice、map 等）自动忽略
	default:
		return
	}
}
