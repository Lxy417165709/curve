package controller

import (
	"curve/src/model"
	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(c *gin.Context) {
	var loginForm model.LoginForm
	if err := c.ShouldBindJSON(&loginForm); err != nil {
		logs.Error(err)
		c.JSON(http.StatusBadRequest, model.Response{
			Err: err.Error(),
		})
		return
	}
	isRight, err := GlobalUserManager.IsPasswordRight(loginForm.Email, loginForm.Password)
	if err != nil {
		logs.Error(err)
		c.JSON(http.StatusBadRequest, model.Response{
			Err: err.Error(),
		})
		return
	}
	if !isRight {
		c.JSON(http.StatusBadRequest, model.Response{
			Err: "账号或密码错误.",
		})
		return
	}
	if err := GlobalUserManager.UpdateLastLoginTime(loginForm.Email); err != nil {
		logs.Error(err)
		c.JSON(http.StatusBadRequest, model.Response{
			Err: err.Error(),
		})
		return
	}

	uid, err := GlobalUserManager.GetUid(loginForm.Email)
	if err != nil {
		logs.Error(err)
		c.JSON(http.StatusBadRequest, model.Response{
			Err: err.Error(),
		})
		return
	}

	secretTokenString, err := GlobalTokenManager.GetSecretTokenString(uid)
	if err != nil {
		logs.Error(err)
		c.JSON(http.StatusBadRequest, model.Response{
			Err: err.Error(),
		})
		return
	}
	const cookieExpiredSecond = 72 * 24 * 60 * 60
	c.SetCookie(
		model.KeyForTokenInCookies,
		secretTokenString,
		cookieExpiredSecond,
		"/",
		"localhost",
		false,
		true,
	)
	userInformation, err := GlobalUserManager.GetUserInformation(uid)
	if err != nil {
		logs.Error(err)
		c.JSON(http.StatusBadRequest, model.Response{
			Err: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, model.Response{
		Msg:  "登录成功.",
		Data: userInformation,
	})
}
