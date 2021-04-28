package types

type YangDay struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

type YinDay struct {
	Year  int
	Month int
	Leap  bool // 闰月标识
	Day   int
}

type Day struct {
	Yang      *YangDay // 公历
	Yin       *YinDay  // 农历
	WeekDay   int      // 星期几（从 1 开始到 7）
	SolarTerm string   // 节气
}

type YinYear struct {
	Num    int    `json:"num"`    // 2020, 2021, ...
	Tian   string `json:"tian"`   // 地支
	Di     string `json:"di"`     // 天干
	Zodiac string `json:"zodiac"` // 属相（生肖）
}

type YinMonth struct {
	Num  int  // 1,2,3...,12
	Leap bool // 闰月
}

// YinCalendar 阴历为主的日历
type YinCalendar struct {
	Days []*Day

	Year2Detail map[int]*YinYear    // 年 -> 详情
	Year2Months map[int][]*YinMonth // 年 -> 月份
	YM2Days     map[string][]*Day   // 年月 -> 日 key format "2020-7-0/1"  0/1 是否为闰月

	Yang2Yin map[string]*Day // 公历 -> 农历 key format "2020-7-1"
	Yin2Yang map[string]*Day // 农历 -> 公历 key format "2020-7-0/1-1"  0/1 是否为闰月
}
