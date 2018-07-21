# opts

[![CircleCI](https://img.shields.io/circleci/project/github/RedSparr0w/node-csgo-parser.svg?style=for-the-badge)](https://circleci.com/gh/acomagu/opts)

The Go Functional Options Generator.

## Example

opts_gen.go:

```
//go:generate opts -type=int -name=Year
//go:generate opts -append -import=time -type=time.Month
//go:generate opts -append -type=int -name=Day
```

And `go generate` generates:

```Go
package main

type Opt interface {
	set(interface{})
}

// Year
type yearSetter interface {
	setYear(int)
}

type yearParam struct {
	withYear bool
	year     int
}

func (p *yearParam) setYear(v int) {
	p.withYear = true
	p.year = v
}

type yearOpt struct {
	v int
}

func Year(v int) *yearOpt {
	return &yearOpt{
		v: v,
	}
}

func (o *yearOpt) set(p interface{}) {
	p.(yearSetter).setYear(o.v)
}

// Month
...
```

They are used like:

```Go
type DateOpt interface {
	Opt
	dateOpt()
}

func (*yearOpt) dateOpt()  {}
func (*monthOpt) dateOpt() {}
func (*dayOpt) dateOpt()   {}

type dateParam struct {
	yearParam
	monthParam
	dayParam
}

func date(opts ...DateOpt) time.Time {
	param := new(dateParam)
	param.year = 2000
	param.month = 1
	param.day = 1

	for _, opt := range opts {
		opt.set(param)
	}

	return time.Date(param.year, param.month, param.day, 0, 0, 0, 0, time.UTC)
}

func main() {
	fmt.Println(date(Year(2018), Month(3), Day(5))) // => 2018-03-05 00:00:00 +0000 UTC
	fmt.Println(date(Month(8)))                     // => 2000-08-01 00:00:00 +0000 UTC
}
```

See [example/](./example) or type `opts --help` for more detail.
