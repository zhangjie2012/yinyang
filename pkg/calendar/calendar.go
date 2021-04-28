package calendar

import (
	"fmt"
	"io/ioutil"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/zhangjie2012/yinyang/pkg/types"
)

var (
	YearStart = 1901
	YearEnd   = 2101

	YinDay2N = map[string]int{
		"初一": 1, "初二": 2, "初三": 3, "初四": 4, "初五": 5, "初六": 6, "初七": 7, "初八": 8, "初九": 9, "初十": 10,
		"十一": 11, "十二": 12, "十三": 13, "十四": 14, "十五": 15, "十六": 16, "十七": 17, "十八": 18, "十九": 19, "二十": 20,
		"廿一": 21, "廿二": 22, "廿三": 23, "廿四": 24, "廿五": 25, "廿六": 26, "廿七": 27, "廿八": 28, "廿九": 29, "三十": 30,
	}
	N2YinDay = map[int]string{
		1: "初一", 2: "初二", 3: "初三", 4: "初四", 5: "初五", 6: "初六", 7: "初七", 8: "初八", 9: "初九", 10: "初十",
		11: "十一", 12: "十二", 13: "十三", 14: "十四", 15: "十五", 16: "十六", 17: "十七", 18: "十八", 19: "十九", 20: "二十",
		21: "廿一", 22: "廿二", 23: "廿三", 24: "廿四", 25: "廿五", 26: "廿六", 27: "廿七", 28: "廿八", 29: "廿九", 30: "三十",
	}

	YinMonth2N = map[string]int{
		"正月": 1,
		"二月": 2, "閏二月": 2,
		"三月": 3, "閏三月": 3,
		"四月": 4, "閏四月": 4,
		"五月": 5, "閏五月": 5,
		"六月": 6, "閏六月": 6,
		"七月": 7, "閏七月": 7,
		"八月": 8, "閏八月": 8,
		"九月": 9, "閏九月": 9,
		"十月": 10, "閏十月": 10,
		"十一月": 11, "閏十一月": 11,
		"十二月": 12, "閏十二月": 12,
	}
	N2YinMonth = map[int]string{
		1:  "正月",
		2:  "二月",
		3:  "三月",
		4:  "四月",
		5:  "五月",
		6:  "六月",
		7:  "七月",
		8:  "八月",
		9:  "九月",
		10: "十月",
		11: "十一月",
		12: "十二月",
	}

	N2WeekDay = map[int]string{
		1: "星期一",
		2: "星期二",
		3: "星期三",
		4: "星期四",
		5: "星期五",
		6: "星期六",
		7: "星期日",
	}
)

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

func ReadFile(fn string) []string {
	bs, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Fatalf("open file failure|%s|%s", fn, err)
	}
	content := string(bs)
	return strings.Split(content, "\n")
}

var yCalendar *types.YinCalendar

func GetCalendar() *types.YinCalendar {
	return yCalendar
}

func ParseRawData(dir string) {
	data := []string{}
	for y := YearStart; y < YearEnd; y++ {
		fn := fmt.Sprintf("%s/%d.txt", dir, y)
		d := ReadFile(fn)
		data = append(data, d...)
	}

	yinYear := 1900
	yinMonth := 11
	leapMonth := false

	days := []*types.Day{}
	year2Detail := map[int]*types.YinYear{}
	year2Months := map[int][]*types.YinMonth{}
	ym2Days := map[string][]*types.Day{}
	yang2Yin := map[string]*types.Day{}
	yin2Yang := map[string]*types.Day{}
	for _, row := range data {
		fields := strings.Split(row, ",")
		if len(fields) < 3 || len(fields) > 4 {
			continue
		}
		yangYMD := fields[0]
		yindayS := fields[1]
		weekDayS := fields[2]
		solarTermS := "" // 二十四节气
		if len(fields) > 3 {
			solarTermS = fields[3]
		}

		yang := NewYang(yangYMD)
		yin := NewYinDay(&yinYear, &yinMonth, &leapMonth, yindayS)
		weekDay := NewWeekDay(weekDayS)
		solarTerm := NewSolarTerm(solarTermS)

		d := &types.Day{
			Yang:      yang,
			Yin:       yin,
			WeekDay:   weekDay,
			SolarTerm: solarTerm,
		}
		days = append(days, d)

		// 年 -> 详情
		if _, ok := year2Detail[yin.Year]; !ok {
			year2Detail[yin.Year] = &types.YinYear{
				Num:    yin.Year,
				Tian:   NewTian(yin.Year),
				Di:     NewDi(yin.Year),
				Zodiac: NewZodiac(yin.Year),
			}
		}

		// 年 -> 月份
		months := year2Months[yin.Year]
		if len(months) == 0 {
			year2Months[yin.Year] = append(year2Months[yin.Year], &types.YinMonth{Num: yin.Month, Leap: leapMonth})
		} else {
			// last is not current, add it
			if months[len(months)-1].Num != yin.Month || months[len(months)-1].Leap != leapMonth {
				year2Months[yin.Year] = append(year2Months[yin.Year], &types.YinMonth{Num: yin.Month, Leap: leapMonth})
			}
		}

		// 年月 -> 日
		ym := fmt.Sprintf("%d-%d-%d", yin.Year, yin.Month, b2i(leapMonth))
		ym2Days[ym] = append(ym2Days[ym], d)

		// 公历 -> 农历
		ymd := fmt.Sprintf("%d-%d-%d", yang.Year, yang.Month, yang.Day)
		yang2Yin[ymd] = d

		// 农历 -> 公历
		ymd = fmt.Sprintf("%d-%d-%d-%d", yin.Year, yin.Month, b2i(leapMonth), yin.Day)
		yin2Yang[ymd] = d
	}

	yCalendar = &types.YinCalendar{
		Days:        days,
		Year2Detail: year2Detail,
		Year2Months: year2Months,
		YM2Days:     ym2Days,
		Yang2Yin:    yang2Yin,
		Yin2Yang:    yin2Yang,
	}
}

