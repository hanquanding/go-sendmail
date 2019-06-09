package main

import (
	"log"
	"testing"

	"io/ioutil"

	"gopkg.in/gomail.v2"
	"gopkg.in/yaml.v2"
)

const CONFIG_FILE = "conf/mail.yml"

type ConfMail struct {
	Mail_host string
	Mail_port int
	Mail_user string
	Mail_pass string
}

var conf *ConfMail

func init() {
	conf = new(ConfMail)
	err := loadYml(CONFIG_FILE, &conf)
	if err != nil {
		log.Printf("%s\n", "Cannot load configuration file")
		panic(err)
	}
	log.Printf("%v\n", conf)
}

func loadYml(path string, t interface{}) error {
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

func TestSendMail(t *testing.T) {
	m := gomail.NewMessage()
	m.SetHeader("From", "545922113@qq.com")
	m.SetHeader("To", "hanquanding@163.com", "703362961@qq.com")
	m.SetAddressHeader("Cc", "545922113@qq.com", "hqd")
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", "Hello <b>hanquanding</b> and <i>703362961</i>!")
	m.Attach("avatar.jpeg")

	d := gomail.NewDialer(conf.Mail_host, conf.Mail_port, conf.Mail_user, conf.Mail_pass)

	// Send the email to hanquanding@163.com, 703362961@qq.com and 545922113@qq.com.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
