package dto

import (
	"time"

	"github.com/anthonyzero/gateway/public"
	"github.com/gin-gonic/gin"
)

type AdminInfoOutput struct {
	ID           int       `json:"id"`           //用户ID
	Name         string    `json:"name"`         //用户名
	LoginTime    time.Time `json:"login_time"`   //登录时间
	Avatar       string    `json:"avatar"`       //头像
	Introduction string    `json:"introduction"` //介绍
	Roles        []string  `json:"roles"`        //角色列表
}

type ChangePwdInput struct {
	Password string `json:"password" form:"password" comment:"密码" example:"123456" validate:"required"` //密码
}

func (param *ChangePwdInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}
