package base

import (
	"html/template"
	. "pms/init"
	"pms/models/base"
	. "pms/utils"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/beego/i18n"
)

var (
	AppVer string
	IsPro  bool
)

type BaseController struct {
	beego.Controller
	IsAdmin   bool
	UserName  string
	URL       string
	LastLogin time.Time
	User      base.User
	i18n.Locale
}

// Prepare implemented Prepare method for baseRouter.
func (ctl *BaseController) Prepare() {
	// Setting properties.
	ctl.StartSession()
	ctl.Data["AppVer"] = AppVer
	ctl.Data["IsPro"] = IsPro
	ctl.Data["xsrf"] = template.HTML(ctl.XSRFFormHTML())
	ctl.Data["PageStartTime"] = time.Now()

	// Redirect to make URL clean.
	if ctl.setLangVer() {
		i := strings.Index(ctl.Ctx.Request.RequestURI, "?")
		ctl.Redirect(ctl.Ctx.Request.RequestURI[:i], 302)
		return
	}

	user := ctl.GetSession("User")
	if user != nil {
		ctl.User = user.(base.User)
		ctl.Data["user"] = user
		ctl.Data["LastLogin"] = ctl.GetSession("LastLogin")
	} else {
		if ctl.Ctx.Request.RequestURI != "/login/in" {
			ctl.Redirect("/login/in", 302)
		}

		ctl.Data["LastLogin"] = ctl.GetSession("LastLogin")
		ctl.Data["LastIp"] = ctl.GetSession("LastIp")
	}

}

// setLangVer sets site language version.
func (ctl *BaseController) setLangVer() bool {
	isNeedRedir := false
	hasCookie := false

	// 1. Check URL arguments.
	lang := ctl.Input().Get("lang")

	// 2. Get language information from cookies.
	if len(lang) == 0 {
		lang = ctl.Ctx.GetCookie("lang")
		hasCookie = true
	} else {
		isNeedRedir = true
	}

	// Check again in case someone modify by purpose.
	if !i18n.IsExist(lang) {
		lang = ""
		isNeedRedir = false
		hasCookie = false
	}

	// 3. Get language information from 'Accept-Language'.
	if len(lang) == 0 {
		al := ctl.Ctx.Request.Header.Get("Accept-Language")
		if len(al) > 4 {
			al = al[:5] // Only compare first 5 letters.
			if i18n.IsExist(al) {
				lang = al
			}
		}
	}

	// 4. Default language is English.
	if len(lang) == 0 {
		lang = "en-US"
		isNeedRedir = false
	}

	curLang := LangType{
		Lang: lang,
	}

	// Save language information in cookies.
	if !hasCookie {
		ctl.Ctx.SetCookie("lang", curLang.Lang, 1<<31-1, "/")
	}

	restLangs := make([]*LangType, 0, len(LangTypes)-1)
	for _, v := range LangTypes {
		if lang != v.Lang {
			restLangs = append(restLangs, v)
		} else {
			curLang.Name = v.Name
		}
	}

	// Set language properties.
	ctl.Lang = lang
	ctl.Data["Lang"] = curLang.Lang
	ctl.Data["CurLang"] = curLang.Name
	ctl.Data["RestLangs"] = restLangs

	return isNeedRedir
}
