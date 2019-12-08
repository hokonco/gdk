package cast

import "strconv"

// StringToFloat32 ...
func StringToFloat32(arg string) (float32, error) {
	f64, e := strconv.ParseFloat(arg, 32)
	return float32(f64), e
}

// StringToFloat64 ...
func StringToFloat64(arg string) (float64, error) {
	return strconv.ParseFloat(arg, 64)
}

// StringToInterface ...
func StringToInterface(arg string) interface{} { return arg }

// StringToInt8 ...
func StringToInt8(arg string) (int8, error) {
	i64, e := strconv.ParseInt(arg, 10, 8)
	return int8(i64), e
}

// StringToInt64 ...
func StringToInt64(arg string) (int64, error) {
	return strconv.ParseInt(arg, 10, 64)
}

// StringToUint8 ...
func StringToUint8(arg string) (uint8, error) {
	u64, e := strconv.ParseUint(arg, 10, 8)
	return uint8(u64), e
}

// StringToUint64 ...
func StringToUint64(arg string) (uint64, error) {
	return strconv.ParseUint(arg, 10, 64)
}
