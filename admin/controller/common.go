/*
Copyright [2018] [jc3wish]

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package controller

import (
	"strings"
	"net/http"
	"github.com/brokercap/Bifrost/admin/xgo"
	"github.com/brokercap/Bifrost/config"
)


type CommonController struct {
	xgo.Controller
}


var writeRequestOp = []string{"/add","/del","/start","/stop","/close","/deal","/update","/export","/import","kill"}

//判断是否为写操作
func (c *CommonController) checkWriteRequest(uri string) bool {
	for _,v := range writeRequestOp{
		if strings.Contains(uri,v){
			return true
		}
	}
	return false
}

func (c *CommonController) Prepare()  {
	if c.Ctx.Request.Header.Get("Authorization") != ""{
		c.basicAuthor()
	}else{
		c.normalAuthor()
		c.Data["Version"] = config.VERSION

	}
}

func (c *CommonController) basicAuthor() bool{
	UserName,Password,ok := c.Ctx.Request.BasicAuth()
	if !ok || UserName == "" {
		c.SetJsonData(ResultDataStruct{Status:-1,Msg:"Author error",Data:nil})
		c.StopServeJSON()
		return false
	}
	pwd := config.GetConfigVal("user",UserName)
	if pwd == Password{
		GroupName := config.GetConfigVal("groups",UserName)
		if GroupName != "administrator" && c.checkWriteRequest(c.Ctx.Request.RequestURI){
			c.SetJsonData(ResultDataStruct{Status:-1,Msg:"user group : [ "+GroupName+" ] no authority",Data:nil})
			c.StopServeJSON()
			return false
		}
		return true
	}else{
		c.SetJsonData(ResultDataStruct{Status:-1,Msg:"Password error",Data:nil})
		c.StopServeJSON()
	}
	return false
}

func (c *CommonController)  normalAuthor() bool{
	var sessionID= c.Ctx.Session.CheckCookieValid(c.Ctx.ResponseWriter, c.Ctx.Request)

	if sessionID != "" {
		if _,ok:=c.Ctx.Session.GetSessionVal(sessionID,"UserName");ok{
			//非administrator用户 用户，没有写操作权限
			Group,_ := c.Ctx.Session.GetSessionVal(sessionID,"Group")
			if Group.(string) != "administrator" && c.checkWriteRequest(c.Ctx.Request.RequestURI){
				c.SetJsonData(ResultDataStruct{Status:-1,Msg:"user group : [ "+Group.(string)+" ] no authority",Data:nil})
				c.StopServeJSON()
				return false
			}
			return true
		}else{
			goto toLogin
		}
	}else{
		goto toLogin
	}

toLogin:
	if c.Ctx.Request.RequestURI != "/login/index" &&  c.Ctx.Request.RequestURI != "/dologin" &&  c.Ctx.Request.RequestURI != "/logout"{
		if c.IsHtmlOutput() {
			http.Redirect(c.Ctx.ResponseWriter, c.Ctx.Request, "/login/index", http.StatusFound)
			return false
		}
		c.SetJsonData(ResultDataStruct{Status:-1,Msg:"session time out",Data:nil})
		c.StopServeJSON()
		return false
	}
	return true
}

func (c *CommonController) SetTitle(title string)  {
	c.SetData("Title",title+" - Bifrost")
}

func (c *CommonController) AddAdminTemplate(tpl ...string)  {
	for _,tplName := range tpl {
		c.AddTemplate(AdminTemplatePath("/template/"+tplName))
	}
}

func (c *CommonController) AddPluginTemplate(tpl ...string)  {
	for _,tplName := range tpl {
		c.AddTemplate(PluginTemplatePath("/plugin/"+tplName))
	}
}
