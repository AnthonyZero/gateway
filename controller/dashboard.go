package controller

import (
	"errors"

	"github.com/anthonyzero/gateway/dao"
	"github.com/anthonyzero/gateway/dto"
	"github.com/anthonyzero/gateway/middleware"
	"github.com/anthonyzero/gateway/public"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
)

type DashboardController struct{}

func DashboardRegister(group *gin.RouterGroup) {
	service := &DashboardController{}
	group.GET("/panel_group_data", service.PanelGroupData)
	group.GET("/flow_stat", service.FlowStat)
	group.GET("/service_stat", service.ServiceStat)
}

// PanelGroupData godoc
// @Summary 指标统计
// @Description 指标统计
// @Tags 首页大盘
// @ID /dashboard/panel_group_data
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.PanelGroupDataOutput} "success"
// @Router /dashboard/panel_group_data [get]
func (service *DashboardController) PanelGroupData(c *gin.Context) {
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, middleware.InternalErrorCode, err)
		return
	}
	//获取服务数量
	serviceInfo := &dao.ServiceInfo{}
	_, serviceNum, err := serviceInfo.PageList(c, tx, &dto.ServiceListInput{PageSize: 1, PageNo: 1})
	if err != nil {
		middleware.ResponseError(c, middleware.InternalErrorCode, err)
		return
	}
	//获取租户数量
	app := &dao.App{}
	_, appNum, err := app.APPList(c, tx, &dto.APPListInput{PageNo: 1, PageSize: 1})
	if err != nil {
		middleware.ResponseError(c, middleware.InternalErrorCode, err)
		return
	}
	// counter, err := public.FlowCounterHandler.GetCounter(public.FlowTotal)
	// if err != nil {
	// 	middleware.ResponseError(c, middleware.InternalErrorCode, err)
	// 	return
	// }
	out := &dto.PanelGroupDataOutput{
		ServiceNum: serviceNum,
		AppNum:     appNum,
		// TodayRequestNum: counter.TotalCount,
		// CurrentQPS:      counter.QPS,
	}
	middleware.ResponseSuccess(c, out)
}

// ServiceStat godoc
// @Summary 服务类型占比统计
// @Description 服务类型占比统计
// @Tags 首页大盘
// @ID /dashboard/service_stat
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.DashServiceStatOutput} "success"
// @Router /dashboard/service_stat [get]
func (service *DashboardController) ServiceStat(c *gin.Context) {
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, middleware.InternalErrorCode, err)
		return
	}
	serviceInfo := &dao.ServiceInfo{}
	list, err := serviceInfo.GroupByLoadType(c, tx)
	if err != nil {
		middleware.ResponseError(c, middleware.InternalErrorCode, err)
		return
	}
	legend := []string{}
	for index, item := range list {
		name, ok := public.LoadTypeMap[item.LoadType]
		if !ok {
			middleware.ResponseError(c, middleware.InternalErrorCode, errors.New("load_type not found"))
			return
		}
		list[index].Name = name
		legend = append(legend, name)
	}
	out := &dto.DashServiceStatOutput{
		Legend: legend,
		Data:   list,
	}
	middleware.ResponseSuccess(c, out)
}

// FlowStat godoc
// @Summary 流量统计
// @Description 流量统计
// @Tags 首页大盘
// @ID /dashboard/flow_stat
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.ServiceStatOutput} "success"
// @Router /dashboard/flow_stat [get]
func (service *DashboardController) FlowStat(c *gin.Context) {
	// counter, err := public.FlowCounterHandler.GetCounter(public.FlowTotal)
	// if err != nil {
	// 	middleware.ResponseError(c, middleware.ParamCheckErrorCode, err)
	// 	return
	// }
	todayList := []int64{}
	// currentTime := time.Now()
	// for i := 0; i <= currentTime.Hour(); i++ {
	// 	dateTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), i, 0, 0, 0, lib.TimeLocation)
	// 	hourData, _ := counter.GetHourData(dateTime)
	// 	todayList = append(todayList, hourData)
	// }

	yesterdayList := []int64{}
	// yesterTime := currentTime.Add(-1 * time.Duration(time.Hour*24))
	// for i := 0; i <= 23; i++ {
	// 	dateTime := time.Date(yesterTime.Year(), yesterTime.Month(), yesterTime.Day(), i, 0, 0, 0, lib.TimeLocation)
	// 	hourData, _ := counter.GetHourData(dateTime)
	// 	yesterdayList = append(yesterdayList, hourData)
	// }
	middleware.ResponseSuccess(c, &dto.ServiceStatOutput{
		Today:     todayList,
		Yesterday: yesterdayList,
	})
}