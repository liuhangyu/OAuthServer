package models

import (
	"encoding/hex"
	"github.com/astaxie/beego"
	"os"
	"strconv"
)

type ReqParam struct {
	AppId     string `json:"appid"`
	Timestamp string `json:"timestamp"`
	Random    string `json:"random"`
}
type ReqToken struct {
	Param     ReqParam `json:"param"`
	Signature string   `json:"signature"`
}

type RegisterInfo struct {
	AppId  string `json:"appid"`
	AppKey string `json:"appkey"`
}

type TokenInfo struct {
	Token   string `json:"access_token"`
	Timeout uint32 `json:"expire_time"`
}

var Parameters *OAuthConfiguration

func init() {
	beego.Debug("OAuthServer Init......")
	Parameters = readConfig()
	if Parameters == nil {
		os.Exit(-1)
		beego.Error("readConfig Error...")
	}
	ret := registerDB(Parameters)
	if ret == false {
		os.Exit(-1)
		beego.Error("register and create table abort...")
	}
}

func GetAppId() []byte {
	randData := RandomData()
	if randData == nil {
		return nil
	}
	return GetMd5(randData)
}

func GetAppKey() []byte {
	randData := RandomData()
	if randData == nil {
		return nil
	}
	return GetMd5(randData)
}

func GetToken(param *string) []byte {
	randData := RandomData()
	data := BytesCombine(randData, S2B(param))
	return GetMd5(data)
}

/*
func (p *RegisterInfo)  AppRegisterInfo(caKey *string) (int, bool) {
	var app AppRegTab
	ret := true
	errCode := 0

	appId :=  GetAppId()
	appKey := GetAppKey()
	if appId == nil || appKey == nil {
		errCode = 7
		ret = false
	} else {
		for {
			isExist := app.exist(hex.EncodeToString(appId[:]))
			if isExist == true {
				appId = GetAppId()
				continue
			}else{
				break
			}
		}
		p.AppId  = hex.EncodeToString(appId[:])
		p.AppKey = hex.EncodeToString(appKey[:])
		id, err := app.add(p.AppId, p.AppKey)
		if err != nil {
			errCode = 7
			ret = false
		} else {
			beego.Debug("insert new data success",id, p.AppId, p.AppKey)
		    errCode = 0
			ret = true
		}
	}
	return errCode,ret
}

func ResponseServer(interface {}) *ResponseMsg {
	var app RegisterInfo
	var response ResponseMsg
	var caKey string
	errCode, ret := app.AppRegisterInfo(&caKey)
	if ret {
		response.Code = ErrTable[errCode].ErrCode
		response.Desc = ErrTable[errCode].ErrDesc
		response.Msg  = app
	} else {
		response.Code = ErrTable[errCode].ErrCode
		response.Desc = ErrTable[errCode].ErrDesc
		response.Msg  = ""
	}
	return &response
}
*/

func AppRegisterInfo(caKey *string) *ResponseMsg {
	var app AppRegTab
	response := new(ResponseMsg)

	appId := GetAppId()
	appKey := GetAppKey()

	if appId == nil || appKey == nil {
		response.Code = ErrTable[7].ErrCode
		response.Desc = ErrTable[7].ErrDesc
		response.Msg = nil
	} else {
		for {
			isExist := app.exist(hex.EncodeToString(appId[:]))
			if isExist == true {
				appId = GetAppId()
				continue
			} else {
				break
			}
		}
		id, err := app.add(hex.EncodeToString(appId[:]), hex.EncodeToString(appKey[:]))
		if err != nil {
			beego.Error("insert data error...", appId, appKey)
			response.Code = ErrTable[7].ErrCode
			response.Desc = ErrTable[7].ErrDesc
			response.Msg = err
		} else {
			beego.Debug("insert new data success", id, appId, appKey)
			response.Code = ErrTable[0].ErrCode
			response.Desc = ErrTable[0].ErrDesc
			response.Msg = RegisterInfo{
				hex.EncodeToString(appId[:]),
				hex.EncodeToString(appKey[:]),
			}
		}
	}
	return response
}

