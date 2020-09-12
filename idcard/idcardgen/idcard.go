package main

import (
	"encoding/json"
	"github.com/brickman-source/golang-utilities/idcard"
	"github.com/brickman-source/golang-utilities/rand"
	"time"
)

type Street struct {
	Code         string `json:"code" xml:"code"`
	Name         string `json:"name" xml:"name"`
	AreaCode     string `json:"areaCode" xml:"areaCode"`
	CityCode     string `json:"cityCode" xml:"cityCode"`
	ProvinceCode string `json:"provinceCode" xml:"provinceCode"`
}
type Province struct {
	Code string `json:"code" xml:"code"`
	Name string `json:"name" xml:"name"`
}

func main() {
	streets := make([]*Street, 0)
	err := json.Unmarshal([]byte(idcard.StreetsJson), &streets)
	if err != nil {
		return
	}
	provinces := make([]*Province, 0)
	err = json.Unmarshal([]byte(idcard.ProvincesJson), &provinces)
	if err != nil {
		return
	}

	for offset := 0; offset < 100; offset++ {
		areaIndex := rand.GetRandInt(0, len(streets))

		var t time.Time

		if offset == 0 {
			t = time.Now()
		} else {
			t = time.Now().Add(-time.Hour*24*time.Duration(rand.GetRandInt(10, 365))).AddDate(-offset, 0, 0)
		}
		var sexStr string
		if rand.GetRandInt(0, 2) == 1 {
			sexStr = "M"
		} else {
			sexStr = "F"
		}
		number := idcard.GenResidentIdCard(streets[areaIndex].AreaCode, t, sexStr)
		provinceName := ""
		for _, province := range provinces {
			if province.Code == streets[areaIndex].ProvinceCode {
				provinceName = province.Name
			}
		}
		valid, _, _ := idcard.IsResidentIdCard(number)
		if !valid {
			panic("invalid id:" + number)
		}
		println(number, t.Format("2006-01-02"), sexStr, provinceName)
	}

}
