package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

var (
	config map[string]string
	dialer *gomail.Dialer
)

// 加载配置
func loadConfig(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return err
	}
	config = make(map[string]string)
	err = json.Unmarshal(content, &config)
	if err != nil {
		fmt.Println(err)
		return err
	}
	port, err := strconv.Atoi(config["port"])
	if err != nil {
		return err
	}
	dialer = gomail.NewDialer(config["host"], port, config["username"], config["password"])
	return nil
}

// 发送邮件：简单方式
func sendMail(subject string, content string, to string, cc string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", config["from"])

	if len(to) == 0 {
		return errors.New("收件人不能为空")
	}
	m.SetHeader("To", to)
	if len(cc) > 0 {
		m.SetHeader("Cc", cc)
	}

	m.SetHeader("Subject", subject)
	m.SetBody("text/html", content)
	return dialer.DialAndSend(m)
}

// 发送邮件：复杂方式
func sendMsg(msg *gomail.Message) error {
	return dialer.DialAndSend(msg)
}

func main() {
	args := os.Args[1:]
	err := loadConfig("./mailto.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(args) == 4 {
		sendMail(args[0], args[1], args[2], args[3])
	} else {
		sendMail(args[0], args[1], args[2], "")
	}
}