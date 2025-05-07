package common

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"gorm.io/datatypes"
)

const (
	ERR_PARSE_VALUE_ENV = "cannot parse value of %v env"
)

var (
	FormatErr = func(prefix string, err error) error {
		return ErrorWrapper(fmt.Sprintf(ERR_PARSE_VALUE_ENV, prefix), err)
	}
)

func ConvertStruct2Map(ctx context.Context, obj interface{}) map[string]interface{} {
	m := map[string]interface{}{}
	if obj != nil {
		values := reflect.ValueOf(obj).Elem()
		types := values.Type()

		for i := 0; i < values.NumField(); i++ {
			key := types.Field(i).Tag.Get("filter")
			fmt.Sprint(key)
			value := values.Field(i)

			if key != "" && !value.IsNil() {
				m[key] = value.Interface()
			}
		}
	}

	return m
}

var ConvertMap2StringSQL = func(cond map[string]interface{}) ([]string, []interface{}) {
	sqls := []string{}
	values := []interface{}{}

	for k, v := range cond {
		operator := "="
		if k != "" && v != nil {
			typeValue := fmt.Sprintf("%T", v)
			if strings.Contains(typeValue, "[]") {
				operator = "IN"
			}
			sqls = append(sqls, fmt.Sprintf("%s %s ?", k, operator))

			values = append(values, v)
		}
	}

	return sqls, values
}

type osENV struct {
	name  string
	value string
}

func (o *osENV) ParseInt() (value int64, err error) {
	v, err := strconv.ParseInt(o.value, 10, 64)
	if err != nil {
		return v, FormatErr(o.name, err)
	}
	return v, nil
}

func (o *osENV) ParseUInt() (value uint64, err error) {
	v, err := strconv.ParseUint(o.value, 10, 64)
	if err != nil {
		return v, FormatErr(o.name, err)
	}
	return v, nil
}

func (o *osENV) ParseString() (value string, err error) {
	return o.value, nil
}

func (o *osENV) ParseBool() (value bool, err error) {
	v, err := strconv.ParseBool(o.value)
	if err != nil {
		return v, FormatErr(o.name, err)
	}
	return v, nil
}

func (o *osENV) ParseFloat() (value float64, err error) {
	v, err := strconv.ParseFloat(o.value, 64)
	if err != nil {
		return v, FormatErr(o.name, err)
	}
	return v, nil
}

func GetOSEnv(envName string) *osENV {
	value := os.Getenv(envName)
	return &osENV{name: envName, value: value}
}

func ValidCronTab(syntax string) (cron.Schedule, error) {
	parse := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	sche, err := parse.Parse(syntax)
	if err != nil {
		return nil, err
	}

	return sche, nil
}

func GetOffset(page int, pageSize int) int {
	return pageSize * (page - 1)
}

func UnmarshalJSON(input string) (datatypes.JSON, error) {
	var result datatypes.JSON
	if input != "" {
		err := json.Unmarshal([]byte(input), &result)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func Contains(slice []int, item int) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func ContainsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func ExtractDomain(origin string) string {
	u, err := url.Parse(origin)
	if err != nil {
		return ""
	}
	parts := strings.Split(u.Hostname(), ".")
	if len(parts) > 2 {
		// Return the last two parts (main domain)
		return strings.Join(parts[len(parts)-2:], ".")
	}
	return u.Hostname()
}

// GetCurrentTime get the current time and check if it's +7 timezone
func GetCurrentTime() (err error, now time.Time) {
	now = time.Now()
	_, offset := now.Zone()
	if offset != 25200 { // 7 hours * 60 minutes * 60 seconds
		err = fmt.Errorf("Timezone is not +7")
	}
	return err, now
}

// GetStartEndOfDay returns 0h0m0s of today and 0h0m0s of tomorrow
func GetStartEndOfDay(tm time.Time) (startTime, endTime time.Time) {
	loc := tm.Location()
	year, month, day := tm.Date()
	startTime = time.Date(year, month, day, 0, 0, 0, 0, loc)
	endTime = time.Date(year, month, day+1, 0, 0, 0, 0, loc)
	return startTime, endTime
}

// GetStartEndOfWeek returns Monday 0h0m0s of current week and Monday 0h0m0s of next week
func GetStartEndOfWeek(tm time.Time) (startTime, endTime time.Time) {
	loc := tm.Location()

	weekday := time.Duration(tm.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	year, month, day := tm.Date()
	startTime = time.Date(year, month, day, 0, 0, 0, 0, loc).Add(-1 * (weekday - 1) * 24 * time.Hour)

	endTime = startTime.AddDate(0, 0, 7)
	return startTime, endTime
}

// GetStartEndOfMonth returns 1st day 0h0m0s of current month and 1st day 0h0m0s of next month
func GetStartEndOfMonth(tm time.Time) (startTime, endTime time.Time) {
	loc := tm.Location()
	year, month, _ := tm.Date()
	startTime = time.Date(year, month, 1, 0, 0, 0, 0, loc)
	endTime = time.Date(year, month+1, 1, 0, 0, 0, 0, loc)
	return startTime, endTime
}

func CheckValidHour(from int, to int) (error, bool) {
	err, now := GetCurrentTime()
	if err != nil {
		return err, false
	}

	hour, _, _ := now.Clock()
	return nil, to >= hour || hour >= from
}

func ConvertUnixToTime(un float64) (error, time.Time) {
	sec, dec := math.Modf(un)
	return nil, time.Unix(int64(sec), int64(dec*(1e9)))
}

// GetCurrentTime get the current time and check if it's +7 timezone
func GetCurrentUnixTime() (err error, unixT int64) {
	now := time.Now()
	_, offset := now.Zone()
	if offset != 25200 { // 7 hours * 60 minutes * 60 seconds
		err = fmt.Errorf("Timezone is not +7")
	}
	return err, now.Unix()
}

func GenerateRandomText(length int) string {
	char := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	sum := ""
	rand.Seed(time.Now().Unix())
	for i := 1; i <= length; i++ {
		randomIndex := rand.Intn(len(char))
		pick := char[randomIndex]
		sum += pick
	}
	return sum
}

func CheckIfSliceContainStr(a string, b []string) bool {
	for _, v := range b {
		if a == v {
			return true
		}
	}
	return false
}

func CheckStringArrOverlap(a []string, b []string) bool {
	for _, va := range a {
		for _, vb := range b {
			if va == vb {
				return true
			}
		}
	}
	return false
}

func ConvertNumArrToString(a []uint) string {
	if len(a) == 0 {
		return ""
	}
	idsStr := ""
	for _, id := range a {
		idsStr += fmt.Sprintf("%v,", id)

	}
	return idsStr[:len(idsStr)-1]
}

// GetOrderID returns a random big.Int number
func GetOrderID() int {
	uuidValue := uuid.New()

	intValue := new(big.Int)
	intValue.SetString(uuidValue.String(), 16)

	currentTime := time.Now()
	dateFormatted := currentTime.Format("0601")

	orderID, err := strconv.Atoi(dateFormatted + intValue.String())
	if err != nil {
		return int(intValue.Int64())
	}
	return orderID
}
