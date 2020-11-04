package dto

import "time"

type AdminInfoOutput struct {
	ID           int       `json:"id"`           //用户ID
	Name         string    `json:"name"`         //用户名
	LoginTime    time.Time `json:"login_time"`   //登录时间
	Avatar       string    `json:"avatar"`       //头像
	Introduction string    `json:"introduction"` //介绍
	Roles        []string  `json:"roles"`        //角色列表
}
