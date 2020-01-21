package auth

import (
	"dijan/exception/system"
	"dijan/utils"
	"dijan/utils/hash"
	"github.com/kataras/iris"
	"time"
)

type DijanAuthAuthorization interface {
	SetCookie(save bool)
	CheckLogin()
	IsLogin() bool
}

type cookieInfo struct {
	ExpireTime int64 `json:"expire_time"`
}

type dijanAuthAuthorization struct {
	isLogin bool
	Context iris.Context
}

var cookieName = "Dijan_Auth_Cookie"

func NewAuthAuthorization(ctx *iris.Context) DijanAuthAuthorization {
	authorization := dijanAuthAuthorization{
		isLogin: false,
		Context: *ctx,
	}
	authorization.loadAuthVerification()
	return &authorization
}

func (r *dijanAuthAuthorization) CheckLogin() {
	if !r.isLogin {
		panic(systemException.SystemApiTokenVerificationFail())
	}
}

func (r *dijanAuthAuthorization) IsLogin() bool {
	return r.isLogin
}

func (r *dijanAuthAuthorization) loadAuthVerification() {
	succ := r.loadHeader()
	if !succ {
		r.loadFromCookie()
	}
}

func (r *dijanAuthAuthorization) loadHeader() bool {
	token := r.Context.GetHeader("token")
	if token == utils.GlobalSystemConfig.Server.SystemApiToken {
		r.isLogin = true
		return true
	}
	return false
}

// 从cookie载入登录信息
func (r *dijanAuthAuthorization) loadFromCookie() bool {
	defer func() {
		recover()
	}()
	cookie := r.Context.GetCookie(cookieName)
	if cookie == "" {
		return false
	}
	var cookieStruct cookieInfo
	hash.DecodeToken(cookie, &cookieStruct)
	if cookieStruct.ExpireTime <= time.Now().Unix() {
		return false
	}
	r.isLogin = true
	return true
}
// 设置cookie
func (r *dijanAuthAuthorization) SetCookie(save bool) {
	if !save {
		r.Context.SetCookieKV(cookieName, "")
		return
	}
	payload := cookieInfo{
		ExpireTime: utils.GlobalSystemConfig.Server.CookieExpires + time.Now().Unix(),
	}

	payloadString := hash.GenerateToken(payload, true)
	r.Context.SetCookieKV(cookieName, payloadString)
}


