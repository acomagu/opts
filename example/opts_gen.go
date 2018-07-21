package main

//go:generate opts -type=int -name=Year
//go:generate opts -append -import=time -type=time.Month
//go:generate opts -append -type=int -name=Day
//go:generate opts -append -type=int -name=Hour
//go:generate opts -append -type=int -name=Min
//go:generate opts -append -type=int -name=Sec
//go:generate opts -append -type=int -name=NSec
