package idcard

import (
	"fmt"
	"github.com/brickman-source/golang-utilities/log"
	"github.com/brickman-source/golang-utilities/rand"
	"strings"
	"time"
)

func IsResidentIdCard(number string) (valid bool, birthDate time.Time, sex string) {
	l := len(number)

	valid = false

	if l != 15 && l != 18 {
		return
	}

	var sexStr string
	var err error

	if l == 15 {
		sexStr = number[14:15]
		birthDate, err = time.ParseInLocation("20060102", "19"+number[6:12], time.Local)
		if err != nil {
			return
		}
	} else {
		sexStr = number[16:17]
		birthDate, err = time.ParseInLocation("20060102", number[6:14], time.Local)
		if err != nil {
			log.Debugf( "get birthdate from idcard(%s) err:%s", number, err.Error())
			return
		}
	}

	if birthDate.After(time.Now()) {
		return
	}

	if sexStr == "0" || sexStr == "2" || sexStr == "4" || sexStr == "6" || sexStr == "8" {
		sex = "F" //女
	} else {
		sex = "M" //男
	}

	if l == 15 {
		valid = true
		return
	}

	factors := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	sum := 0
	for idx := 0; idx < 17; idx++ {
		sum = sum + int(number[idx]-'0')*factors[idx]
	}
	code := "10X98765432"
	h := code[sum%11]
	number = strings.ToUpper(number)
	if h != number[17] {
		//println("need:", string(h))
		return
	}

	valid = true
	return
}
func UpgradeResidentIdCard(idCardNo string) string {
	l := len(idCardNo)

	if l != 15 {
		return idCardNo
	}
	newIdCard := idCardNo[:6] + "19" + idCardNo[6:]

	factors := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	sum := 0
	for idx := 0; idx < 17; idx++ {
		sum = sum + int(newIdCard[idx]-'0')*factors[idx]
	}
	code := "10X98765432"
	h := code[sum%11]
	return newIdCard[:17] + string(h)
}
func GenResidentIdCard(areaCode string, birthDate time.Time, sex string) (number string) {
	var sexStr string
	if sex == "M" { // 男
		sexStr = fmt.Sprintf("%d", rand.GetRandInt(1, 1000)*2%10)
	} else {
		sexStr = fmt.Sprintf("%d", (rand.GetRandInt(1, 1000)*2+1)%10)
	}

	birthDateStr := birthDate.Format("20060102")

	number = fmt.Sprintf("%s%s%d%s", areaCode, birthDateStr, rand.GetRandInt(10, 99), sexStr)
	factors := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	sum := 0
	for idx := 0; idx < 17; idx++ {
		sum = sum + int(number[idx]-'0')*factors[idx]
	}
	code := "10X98765432"
	h := code[sum%11]
	number = strings.ToUpper(number) + string(h)
	return
}
