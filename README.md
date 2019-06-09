# go-sendmail
golang定时发送邮件

#### 使用的第三方包

`gopkg.in/gomail.v2`

`gopkg.in/yaml.v2`

`github.com/robfig/cron`

#### 配置文件

`数据库配置`

```
db:
  db_host: 127.0.0.1
  db_user: root
  db_pass: 123456
  db_port: 3306
  db_name: test
```

`发邮件配置`

```
mail_host: smtp.qq.com
mail_port: 587
mail_user: example@qq.com
mail_pass: qqq
```

`接收用户账号`

```
545922113=545922113@qq.com
703362961=703362961@qq.com
.
.
.

```

#### 测试文件

![001](https://raw.githubusercontent.com/hanquanding/go-sendmail/master/img/001.jpg)


![002](https://raw.githubusercontent.com/hanquanding/go-sendmail/master/img/002.jpg)
