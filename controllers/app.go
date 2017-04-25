package controllers

import (
	"OAuthServer/models"
	"encoding/json"
	"github.com/astaxie/beego"
)

// Operations about App
type AppController struct {
	beego.Controller
}

// @Title CreateUser
// @Description create users
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.User.Id
// @Failure 403 body is empty
// @router /regInfo [post]
func (u *AppController) PostRegInfo() {
	caKay := u.GetString("cakey")
	msg := models.AppRegisterInfo(&caKay)
	u.Data["json"] = *msg
	u.ServeJSON()
}

// @router /getToken [post]
func (u *AppController) PostToken() {
	var reqParam models.ReqToken
	json.Unmarshal(u.Ctx.Input.RequestBody, &reqParam)
	msg := models.AccessToken(&reqParam)
	u.Data["json"] = msg
	u.ServeJSON()
}
