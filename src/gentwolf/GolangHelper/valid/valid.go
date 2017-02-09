package valid

import (
	"gentwolf/GolangHelper/convert"
	"regexp"
)

var (
	items = make(map[string]*regexp.Regexp)
)

func IsUsername(str string, args ...int) bool {
	minLen := "6"
	maxLen := "20"
	switch len(args) {
	case 1:
		minLen = convert.ToStr(args[0])
	case 2:
		minLen = convert.ToStr(args[0])
		maxLen = convert.ToStr(args[1])
	}
	key := "IsUsername_" + minLen + "-" + maxLen

	reg, bl := items[key]
	if !bl {
		reg = regexp.MustCompile("^[A-Za-z0-9-]{" + minLen + "," + maxLen + "}$")
		items[key] = reg
	}

	return reg.MatchString(str)
}

func IsString(str string, minLen, maxLength int) bool {
	strLen := len(str)
	return strLen >= minLen && strLen <= maxLength
}

func IsDigit(str string, minLength int, maxLength int) bool {
	min := convert.ToStr(minLength)
	max := convert.ToStr(maxLength)

	key := "isDigit_" + min + "_" + max
	reg, bl := items[key]
	if !bl {
		reg = regexp.MustCompile("^[0-9]{" + min + "," + max + "}$")
		items[key] = reg
	}

	return reg.MatchString(str)
}

func IsAlpha(str string, minLength int, maxLength int) bool {
	min := convert.ToStr(minLength)
	max := convert.ToStr(maxLength)

	key := "isAlpha_" + min + "_" + max
	reg, bl := items[key]
	if !bl {
		reg = regexp.MustCompile("^[a-zA-Z]{" + min + "," + max + "}$")
		items[key] = reg
	}

	return reg.MatchString(str)
}

func IsYearMonth(str string, args ...string) bool {
	split := "-"
	if len(args) > 0 {
		split = args[0]
	}

	key := "IsYearMonth" + split
	reg, bl := items[key]
	if !bl {
		reg = regexp.MustCompile("^2[0-9]{3}" + split + "([0][1-9]|[1][0-2])$")
		items[key] = reg
	}
	return reg.MatchString(str)
}

func IsDate(str string, args ...string) bool {
	split := "-"
	if len(args) > 0 {
		split = args[0]
	}

	key := "IsDate" + split
	reg, bl := items[key]
	if !bl {
		reg = regexp.MustCompile("^2[0-9]{3}" + split + "([0][1-9]|[1][0-2])" + split + "([0][1-9]|[1-2][0-9]|[3][0-1])$")
		items[key] = reg
	}
	return reg.MatchString(str)
}

func IsMobile(str string) bool {
	key := "IsMobile"
	reg, bl := items[key]
	if !bl {
		reg = regexp.MustCompile("^1[345789][0-9]{9}$")
		items[key] = reg
	}
	return reg.MatchString(str)
}
