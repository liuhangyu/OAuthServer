package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	//_ "github.com/mattn/go-sqlite3"
)

type AppRegTab struct {
	Id         int64  `orm:"column(id);pk;auto"`
	AppId      string `orm:"column(appid);size(35);unique;index"`
	AppKey     string `orm:"column(appkey);size(35)"`
	RegisTime  string `orm:"column(registime)"`
	TokenCount int    `orm:"column(tokencount)"`
}

type TokenInfoTab struct {
	Id         int64  `orm:"column(id);pk;auto"`
	AppId      string `orm:"column(appid);size(35);index"`
	Token      string `orm:"column(token);size(35)"`
	ProdToken  string `orm:"column(producetoken)"`
	ExpireTime uint32 `orm:"column(expiretime)"`
}

func registerDB(config *OAuthConfiguration) bool {
	var err error
	ret := true

	switch config.Dbdriver {
	case "mysql":
		connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", config.Dbuser, config.Dbpasswd, config.Dbhost, config.Dbport, config.Dbname)
		err = orm.RegisterDriver("mysql", orm.DRMySQL)
		if err != nil {
			beego.Error("register mysql driver abort", err)
		}
		//err := orm.RegisterDataBase("default", "mysql", "root:mobile@tcp(127.0.0.1:3306)/onchain?charset=utf8&loc=Asia%2FShanghai",5,5)
		err = orm.RegisterDataBase("default", "mysql", connStr)
		if err != nil {
			beego.Error("connect mysql abort", err)
		}
		orm.SetMaxIdleConns("default", config.Maxidleconn)
		orm.SetMaxOpenConns("default", config.Maxopenconn)
	case "pg":
		err = orm.RegisterDriver("postgres", orm.DRPostgres)
		if err != nil {
			beego.Error("register postgres driver abort", err)
		}
		err = orm.RegisterDataBase("pg", "postgres", "user=postgres password=mobile host=127.0.0.1 port=5433 dbname=onchain sslmode=disable")
		if err != nil {
			beego.Error("connect mysql abort", err)
		}
	default:
		beego.Error("connect database driver is ", config.Dbdriver)
		return true
	}
	orm.RegisterModel(new(AppRegTab), new(TokenInfoTab))
	err = orm.RunSyncdb("default", true, true)
	if err != nil {
		beego.Error(err)
		ret = false
	}

	if Parameters.Sqldebug {
		orm.Debug = true
	}
	//orm.RunCommand()
	return ret
}

func insertTokenInfo(data *TokenInfoTab) error {
	o := orm.NewOrm()
	o.Using("default")
	sql := fmt.Sprintf("insert into token_info_tab(appid, token, producetoken, expiretime) values('%s','%s','%s',%d)", data.AppId, data.Token, data.ProdToken, data.ExpireTime)
	_, err := o.Raw(sql).Exec()
	if err != nil {
		beego.Error(err)
		return err
	}
	return nil
}

func cleanDueToken(appid string) (int, int, error) {
	o := orm.NewOrm()
	o.Using("default")
	var result []TokenInfoTab
	delCount := 0

	sql := fmt.Sprintf("select * from token_info_tab where appid = '%s'", appid)
	num, err := o.Raw(sql).QueryRows(&result)
	if err == nil && num > 0 {
		for _, v := range result {
			diffTime := GetUTCTimeSecond() - ParseUTCTimeFromTimestamp(v.ProdToken)
			if diffTime >= int64(v.ExpireTime) {
				sql := fmt.Sprintf("delete from token_info_tab where id = %d", v.Id)
				_, err := o.Raw(sql).Exec()
				if err != nil {
					beego.Error("delete due token failt", err, appid)
				} else {
					delCount++
				}
			}
		}
	}
	return int(num), delCount, nil
}

func updateAppReg(data AppRegTab) bool {
	ret := true
	o := orm.NewOrm()
	o.Using("default")
	sql := fmt.Sprintf("update app_reg_tab set appid = '%s', appkey = '%s'", data.AppId, data.AppKey)
	_, err := o.Raw(sql).Exec()
	if err != nil {
		beego.Error(err)
		ret = false
	}
	return ret
}

func queryAppReg() *[]AppRegTab {
	var result []AppRegTab
	o := orm.NewOrm()
	o.Using("default")

	sql := fmt.Sprintf("select appid,appkey from app_reg_tab")
	num, err := o.Raw(sql).QueryRows(&result)
	if err == nil && num > 0 {
		for _, v := range result {
			fmt.Println(v.AppId)
		}
	}
	return &result
}

func queryAppRegByCond(appId string) *AppRegTab {
	var result AppRegTab
	o := orm.NewOrm()
	o.Using("default")

	sql := fmt.Sprintf("select appid,appkey from app_reg_tab where appid = '%s'", appId)
	num, err := o.Raw(sql).QueryRows(&result)
	if err == nil && num > 0 {
		fmt.Println(result.AppId)
	}
	return &result
}

func deleteAppRegByCond(appId string) bool {
	ret := false
	o := orm.NewOrm()
	o.Using("default")

	sql := fmt.Sprintf("delete from app_reg_tab where appid = '%s'", appId)
	res, err := o.Raw(sql).Exec()
	if err == nil {
		num, _ := res.RowsAffected()
		fmt.Println("mysql row affected nums: ", num)
		ret = true
	}
	return ret
}

func (this *AppRegTab) exist(addId string) bool {
	o := orm.NewOrm()
	o.Using("default")

	exist := o.QueryTable("app_reg_tab").Filter("appid", addId).Exist()
	return exist
}

func (this *AppRegTab) queryOne(appId string) (*AppRegTab, error) {
	o := orm.NewOrm()
	o.Using("default")

	var data AppRegTab
	err := o.QueryTable("app_reg_tab").Filter("appid", appId).One(&data)
	if err == orm.ErrNoRows {
		return nil, err
	}
	return &data, nil
}

func (this *AppRegTab) query() ([]AppRegTab, error) {
	o := orm.NewOrm()
	o.Using("default")

	var data []AppRegTab
	_, err := o.QueryTable("app_reg_tab").Filter("id__gt", 0).All(&data)
	return data, err
}

func (this *AppRegTab) add(appId string, appKey string) (int64, error) {
	o := orm.NewOrm()
	o.Using("default")

	appReg := AppRegTab{AppId: appId, AppKey: appKey, RegisTime: GetUTCTimeStr(), TokenCount: 0}
	qs := o.QueryTable("app_reg_tab")
	i, _ := qs.PrepareInsert()
	id, err := i.Insert(&appReg)
	if err != nil {
		beego.Error("insert new data abort", err, appId, appKey)
	}
	i.Close()
	return id, err
}

func (this *AppRegTab) update(count int, app *AppRegTab) (bool, error) {
	ret := true
	o := orm.NewOrm()
	o.Using("default")
	sql := fmt.Sprintf("update app_reg_tab set tokencount = %d where id = %d and appid = '%s'", count, app.Id, app.AppId)
	_, err := o.Raw(sql).Exec()
	if err != nil {
		beego.Error(err)
		ret = false
	}
	return ret, err
}

func (this *AppRegTab) delete(app *AppRegTab) error {
	o := orm.NewOrm()
	o.Using("default")
	appReg := AppRegTab{Id: app.Id}
	_, err := o.Delete(&appReg)
	return err
}
