package utils

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"
	"time"
)

func StringToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func IntToString(i int) string {
	return strconv.Itoa(i)
}

func StringToFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func FloatToString(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}

func RoundFloat(val float64, precision int) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func GetCurrentTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func ParseJSON(data string) map[string]interface{} {
	var result map[string]interface{}
	json.Unmarshal([]byte(data), &result)
	return result
}

func ToJSON(data interface{}) string {
	bytes, _ := json.Marshal(data)
	return string(bytes)
}

func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func Sleep(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
