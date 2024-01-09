package general

import "time"

/*
获取当前时间对象（默认返回当前时区）
now：时间对象
*/
func Now() (now time.Time) {
	now = time.Now().Local()
	return
}

/*
获取当前秒级时间戳
now：当前秒级时间戳
*/
func NowUnix() (now int64) {
	now = Now().Unix()
	return
}

/*
获取当前毫秒级时间戳
now：当前毫秒级时间戳
*/
func NowMilli() (now int64) {
	now = Now().UnixMilli()
	return
}

/*
获取当前微秒级时间戳
now：当前微秒级时间戳
*/
func NowMicro() (now int64) {
	now = Now().UnixMicro()
	return
}

/*
获取当前纳秒级时间戳
now：当前纳秒级时间戳
*/
func NowNano() (now int64) {
	now = Now().UnixNano()
	return
}

/*
获取偏移时间对象
offset：偏移量
offsetTime：偏移后时间对象
*/
func NowAdd(offset time.Duration) (offsetTime time.Time) {
	offsetTime = Now().Add(offset)
	return
}

/*
获取偏移天数对象
year：偏移年数
month：偏移月数
day：偏移天数
offsetTime：偏移后时间对象
*/
func NowAddDate(year, month, day int) (offsetTime time.Time) {
	offsetTime = Now().AddDate(year, month, day)
	return
}

/*
获取指定天数0时时间对象
year：偏移年数
month：偏移月数
day：偏移天数
startTime：偏移后天数0时时间对象
*/
func DayStart(year, month, day int) (startTime time.Time) {
	target := Now().AddDate(year, month, day)
	startTime = time.Date(target.Year(), target.Month(), target.Day(), 0, 0, 0, 0, time.Local)
	return
}

/*
获取指定天数最后一毫秒时间对象
year：偏移年数
month：偏移月数
day：偏移天数
EndTime：偏移后天数最后一毫秒时间对象
*/
func DayEnd(year, month, day int) (EndTime time.Time) {
	target := Now().AddDate(year, month, day)
	EndTime = time.Date(target.Year(), target.Month(), target.Day(), 23, 59, 59, 999, time.Local)
	return
}

/*
时间对象转字符串
target：时间对象
format：格式
result：转换结果
*/
func TimeToString(target time.Time, format string) (result string) {
	result = target.Format(format)
	return
}

/*
字符串转时间对象
timeString：时间字符串
format：格式
result：转换结果
*/
func StringToTime(timeString string, format string) (result time.Time) {
	if t, err := time.ParseInLocation(format, timeString, time.Local); err == nil {
		result = t
	}
	return
}
