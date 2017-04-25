1:环境设置：
beego 编译环境设置：设置GOPATH环境之后，执行go get github.com/astaxie/beego
会在GOPATH bin目录下安装bee命名

GOPATH已经设置，在src目录下放置依赖go语言包: 将当前vendor目录中的github.com文件夹拷贝到src目录下

如下目录平级关系：
liuhy@liuhy ~/linuxshare/work/src $ tree -d -L 1 github.com OAuthServer
github.com
├── astaxie
├── beego
├── gorilla
├── go-sql-driver
├── lib
└── mattn
OAuthServer
├── conf
├── controllers
├── doc
├── models
├── routers
├── tests
├── tmp
└── vendor

2：OAuthServer 编译运行 bee run  
首先需要安装mysql数据库
mysql安装参照：
http://www.cnblogs.com/bookwed/p/5896619.html

3：逻辑部分

OAuthServer\conf\app.conf 配置文件：
appname = OAuthServer   程序名称
httpport = 8080         rest接口端口
runmode = dev           程序发布模式 框架默认值，不用设置
autorender = false        框架默认值，不用设置
copyrequestbody = true    框架默认值，不用设置
EnableDocs = true         框架默认值，不用设置

dbdriver =  mysql         表示连接的是mysql数据库
dbname = onchain          数据库名称
dbuser = root             连接数据库的用户名
dbpasswd = mobile         连接数据库的密码
dbhost = 127.0.0.1        连接数据库的地址
dbport = 3306             数据库端口
maxidleconn = 5           数据库连接池最小连接数
maxopenconn = 5           数据库连接池最大连接数
sqldebug = true           是否开启执行sql日志 true为开启
verifytime = 300000       申请token时，应用程序和本地时间之间的误差值 单位秒
expiretime = 30           token过期时间 单位秒


bee run 启动OAuthServer程序后：
可以使用postman或者狐火的psoter进行测试：

1》申请appid 、 appkey 接口

例如：使用poster访问如下接口
http://10.0.1.45:8080/v1/app/regInfo

返回json格式数据如下：
{
  "code": "r0000",
  "desc": "success",
  "msg": {
    "appid": "72592fd77f0a1153bdbcfd97b6ce36ab",
    "appkey": "5b722d7b5029ae69cf4894425a5d7203"
  }
}

可以进入mysql 命令行查看数据：
mysql -uroot -pmobile
mysql> use onchain

mysql> show tables;
+-------------------+
| Tables_in_onchain |
+-------------------+
| app_reg_tab       |
| my_test           |
| token_info_tab    |
| user              |
+-------------------+
4 rows in set (0.00 sec)

mysql> select * from app_reg_tab;
+----+----------------------------------+----------------------------------+---------------------------+------------+
| id | appid                            | appkey                           | registime                 | tokencount |
+----+----------------------------------+----------------------------------+---------------------------+------------+
|  1 | 72592fd77f0a1153bdbcfd97b6ce36ab | 5b722d7b5029ae69cf4894425a5d7203 | 2017-04-25T15:25:00+08:00 |          0 |
+----+----------------------------------+----------------------------------+---------------------------+------------+
1 row in set (0.00 sec)

其中registime是申请appid appkey时刻的UTC时间
tokencount表示这个app使用token的次数（每个应用app最多可以申请5个token）

2》申请token接口

例如：使用poster访问如下接口
http://10.0.1.45:8080/v1/app/getToken
同时传递参数如下：
{
  "param": {
    "appid": "72592fd77f0a1153bdbcfd97b6ce36ab",
    "timestamp": "2017-04-25T15:33:51+08:00",
    "random": "123456"
  },
  "signature": "b385dbf71786e554ad2754832404d6d1"
}

说明：
appid是应用程序之前申请的appid
timestamp是应用程序是utc时间
random是6位的数字字符串
signature是将appid、timestamp、random拼接之后的md5摘要


返回的json数据如下：
{
  "code": "r0000",
  "desc": "success",
  "msg": {
    "access_token": "ac08d0051c4d25aaab58d003e0397ffe",
    "expire_time": 30
  }
}

说明:
access_token 表示生成的token随机串
expire_time是表示token的失效时间

数据表token_info_tab记录生成的token信息：

mysql> select * from token_info_tab;
+----+----------------------------------+----------------------------------+---------------------------+------------+
| id | appid                            | token                            | producetoken              | expiretime |
+----+----------------------------------+----------------------------------+---------------------------+------------+
|  2 | 72592fd77f0a1153bdbcfd97b6ce36ab | 2257565acd67a99ebf60662fe150e958 | 2017-04-25T15:38:17+08:00 |         30 |
|  3 | 72592fd77f0a1153bdbcfd97b6ce36ab | 0920b922853a5072d6de830366f18bbc | 2017-04-25T15:38:19+08:00 |         30 |
|  4 | 72592fd77f0a1153bdbcfd97b6ce36ab | 7a35c1a12cdd9a65e9c087503166fcef | 2017-04-25T15:38:21+08:00 |         30 |
|  5 | 72592fd77f0a1153bdbcfd97b6ce36ab | 8533eeeeafe1e66cdc5a4705bd7aaee8 | 2017-04-25T15:38:22+08:00 |         30 |
|  6 | 72592fd77f0a1153bdbcfd97b6ce36ab | aca5bcd5ed443f0c61bcae62fed88ed4 | 2017-04-25T15:38:24+08:00 |         30 |
+----+----------------------------------+----------------------------------+---------------------------+------------+
5 rows in set (0.00 sec)

mysql> select * from app_reg_tab;
+----+----------------------------------+----------------------------------+---------------------------+------------+
| id | appid                            | appkey                           | registime                 | tokencount |
+----+----------------------------------+----------------------------------+---------------------------+------------+
|  1 | 72592fd77f0a1153bdbcfd97b6ce36ab | 5b722d7b5029ae69cf4894425a5d7203 | 2017-04-25T15:25:00+08:00 |          5 |
+----+----------------------------------+----------------------------------+---------------------------+------------+
1 row in set (0.00 sec)


说明：如果appid：72592fd77f0a1153bdbcfd97b6ce36ab申请了5个token 则app_reg_tab字段tokencount为5
在72592fd77f0a1153bdbcfd97b6ce36ab申请的5个token没有失效的情况下，再次申请则返回失败数据：
{
  "code": "r1003",
  "desc": "token overrun appid application, at the same time can only apply for 5 a valid token by default",
  "msg": ""
}