func InsertToken(appInfo *AppRegTab, params *string) ResponseMsg {
	var response ResponseMsg
	var app AppRegTab
	var tokenInfo TokenInfoTab

	if appInfo.TokenCount < 5 {
		tokenInfo.AppId = appInfo.AppId
		tokenInfo.Token = GetTokenHexStr(params)
		tokenInfo.ProdToken = GetUTCTimeStr()
		tokenInfo.ExpireTime = Parameters.Expiretime

		err := insertTokenInfo(&tokenInfo)
		if err != nil {
			response.Code = ErrTable[12].ErrCode
			response.Desc = ErrTable[12].ErrDesc
			response.Msg = err
			beego.Error(tokenInfo.AppId, tokenInfo.Token, tokenInfo.ProdToken, tokenInfo.ExpireTime)
		} else {
			app.update(appInfo.TokenCount+1, appInfo)
			response.Code = ErrTable[0].ErrCode
			response.Desc = ErrTable[0].ErrDesc
			var respon TokenInfo
			respon.Token = tokenInfo.Token
			respon.Timeout = tokenInfo.ExpireTime
			response.Msg = respon
		}
	} else {
		size, num, _ := cleanDueToken(appInfo.AppId)
		if size-num < 5 {
			//TO insert new data
			tokenInfo.AppId = appInfo.AppId
			tokenInfo.Token = GetTokenHexStr(params)
			tokenInfo.ProdToken = GetUTCTimeStr()
			tokenInfo.ExpireTime = Parameters.Expiretime

			err := insertTokenInfo(&tokenInfo)
			if err != nil {
				response.Code = ErrTable[12].ErrCode
				response.Desc = ErrTable[12].ErrDesc
				response.Msg = err
				beego.Error(tokenInfo.AppId, tokenInfo.Token, tokenInfo.ProdToken, tokenInfo.ExpireTime)
			} else {
				app.update(size-num+1, appInfo)
				response.Code = ErrTable[0].ErrCode
				response.Desc = ErrTable[0].ErrDesc
				var respon TokenInfo
				respon.Token = tokenInfo.Token
				respon.Timeout = tokenInfo.ExpireTime
				response.Msg = respon
			}
		} else {
			response.Code = ErrTable[3].ErrCode
			response.Desc = ErrTable[3].ErrDesc
			response.Msg = ""
		}
	}
	return response
}

func AccessToken(reqParams *ReqToken) *ResponseMsg {
	var app AppRegTab
	response := new(ResponseMsg)

	if reqParams == nil ||
		reqParams.Param.AppId == "" ||
		reqParams.Param.Timestamp == "" ||
		reqParams.Param.Random == "" ||
		reqParams.Signature == "" {
		response.Code = ErrTable[1].ErrCode
		response.Desc = ErrTable[1].ErrDesc
		response.Msg = "part of the parameter is empty"
	}

	for {
		appInfo, _ := app.queryOne(reqParams.Param.AppId)
		if appInfo == nil {
			response.Code = ErrTable[1].ErrCode
			response.Desc = ErrTable[1].ErrDesc
			response.Msg = "appid is not exist"
			break
		}

		diffTime := GetUTCTimeSecond() - ParseUTCTimeFromTimestamp(reqParams.Param.Timestamp)
		if diffTime > Parameters.Verifytime || diffTime < 0 {
			beego.Debug("diffTime:", diffTime, GetUTCTimeSecond(), ParseUTCTimeFromTimestamp(reqParams.Param.Timestamp))
			response.Code = ErrTable[1].ErrCode
			response.Desc = ErrTable[1].ErrDesc
			if diffTime < 0 {
				response.Msg = "remote host time error"
			} else {
				response.Msg = "timestamp timeout"
			}
			beego.Debug("sys utc time:", GetUTCTimeStr())
			beego.Debug("param utc time:", reqParams.Param.Timestamp)
			break

		}

		random := reqParams.Param.Random
		_, error := strconv.Atoi(random)
		if len(random) != 6 || error != nil {
			response.Code = ErrTable[1].ErrCode
			response.Desc = ErrTable[1].ErrDesc
			response.Msg = "random random param"
			break
		}

		paramStr := reqParams.Param.AppId + reqParams.Param.Timestamp + reqParams.Param.Random
		ret := VerifyMd5Info(&paramStr, reqParams.Signature)
		if ret {
			*response = InsertToken(appInfo, &paramStr)
		} else {
			response.Code = ErrTable[1].ErrCode
			response.Desc = ErrTable[1].ErrDesc
			response.Msg = "verify the signature failure"
		}
		break
	}
	return response
}

/*

http://10.0.1.45:8080/v1/app/regInfo
{
  "code": "r0000",
  "desc": "success",
  "msg": {
    "appid": "637ee570a5a130b9a07bc9df79c16cc6",
    "appkey": "12b51e3abf7416bea5d8a2d3fd20aa0f"
  }
}


http://10.0.1.45:8080/v1/app/getToken

{
  "param": {
    "appid": "637ee570a5a130b9a07bc9df79c16cc6",
    "timestamp": "2017-04-21T11:52:04+08:00",
    "random": "123456"
  },
  "signature": "90816513a173646c3bc758a83322e761"
}

------------

{
  "code": "r0000",
  "desc": "success",
  "msg": {
    "access_token": "e3003ea902e1848349364704d42574d7",
    "expire_time": 600
  }
}




*/