// NewYang "1902年12月15日" || "2011年1月1日"
func NewYang(raw string) *types.YangDay {
	var year, month, day int
	if _, err := fmt.Sscanf(raw, "%d年%d月%d日", &year, &month, &day); err != nil {
		log.Fatalf("wrong yang fomat|%s", raw)
	}

	return &types.YangDay{
		Year:  int(year),
		Month: int(month),
		Day:   int(day),
	}
}

// NewWeekDay 星期转数字 1-7
func NewWeekDay(raw string) int {
	switch raw {
	case "星期一":
		return 1
	case "星期二":
		return 2
	case "星期三":
		return 3
	case "星期四":
		return 4
	case "星期五":
		return 5
	case "星期六":
		return 6
	case "星期日":
		return 7
	default:
		log.Fatalf("invalid week day string|%s", raw)
		return -1
	}
}

// NewSolarTerm 繁体转简体
func NewSolarTerm(raw string) string {
	switch raw {
	case "驚蟄":
		return "惊蛰"
	case "穀雨":
		return "谷雨"
	case "小滿":
		return "小满"
	case "芒種":
		return "芒种"
	case "處暑":
		return "处暑"
	default:
		// 其它的都是简体
		return raw
	}
}

func NewYinDay(yinYear *int, yinMonth *int, leapMonth *bool, raw string) *types.YinDay {
	day := 0

	// 正常日子，无年份和月份跨越
	day, exist := YinDay2N[raw]
	if exist {
		return &types.YinDay{
			Year:  *yinYear,
			Month: *yinMonth,
			Leap:  *leapMonth,
			Day:   day,
		}
	}

	// 月份跨越
	day = 1 // 肯定是月初
	month, exist := YinMonth2N[raw]
	if !exist {
		log.Fatalf("wrong yin day|%s", raw)
	}
	*yinMonth = month

	// 闰月处理
	// 1. 包含 "闰" 开始算闰月
	// 2. 如果之前是闰月，有月份跨越了，取消闰月
	if strings.Contains(raw, "閏") {
		*leapMonth = true
	} else {
		if *leapMonth {
			*leapMonth = false
		}
	}

	// 年份跨越
	if month == 1 {
		*yinYear++
	}

	return &types.YinDay{
		Year:  *yinYear,
		Month: *yinMonth,
		Leap:  *leapMonth,
		Day:   day,
	}
}

// NewTian 天干计算
func NewTian(year int) string {
	d := []string{"甲", "乙", "丙", "丁", "戊", "已", "庚", "辛", "壬", "癸"}
	r := (year - 4) % 10
	return d[r]
}

// NewDi 地支计算
func NewDi(year int) string {
	d := []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}
	r := (year - 4) % 12
	return d[r]
}

// NewZodiac 属相计算
func NewZodiac(year int) string {
	d := []string{"鼠", "牛", "虎", "兔", "龙", "蛇", "马", "羊", "猴", "鸡", "狗", "猪"}
	r := (year - 4) % 12
	return d[r]
}
