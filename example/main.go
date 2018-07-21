package main

import (
	"fmt"
	"time"
)

type DateOpt interface {
	Opt
	dateOpt()
}

func (*yearOpt) dateOpt()  {}
func (*monthOpt) dateOpt() {}
func (*dayOpt) dateOpt()   {}
func (*hourOpt) dateOpt()  {}
func (*minOpt) dateOpt()   {}
func (*secOpt) dateOpt()   {}
func (*nSecOpt) dateOpt()  {}

type dateParam struct {
	yearParam
	monthParam
	dayParam
	hourParam
	minParam
	secParam
	nSecParam
}

func date(opts ...DateOpt) time.Time {
	param := new(dateParam)
	param.year = 2000
	param.month = 1
	param.day = 1

	for _, opt := range opts {
		opt.set(param)
	}

	return time.Date(param.year, param.month, param.day, param.hour, param.min, param.sec, param.nSec, time.UTC)
}

type AddDateOpt interface {
	Opt
	addDateOpt()
}

// Share the opts with date()
func (*yearOpt) addDateOpt()  {}
func (*monthOpt) addDateOpt() {}
func (*dayOpt) addDateOpt()   {}

type addDateParam struct {
	yearParam
	monthParam
	dayParam
}

func addDate(t time.Time, opts ...AddDateOpt) time.Time {
	param := new(addDateParam)
	for _, opt := range opts {
		opt.set(param)
	}

	return t.AddDate(param.year, int(param.month), param.day)
}

func main() {
	fmt.Println(date(Year(2018), Month(3), Day(5))) // => 2018-03-05 00:00:00 +0000 UTC
	fmt.Println(date(Min(10)))                      // => 2000-01-01 00:10:00 +0000 UTC

	fmt.Println(addDate(date(), Year(20), Day(3))) // => 2020-01-04 00:00:00 +0000 UTC
}
