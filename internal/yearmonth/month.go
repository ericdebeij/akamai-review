package yearmonth

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func Parse(ym string) (year, month int) {
	x := strings.Split(string(ym), "-")
	year, _ = strconv.Atoi(x[0])
	month, _ = strconv.Atoi(x[1])
	return
}

func FromTime(t time.Time) (ym string) {
	ym = fmt.Sprintf("%04d-%02d", t.Year(), t.Month())
	return
}

func Add(ym string, a int) string {
	y, m := Parse(ym)
	y += a / 12
	m += a % 12
	if m > 12 {
		y += 1
		m -= 12
	} else if m < 1 {
		y -= 1
		m += 12
	}
	return fmt.Sprintf("%04d-%02d", y, m)
}
