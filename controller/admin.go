package controller

import (
	"encoding/json"
	"fmt"

	"github.com/anthonyzero/gateway/dao"
	"github.com/e421083458/golang_common/lib"

	"github.com/anthonyzero/gateway/middleware"

	"github.com/anthonyzero/gateway/dto"
	"github.com/anthonyzero/gateway/public"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AdminController struct{}

func AdminRegister(group *gin.RouterGroup) {
	admin := &AdminController{}
	group.GET("/user_info", admin.GetUserInfo)
	group.POST("/change_pwd", admin.ChangePwd)
}

// GetUserInfo godoc
// @Summary 管理员信息
// @Description 管理员信息
// @Tags 管理员接口
// @ID /admin/user_info
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.AdminInfoOutput} "success"
// @Router /admin/user_info [get]
func (admin *AdminController) GetUserInfo(c *gin.Context) {
	session := sessions.Default(c)
	sessInfo := session.Get(public.AdminSessionInfoKey)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), adminSessionInfo); err != nil {
		middleware.ResponseError(c, middleware.BusinessErrorCode, err)
		return
	}
	out := &dto.AdminInfoOutput{
		ID:           adminSessionInfo.ID,
		Name:         adminSessionInfo.UserName,
		LoginTime:    adminSessionInfo.LoginTime,
		Avatar:       "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		Introduction: "I am a super administrator",
		Roles:        []string{"admin"},
	}
	middleware.ResponseSuccess(c, out)
}

// ChangePwd godoc
// @Summary 修改密码
// @Description 修改密码
// @Tags 管理员接口
// @ID /admin/change_pwd
// @Accept  json
// @Produce  json
// @Param body body dto.ChangePwdInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin/change_pwd [post]
func (adminlogin *AdminController) ChangePwd(c *gin.Context) {
	param := &dto.ChangePwdInput{}
	if err := param.BindValidParam(c); err != nil {
		middleware.ResponseError(c, middleware.ParamCheckErrorCode, err)
		return
	}

	//session读取用户信息到结构体
	sess := sessions.Default(c)
	sessInfo := sess.Get(public.AdminSessionInfoKey)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), adminSessionInfo); err != nil {
		middleware.ResponseError(c, middleware.BusinessErrorCode, err)
		return
	}

	//从数据库中读取用户信息
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, middleware.InternalErrorCode, err)
		return
	}
	user := &dao.Admin{}
	user, err = user.Find(c, tx, &dao.Admin{
		UserName: adminSessionInfo.UserName,
	})
	if err != nil {
		middleware.ResponseError(c, middleware.BusinessErrorCode, err)
		return
	}
	//生成新密码 saltPassword
	saltPassword := public.GenSaltPassword(user.Salt, param.Password)
	user.Password = saltPassword
	//执行数据保存
	if err := user.Save(c, tx); err != nil {
		middleware.ResponseError(c, middleware.BusinessErrorCode, err)
		return
	}
	middleware.ResponseSuccess(c, "")
}
