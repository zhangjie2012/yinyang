package api

import "github.com/gin-gonic/gin"

func (s *Server) RegisterRouter(r *gin.Engine) {
	apiv1 := r.Group("/api/v1")
	s.RegisterRouterV1(apiv1)
}

func (s *Server) RegisterRouterV1(v1 *gin.RouterGroup) {
	v1.GET("/years", s.ListYearH)
	v1.GET("/years/:year/months", s.ListYearMonthH)
	v1.GET("/years/:year/months/:month/:leap/days", s.ListYearMonthDayH)

	v1.GET("/conv/yang-yin/:year/:month/:day", s.ConvertYang2YinH)
	v1.GET("/conv/yin-yang/:year/:month/:leap/:day", s.ConvertYin2YangH)
}
