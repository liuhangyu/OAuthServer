1:�������ã�
beego ���뻷�����ã�����GOPATH����֮��ִ��go get github.com/astaxie/beego
����GOPATH binĿ¼�°�װbee����

GOPATH�Ѿ����ã���srcĿ¼�·�������go���԰�: ����ǰvendorĿ¼�е�github.com�ļ��п�����srcĿ¼��

����Ŀ¼ƽ����ϵ��
liuhy@liuhy ~/linuxshare/work/src $ tree -d -L 1 github.com OAuthServer
github.com
������ astaxie
������ beego
������ gorilla
������ go-sql-driver
������ lib
������ mattn
OAuthServer
������ conf
������ controllers
������ doc
������ models
������ routers
������ tests
������ tmp
������ vendor

2��OAuthServer �������� bee run  
������Ҫ��װmysql���ݿ�
mysql��װ���գ�
http://www.cnblogs.com/bookwed/p/5896619.html

3���߼�����

OAuthServer\conf\app.conf �����ļ���
appname = OAuthServer   ��������
httpport = 8080         rest�ӿڶ˿�
runmode = dev           ���򷢲�ģʽ ���Ĭ��ֵ����������
autorender = false        ���Ĭ��ֵ����������
copyrequestbody = true    ���Ĭ��ֵ����������
EnableDocs = true         ���Ĭ��ֵ����������

dbdriver =  mysql         ��ʾ���ӵ���mysql���ݿ�
dbname = onchain          ���ݿ�����
dbuser = root             �������ݿ���û���
dbpasswd = mobile         �������ݿ������
dbhost = 127.0.0.1        �������ݿ�ĵ�ַ
dbport = 3306             ���ݿ�˿�
maxidleconn = 5           ���ݿ����ӳ���С������
maxopenconn = 5           ���ݿ����ӳ����������
sqldebug = true           �Ƿ���ִ��sql��־ trueΪ����
verifytime = 300000       ����tokenʱ��Ӧ�ó���ͱ���ʱ��֮������ֵ ��λ��
expiretime = 30           token����ʱ�� ��λ��


bee run ����OAuthServer�����
����ʹ��postman���ߺ����psoter���в��ԣ�

1������appid �� appkey �ӿ�

���磺ʹ��poster�������½ӿ�
http://10.0.1.45:8080/v1/app/regInfo

����json��ʽ�������£�
{
  "code": "r0000",
  "desc": "success",
  "msg": {
    "appid": "72592fd77f0a1153bdbcfd97b6ce36ab",
    "appkey": "5b722d7b5029ae69cf4894425a5d7203"
  }
}

���Խ���mysql �����в鿴���ݣ�
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

����registime������appid appkeyʱ�̵�UTCʱ��
tokencount��ʾ���appʹ��token�Ĵ�����ÿ��Ӧ��app����������5��token��

2������token�ӿ�

���磺ʹ��poster�������½ӿ�
http://10.0.1.45:8080/v1/app/getToken
ͬʱ���ݲ������£�
{
  "param": {
    "appid": "72592fd77f0a1153bdbcfd97b6ce36ab",
    "timestamp": "2017-04-25T15:33:51+08:00",
    "random": "123456"
  },
  "signature": "b385dbf71786e554ad2754832404d6d1"
}

˵����
appid��Ӧ�ó���֮ǰ�����appid
timestamp��Ӧ�ó�����utcʱ��
random��6λ�������ַ���
signature�ǽ�appid��timestamp��randomƴ��֮���md5ժҪ


���ص�json�������£�
{
  "code": "r0000",
  "desc": "success",
  "msg": {
    "access_token": "ac08d0051c4d25aaab58d003e0397ffe",
    "expire_time": 30
  }
}

˵��:
access_token ��ʾ���ɵ�token�����
expire_time�Ǳ�ʾtoken��ʧЧʱ��

���ݱ�token_info_tab��¼���ɵ�token��Ϣ��

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


˵�������appid��72592fd77f0a1153bdbcfd97b6ce36ab������5��token ��app_reg_tab�ֶ�tokencountΪ5
��72592fd77f0a1153bdbcfd97b6ce36ab�����5��tokenû��ʧЧ������£��ٴ������򷵻�ʧ�����ݣ�
{
  "code": "r1003",
  "desc": "token overrun appid application, at the same time can only apply for 5 a valid token by default",
  "msg": ""
}





