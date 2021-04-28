package api

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/zhangjie2012/cbl-go"
	"github.com/zhangjie2012/yinyang/pkg/calendar"
	"github.com/zhangjie2012/yinyang/pkg/types"
)

// ListYearH 返回年份数据
func (s *Server) ListYearH(c *gin.Context) {
	yc := calendar.GetCalendar()
	years := []*types.YinYear{}
	for y := calendar.YearStart; y < calendar.YearEnd; y++ {
		years = append(years, yc.Year2Detail[y])
	}
	log.WithField("client_ip", c.ClientIP()).Debugf("list year success")
	cbl.SuccessResponse(c, years)
}

func (s *Server) ListYearMonthH(c *gin.Context) {
	year64, err := strconv.ParseInt(c.Param("year"), 10, 32)
	if err != nil {
		log.WithField("client_ip", c.ClientIP()).Warnf("list year month invalid params|%s", err)
		cbl.ErrorResponse(c, cbl.ErrBadRequest)
		return
	}
	year := int(year64)

	if year < calendar.YearStart || year >= calendar.YearEnd {
		log.WithField("client_ip", c.ClientIP()).Warnf("list year month out of range|%d", year)
		cbl.ErrorResponse(c, cbl.ErrBadRequest)
		return
	}

	type YinMonth struct {
		Num  int    `json:"num"`
		Leap bool   `json:"leap"`
		Name string `json:"name"`
	}
	months := []*YinMonth{}
	for _, m := range calendar.GetCalendar().Year2Months[year] {
		name := calendar.N2YinMonth[m.Num]
		if m.Leap {
			name = "闰" + name
		}
		months = append(months, &YinMonth{
			Num:  m.Num,
			Leap: m.Leap,
			Name: name,
		})
	}

	log.WithField("client_ip", c.ClientIP()).Infof("list year month success|%d", year)
	cbl.SuccessResponse(c, months)
}

func (s *Server) ListYearMonthDayH(c *gin.Context) {
	year := c.Param("year")
	month := c.Param("month")
	leap := c.Param("leap")

	ym := fmt.Sprintf("%s-%s-%s", year, month, leap)
	source, ok := calendar.GetCalendar().YM2Days[ym]
	if !ok {
		log.WithField("client_ip", c.ClientIP()).Warnf("list year month day out of range|%s", ym)
		cbl.ErrorResponse(c, cbl.ErrOutOfRange)
		return
	}

	type YinDay struct {
		Num       int    `json:"num"`
		Name      string `json:"name"`
		WeekDay   string `json:"weekday"`
		SolarTerm string `json:"solarterm"`
	}

	days := []*YinDay{}
	for _, m := range source {
		days = append(days, &YinDay{
			Num:       m.Yin.Day,
			Name:      calendar.N2YinDay[m.Yin.Day],
			WeekDay:   calendar.N2WeekDay[m.WeekDay],
			SolarTerm: m.SolarTerm,
		})
	}

	log.WithField("client_ip", c.ClientIP()).Infof("list year month day success|%s", ym)
	cbl.SuccessResponse(c, days)
}

func (s *Server) ConvertYang2YinH(c *gin.Context) {
	year := c.Param("year")
	month := c.Param("month")
	day := c.Param("day")

	ymd := fmt.Sprintf("%s-%s-%s", year, month, day)
	yd, ok := calendar.GetCalendar().Yang2Yin[ymd]
	if !ok {
		log.WithField("client_ip", c.ClientIP()).Warnf("convert yang to yin out of range|%s", ymd)
		cbl.ErrorResponse(c, cbl.ErrOutOfRange)
		return
	}

	type YinDay struct {
		YearNum    int    `json:"year_num"`
		YearTian   string `json:"year_tian"`
		YearDi     string `json:"year_di"`
		YearZodiac string `json:"year_zodiac"`
		MonthNum   int    `json:"month_num"`
		MonthLeap  bool   `json:"month_leap"`
		MonthName  string `json:"month_name"`
		DayNum     int    `json:"day_num"`
		DayName    string `json:"day_name"`
		WeekDay    string `json:"weekday"`
		SolarTerm  string `json:"solarterm"`
	}
	monthName := calendar.N2YinMonth[yd.Yin.Month]
	if yd.Yin.Leap {
		monthName = "闰" + monthName
	}
	yearDetail := calendar.GetCalendar().Year2Detail[yd.Yin.Year]
	resp := YinDay{
		YearNum:    yd.Yin.Year,
		YearTian:   yearDetail.Tian,
		YearDi:     yearDetail.Di,
		YearZodiac: yearDetail.Zodiac,
		MonthNum:   yd.Yin.Month,
		MonthLeap:  yd.Yin.Leap,
		MonthName:  monthName,
		DayNum:     yd.Yin.Day,
		DayName:    calendar.N2YinDay[yd.Yin.Day],
		WeekDay:    calendar.N2WeekDay[yd.WeekDay],
		SolarTerm:  yd.SolarTerm,
	}
	log.WithField("client_ip", c.ClientIP()).Infof("convert yang to yin success|%s", ymd)
	cbl.SuccessResponse(c, &resp)
}

func (s *Server) ConvertYin2YangH(c *gin.Context) {
	year := c.Param("year")
	month := c.Param("month")
	leap := c.Param("leap")
	day := c.Param("day")

	ymd := fmt.Sprintf("%s-%s-%s-%s", year, month, leap, day)
	yd, ok := calendar.GetCalendar().Yin2Yang[ymd]
	if !ok {
		log.WithField("client_ip", c.ClientIP()).Warnf("convert yin to yang out of range|%s", ymd)
		cbl.ErrorResponse(c, cbl.ErrOutOfRange)
		return
	}

	log.WithField("client_ip", c.ClientIP()).Infof("convert yin to yang success|%s", ymd)
	cbl.SuccessResponse(c, yd.Yang)
}
