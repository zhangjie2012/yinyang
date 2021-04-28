package calendar

import (
	"testing"
)

func TestNewYang(t *testing.T) {
	raw1 := "1902年12月15日"
	yang1 := NewYang(raw1)
	t.Log(yang1)

	raw2 := "2011年1月1日"
	yang2 := NewYang(raw2)
	t.Log(yang2)
}

func TestTianDiZodiac(t *testing.T) {
	{
		year := 2020
		tian := NewTian(year)
		di := NewDi(year)
		zo := NewZodiac(year)
		t.Log(tian, di, zo)
	}
	{
		year := 2021
		tian := NewTian(year)
		di := NewDi(year)
		zo := NewZodiac(year)
		t.Log(tian, di, zo)
	}
}
