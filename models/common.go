package models

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/astaxie/beego"
	"io"
	"math/rand"
	"reflect"
	"time"
	"unsafe"
)

const (
	RFC3339 = "2006-01-02T15:04:05Z07:00"
)

type ErrMsg struct {
	index   int
	ErrCode string
	ErrDesc string
}

var ErrTable = [15]ErrMsg{
	{0, "r0000", "success"},
	{1, "r1001", "illegal parameter"},
	{2, "r1002", "invalid AppId and AppKey"},
	{3, "r1003", "token overrun appid application, at the same time can only apply for 5 a valid token by default"},
	{4, "r1004", "invalid token"},
	{5, "r1005", "token expired"},
	{6, "r1006", "repeat application AppId"},
	{7, "r1007", "apply for AppId failure"},
	{8, "r1008", "authentication token failure"},
	{9, "r1009", "invalid file hash"},
	{10, "r1010", "invalid transaction number"},
	{11, "r0030", "serialized data failed"},
	{12, "r2001", "database operation error"},
	{13, "r2002", "network connection error"},
	{14, "r2003", "send mail failed"},
}

type ResponseMsg struct {
	Code string      `json:"code"`
	Desc string      `json:"desc"`
	Msg  interface{} `json:"msg"`
}

func GetMd5(data []byte) []byte {
	h := md5.New()
	io.WriteString(h, string(data[:]))
	//md5Value := fmt.Sprintf("%x", h.Sum(nil))
	return h.Sum(nil)
}

func RandomData() []byte {
	data := make([]byte, 32)
	if _, err := rand.Read(data); err != nil {
		beego.Error("rand.Read() error:", err)
		return nil
	}
	//secordNano := GetUTCNanoTimeStr()
	//byteArry := S2B(&secordNano)
	byteArry := Int64ToBytes(GetLocTimeSecord())
	return BytesCombine(data, byteArry)
}

func ParseLocTimeFromTimestamp(timestamp string) int64 {
	time, _ := time.ParseInLocation("05/01/2017", timestamp, time.Local)
	return time.Unix()
}

func GetLocTimeSecord() int64 {
	timestamp := time.Now().Unix()
	return timestamp
}

func GetLocTimeStr() string {
	second := time.Now()
	year, mon, day := second.Date()
	hour, min, sec := second.Clock()
	zone, _ := second.Zone()
	connStr := fmt.Sprintf("%d-%d-%d %02d:%02d:%02d %s", year, mon, day, hour, min, sec, zone)
	return connStr

}

func ParseUTCTimeFromTimestamp(timestamp string) int64 {
	time, err := time.Parse(RFC3339, timestamp)
	if err != nil {
		beego.Error("Parse time error:", err)
	}
	return time.Unix()
}

func GetUTCTimeSecond() int64 {
	now := time.Now().UTC().Unix()
	return now
}

func GetUTCTimeStr() string {
	now := time.Now()
	timeStr := now.Format(time.RFC3339)
	return timeStr
}

func GetUTCNanoTimeStr() string {
	now := time.Now()
	timeStr := now.Format(time.RFC3339Nano)
	return timeStr
}

func B2S(buf []byte) string {
	return *(*string)(unsafe.Pointer(&buf))
}

func S2B(s *string) []byte {
	return *(*[]byte)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(s))))
}

func BytesCombine(pBytes ...[]byte) []byte {
	len := len(pBytes)
	s := make([][]byte, len)
	for index := 0; index < len; index++ {
		s[index] = pBytes[index]
	}
	sep := []byte("")
	return bytes.Join(s, sep)
}

func VerifyMd5Info(param *string, signature string) bool {
	ret := false
	data := S2B(param)
	hash := md5.New()
	hash.Write(data)
	md := hash.Sum(nil)
	verifyHash := hex.EncodeToString(md)
	if signature == verifyHash {
		ret = true
	} else {
		beego.Debug("not expect hash value:", verifyHash)
	}
	return ret
}

func GetTokenHexStr(param *string) string {
	token := GetToken(param)
	return hex.EncodeToString(token[:])
}

func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}
