package general

import (
	"net"
	"net/http"
	"time"
)

// 科学计算法计算类型
const (
	Add DecimalType = iota
	Sub
	Multiply
	Divide
)

type DecimalType int

// 随机字符串数据源
const (
	Char = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// Time format常量
const (
	Layout      = "01/02 03:04:05PM '06 -0700"
	ANSIC       = "Mon Jan _2 15:04:05 2006"
	UnixDate    = "Mon Jan _2 15:04:05 MST 2006"
	RubyDate    = "Mon Jan 02 15:04:05 -0700 2006"
	RFC822      = "02 Jan 06 15:04 MST"
	RFC822Z     = "02 Jan 06 15:04 -0700"
	RFC850      = "Monday, 02-Jan-06 15:04:05 MST"
	RFC1123     = "Mon, 02 Jan 2006 15:04:05 MST"
	RFC1123Z    = "Mon, 02 Jan 2006 15:04:05 -0700"
	RFC3339     = "2006-01-02T15:04:05Z07:00"
	RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
	Kitchen     = "3:04PM"
	Stamp       = "Jan _2 15:04:05"
	StampMilli  = "Jan _2 15:04:05.000"
	StampMicro  = "Jan _2 15:04:05.000000"
	StampNano   = "Jan _2 15:04:05.000000000"
	DateTime    = "2006-01-02 15:04:05"
	DateOnly    = "2006-01-02"
	TimeOnly    = "15:04:05"
)

// Ordered 可排序类型
type Ordered interface {
	Number | ~string
}

// Number 数字类型
type Number interface {
	Integer | Float
}

// SignedNumber 有符号数字类型
type SignedNumber interface {
	Signed | Float
}

// Integer 整数类型
type Integer interface {
	Signed | Unsigned
}

// Signed 有符号整数类型
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Unsigned 无符号整数类型
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Float 浮点类型
type Float interface {
	~float32 | ~float64
}

var (
	// Http请求配置
	httpClient = &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   60 * time.Second,
				KeepAlive: 60 * time.Second,
			}).DialContext,
			DisableKeepAlives:     true,
			TLSHandshakeTimeout:   60 * time.Second,
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   100,
			MaxConnsPerHost:       100,
			IdleConnTimeout:       60 * time.Second,
			ResponseHeaderTimeout: 60 * time.Second,
			ExpectContinueTimeout: 60 * time.Second,
		},
		Timeout: 60 * time.Second,
	}
)
