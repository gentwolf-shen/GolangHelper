package valid

import (
	"gentwolf/GolangHelper/convert"
	"regexp"
)

var (
	items = make(map[string]*regexp.Regexp)
)

func testStr(key, ptn, str string) bool {
	reg, bl := items[key]
	if !bl {
		reg = regexp.MustCompile(ptn)
		items[key] = reg
	}
	return reg.MatchString(str)
}

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
	ptn := "^[A-Za-z0-9-]{" + minLen + "," + maxLen + "}$"
	return testStr(key, ptn, str)
}

func IsString(str string, minLen, maxLength int) bool {
	strLen := len(str)
	return strLen >= minLen && strLen <= maxLength
}

func IsDigit(str string, minLength int, maxLength int) bool {
	min := convert.ToStr(minLength)
	max := convert.ToStr(maxLength)

	key := "isDigit_" + min + "_" + max
	ptn := "^[0-9]{" + min + "," + max + "}$"
	return testStr(key, ptn, str)
}

func IsAlpha(str string, minLength int, maxLength int) bool {
	min := convert.ToStr(minLength)
	max := convert.ToStr(maxLength)

	key := "isAlpha_" + min + "_" + max
	ptn := "^[a-zA-Z]{" + min + "," + max + "}$"
	return testStr(key, ptn, str)
}

func IsYearMonth(str string, args ...string) bool {
	split := "-"
	if len(args) > 0 {
		split = args[0]
	}

	key := "IsYearMonth" + split
	ptn := "^2[0-9]{3}" + split + "([0][1-9]|[1][0-2])$"
	return testStr(key, ptn, str)
}

func IsDate(str string, args ...string) bool {
	split := "-"
	if len(args) > 0 {
		split = args[0]
	}

	key := "IsDate" + split
	ptn := "^2[0-9]{3}" + split + "([0][1-9]|[1][0-2])" + split + "([0][1-9]|[1-2][0-9]|[3][0-1])$"
	return testStr(key, ptn, str)
}

func IsMobile(str string) bool {
	key := "IsMobile"
	ptn := "^1[345789][0-9]{9}$"
	return testStr(key, ptn, str)
}

func IsEmail(str string) bool {
	key := "IsEmail"
	ptn := `^[a-z0-9][.a-z0-9_-]*@[a-z0-9][a-z0-9-]*(\.[a-z0-9]{2,10})+$`
	return testStr(key, ptn, str)
}
