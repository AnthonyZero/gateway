package controller

import (
	"encoding/json"
	"time"

	"github.com/anthonyzero/gateway/dao"
	"github.com/anthonyzero/gateway/dto"
	"github.com/anthonyzero/gateway/middleware"
	"github.com/anthonyzero/gateway/public"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AdminLoginController struct{}

//路由注册
func AdminLoginRegister(group *gin.RouterGroup) {
	adminLogin := &AdminLoginController{}
	group.POST("/login", adminLogin.Login)
}

// login godoc
// @Summary 管理员登陆
// @Description 管理员登陆
// @Tags 管理员接口
// @ID /admin_login/login
// @Accept  json
// @Produce  json
// @Param body body dto.AdminLoginInput true "body"
// @Success 200 {object} middleware.Response{data=dto.AdminLoginOutput} "success"
// @Router /admin_login/login [post]
func (adminLogin *AdminLoginController) Login(c *gin.Context) {
	params := &dto.AdminLoginInput{}
	if err := params.BindValidParam(c); err != nil {
		//参数校验失败
		middleware.ResponseError(c, middleware.ParamCheckErrorCode, err)
		return
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, middleware.InternalErrorCode, err)
	}
	//查询用户信息
	userInfo := &dao.Admin{}
	userInfo, err = userInfo.LoginCheck(c, tx, params)
	if err != nil {
		middleware.ResponseError(c, middleware.CustomizeCode, err)
		return
	}

	//设置session信息
	sessInfo := &dto.AdminSessionInfo{
		ID:        userInfo.Id,
		UserName:  userInfo.UserName,
		LoginTime: time.Now(),
	}
	sessBts, err := json.Marshal(sessInfo) //结构体转json
	if err != nil {
		middleware.ResponseError(c, middleware.CustomizeCode, err)
		return
	}
	sess := sessions.Default(c)
	sess.Set(public.AdminSessionInfoKey, string(sessBts))
	sess.Save()

	out := &dto.AdminLoginOutput{Token: userInfo.UserName}
	middleware.ResponseSuccess(c, out)
}
