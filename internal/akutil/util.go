package akutil

import (
	"regexp"
	"strconv"

	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

func FindString(a []string, s string) int {
	for i := range a {
		if a[i] == s {
			return i
		}
	}
	return -1
}

func FindInt(a []int, s int) int {
	for i := range a {
		if a[i] == s {
			return i
		}
	}
	return -1
}

type Number interface {
	constraints.Float | constraints.Integer
}

func Percentile[T Number](data []T, perc float64) T {
	dataCopy := make([]T, len(data))
	copy(dataCopy, data)

	slices.Sort(dataCopy)

	var percentile T
	l := len(dataCopy)
	if l == 0 {
		return 0
	} else {
		i := int(float64(l) * perc)

		if i >= len(dataCopy) {
			i = len(dataCopy) - 1
		}

		percentile = dataCopy[i]
	}

	return percentile
}

type HumanRounded struct {
	Number  float64
	Rounded float64
	Factor  float64
	Unit    string
}

func HumanRound(number float64, breaks ...string) (hr HumanRounded) {
	hr.Number = number

	hr.Factor = 1.0
	for _, b := range breaks {
		t := hr.Factor * 1e3
		if number < t {
			hr.Rounded = number / hr.Factor
			hr.Unit = b
			return
		}
		hr.Factor = t
	}
	hr.Rounded = number
	hr.Unit = breaks[0]
	return
}

func StripNonAscii(s string) string {
	re := regexp.MustCompile("[[:^ascii:]]")
	t := re.ReplaceAllLiteralString(s, "")
	return t
}

func JoinIntA(ia []int, del string) (s string) {
	for ii, iv := range ia {
		if ii > 0 {
			s += del
		}
		s += strconv.Itoa(iv)
	}
	return
}
