package controller

import (
	"encoding/base64"
	"errors"
	"strings"
	"time"

	"github.com/anthonyzero/gateway/dao"
	"github.com/anthonyzero/gateway/dto"
	"github.com/anthonyzero/gateway/middleware"
	"github.com/anthonyzero/gateway/public"
	"github.com/dgrijalva/jwt-go"
	"github.com/e421083458/golang_common/lib"

	"github.com/gin-gonic/gin"
)

type OAuthController struct{}

func OAuthRegister(group *gin.RouterGroup) {
	oauth := &OAuthController{}
	group.POST("/tokens", oauth.Tokens)
}

// Tokens godoc
// @Summary 获取TOKEN
// @Description 获取TOKEN
// @Tags OAUTH
// @ID /oauth/tokens
// @Accept  json
// @Produce  json
// @Param body body dto.TokensInput true "body"
// @Success 200 {object} middleware.Response{data=dto.TokensOutput} "success"
// @Router /oauth/tokens [post]
func (oauth *OAuthController) Tokens(c *gin.Context) {
	params := &dto.TokensInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, middleware.ParamCheckErrorCode, err)
		return
	}

	//获取header头 空格分开 Auth信息
	splits := strings.Split(c.GetHeader("Authorization"), " ")
	if len(splits) != 2 {
		middleware.ResponseError(c, middleware.BusinessErrorCode, errors.New("用户名或密码格式错误"))
		return
	}

	//对第二个元素 解密获取appsecret( 类似username:password 这样的组成)
	appSecret, err := base64.StdEncoding.DecodeString(splits[1])
	if err != nil {
		middleware.ResponseError(c, middleware.BusinessErrorCode, err)
		return
	}

	//fmt.Println("appSecret", string(appSecret))
	//  从header中取出 app_id secret
	//  取出 app_list
	//  匹配 app_id
	//  基于appid 生成jwt token
	//  生成 output
	parts := strings.Split(string(appSecret), ":")
	if len(parts) != 2 {
		middleware.ResponseError(c, middleware.BusinessErrorCode, errors.New("用户名或密码格式错误"))
		return
	}

	appList := dao.AppManagerHandler.GetAppList()
	for _, appInfo := range appList {
		//对应header头上的appid 和secret
		if appInfo.AppID == parts[0] && appInfo.Secret == parts[1] {
			claims := jwt.StandardClaims{
				Issuer:    appInfo.AppID,
				ExpiresAt: time.Now().Add(public.JwtExpires * time.Second).In(lib.TimeLocation).Unix(), //3600秒
			}
			token, err := public.JwtEncode(claims)
			if err != nil {
				middleware.ResponseError(c, middleware.InternalErrorCode, err)
				return
			}
			//返回token 信息
			output := &dto.TokensOutput{
				ExpiresIn:   public.JwtExpires,
				TokenType:   "Bearer",
				AccessToken: token,
				Scope:       "read_write",
			}
			middleware.ResponseSuccess(c, output)
			return
		}
	}
	middleware.ResponseError(c, middleware.BusinessErrorCode, errors.New("未匹配正确APP信息"))
}
