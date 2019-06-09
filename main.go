package main

import (
	"bufio"
	"crypto/tls"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"html/template"

	"time"

	"bytes"

	"github.com/robfig/cron"
	"gopkg.in/gomail.v2"
	"gopkg.in/yaml.v2"
)

const CONFIG_PATH_DB = "conf/db.yml"
const CONFIG_PATH_MAIl = "conf/mail.yml"

type ConfigDB struct {
	Db struct {
		Db_host string
		Db_user string
		Db_pass string
		Db_port string
		Db_name string
	}
}

type ConfigMail struct {
	Mail_host string
	Mail_port int
	Mail_user string
	Mail_pass string
}

var (
	confDb   *ConfigDB
	confMail *ConfigMail
)

var subject string = "Hello! this is a test!!"

var tpl = template.Must(template.ParseFiles("templates/tpl.html"))

func init() {
	var err error

	confDb = new(ConfigDB)
	confMail = new(ConfigMail)

	err = LoadYml(CONFIG_PATH_MAIl, &confMail)
	if err != nil {
		log.Printf("%s\n", "Cannot load configuration file")
		panic(err)
	}
	log.Printf("%v\n", confMail)

	err = LoadYml(CONFIG_PATH_DB, &confDb)
	if err != nil {
		log.Printf("%s\n", "Cannot load configuration file")
		panic(err)
	}
	log.Printf("%v\n", confDb)
}

func LoadYml(path string, t interface{}) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(b, t)
	if err != nil {
		return err
	}
	return nil
}

func LoadConf(path string) map[string]string {
	confMap := make(map[string]string)

	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		panic(err)
	}

	r := bufio.NewReader(file)
	// 循环读取
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		// 过滤两端空格
		str := strings.TrimSpace(string(line))

		// 等号（=）的位置，没有找到跳过
		index := strings.Index(str, "=")
		if index < 0 {
			continue
		}

		// 等号（=）左边的值
		name := strings.TrimSpace(str[:index])
		if len(name) == 0 {
			continue
		}

		// 等号（=）右边的值
		address := strings.TrimSpace(str[index+1:])
		if len(address) == 0 {
			continue
		}

		confMap[name] = address
	}
	return confMap

}

func dateFormat(t time.Time) string {
	layout := "2006-01-02 15:04:05"
	return t.Format(layout)
}

var funcMap = template.FuncMap{
	"dateFormat": dateFormat,
}

func sendMail() {
	accountMap := LoadConf("conf/account.txt")
	// fmt.Println(accountMap)
	data := struct {
		Date      string
		SendDate  string
		TotalUser int32
		NumUser   int32
		Volume    int32
	}{
		Date:      dateFormat(time.Now()),
		SendDate:  time.Now().Format("2006年01月02日"),
		TotalUser: 222,
		NumUser:   22,
		Volume:    22,
	}

	// 存储已经执行的模板输出
	var buffer bytes.Buffer
	err := tpl.Execute(&buffer, data)
	if err != nil {
		log.Println(err)
		return
	}

	body := buffer.String()
	// fmt.Println(buffer.String())

	i := 0
	m := gomail.NewMessage(gomail.SetCharset("UTF-8"))
	for name, address := range accountMap {
		i++
		m.SetHeader("From", confMail.Mail_user)
		m.SetAddressHeader("To", address, name)
		m.SetAddressHeader("Cc", confMail.Mail_user, "hqd") // 抄送
		m.SetHeader("Subject", subject)
		m.SetBody("text/html", body)
		m.Attach("avatar.jpeg")

		d := gomail.NewDialer(confMail.Mail_host, confMail.Mail_port, confMail.Mail_user, confMail.Mail_pass)
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

		// send the email
		if err := d.DialAndSend(m); err != nil {
			log.Printf("Could not send email to %q: %v", address, err)
			panic(err)
		}
		m.Reset()
		log.Println("cron running:", i)
	}
}

func main() {
	// 每10秒发送一次邮件
	spec := "*/10 * * * * *"
	c := cron.New()
	c.AddFunc(spec, sendMail)
	c.Start()
	select {}

}
