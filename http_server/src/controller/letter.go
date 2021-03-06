package controller

import (
	"curve/src/model"
	"github.com/astaxie/beego/logs"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func HadSentLetter(c *gin.Context) {
	uid, err := checkAndGetUid(c)
	if err != nil {
		logs.Error(err)
		c.JSON(http.StatusBadRequest, model.Response{
			Err: err.Error(),
		})
		return
	}
	letters, err := GlobalLetterManager.GetHadSentLetters(uid)
	if err != nil {
		logs.Error(err)
		c.JSON(http.StatusBadRequest, model.Response{
			Err: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, model.Response{
		Msg:  "获取已发送消息成功.",
		Data: letters,
	})
}

func HadReceivedLetter(c *gin.Context) {
	senderUIDString := c.Param("senderUID")
	senderUID, err := strconv.Atoi(senderUIDString)
	if err != nil {
		logs.Error(err)
		c.JSON(http.StatusBadRequest, model.Response{
			Err: err.Error(),
		})
		return
	}
	receiverUID, err := checkAndGetUid(c)
	if err != nil {
		logs.Error(err)
		c.JSON(http.StatusBadRequest, model.Response{
			Err: err.Error(),
		})
		return
	}
	letters, err := GlobalLetterManager.GetHadReceivedLetters(senderUID, receiverUID)
	if err != nil {
		logs.Error(err)
		c.JSON(http.StatusBadRequest, model.Response{
			Err: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, model.Response{
		Msg:  "获取消息成功.",
		Data: letters,
	})
}

func SendLetter(c *gin.Context) {
	var sendLetterForm model.SendLetterForm
	if err := c.ShouldBindJSON(&sendLetterForm); err != nil {
		logs.Error(err)
		c.JSON(http.StatusBadRequest, model.Response{
			Err: err.Error(),
		})
		return
	}
	secretTokenString, err := c.Cookie(model.KeyForTokenInCookies)
	if err != nil {
		logs.Error(err)
		c.JSON(http.StatusBadRequest, model.Response{
			Err: err.Error(),
		})
		return
	}
	senderUID, err := GlobalTokenManager.GetUidFromSecretTokenString(secretTokenString)
	if err != nil {
		logs.Error(err)
		c.JSON(http.StatusBadRequest, model.Response{
			Err: err.Error(),
		})
		return
	}
	if err := GlobalLetterManager.StoreLetter(
		senderUID,
		sendLetterForm.ReceiverUID,
		sendLetterForm.Content,
	); err != nil {
		logs.Error(err)
		c.JSON(http.StatusBadRequest, model.Response{
			Err: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, model.Response{
		Msg:  "消息发送成功.",
	})
}
