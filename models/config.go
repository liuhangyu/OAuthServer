package models

import (
	"github.com/astaxie/beego"
	"os"
)

type OAuthConfiguration struct {
	Appname         string
	Httpport        int64
	Runmode         string
	Autorender      bool
	Copyrequestbody bool
	EnableDocs      bool
	Dbdriver        string
	Dbname          string
	Dbuser          string
	Dbpasswd        string
	Dbhost          string
	Dbport          int64
	Maxidleconn     int
	Maxopenconn     int
	Sqldebug        bool
	Verifytime      int64
	Expiretime      uint32
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func readConfig() *OAuthConfiguration {
	Parameters = new(OAuthConfiguration)
	Parameters.Appname = beego.AppConfig.DefaultString("appname", "OAuthServer")
	Parameters.Httpport = beego.AppConfig.DefaultInt64("httpport", 8080)
	Parameters.Runmode = beego.AppConfig.DefaultString("runmode", "dev")
	Parameters.Autorender = beego.AppConfig.DefaultBool("autorender", false)
	Parameters.Copyrequestbody = beego.AppConfig.DefaultBool("copyrequestbody", true)
	Parameters.EnableDocs = beego.AppConfig.DefaultBool("enableDocs", true)

	Parameters.Dbdriver = beego.AppConfig.DefaultString("dbdriver", "mysql")
	Parameters.Dbname = beego.AppConfig.DefaultString("dbname", "onchain")
	Parameters.Dbuser = beego.AppConfig.DefaultString("dbuser", "mysql")
	Parameters.Dbpasswd = beego.AppConfig.DefaultString("dbpasswd", "mysql")
	Parameters.Dbhost = beego.AppConfig.DefaultString("dbhost", "mysql")
	Parameters.Dbport = beego.AppConfig.DefaultInt64("dbport", 3306)
	Parameters.Maxidleconn = beego.AppConfig.DefaultInt("maxidleconn", 30)
	Parameters.Maxopenconn = beego.AppConfig.DefaultInt("maxopenconn", 30)
	Parameters.Sqldebug = beego.AppConfig.DefaultBool("sqldebug", false)
	Parameters.Verifytime = beego.AppConfig.DefaultInt64("verifytime", 30)
	Parameters.Expiretime = uint32(beego.AppConfig.DefaultInt64("expiretime", 30))
	return Parameters
}
