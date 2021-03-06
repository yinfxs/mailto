package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

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
	if len(config["contentType"]) == 0 {
		config["contentType"] = "text/plain"
	}

	if err != nil {
		return err
	}
	dialer = gomail.NewDialer(config["host"], port, config["username"], config["password"])
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return nil
}

// 发送邮件：简单方式
func sendMail(subject string, content string, to string, cc string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", config["from"])

	if len(to) == 0 {
		return errors.New("收件人不能为空")
	}
	toArray := strings.Split(to, ",")
	m.SetHeader("To", toArray...)
	if len(cc) > 0 {
		ccArray := strings.Split(cc, ",")
		m.SetHeader("Cc", ccArray...)
	}

	m.SetHeader("Subject", subject)
	m.SetBody(config["contentType"], content)
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
		panic(err)
	}
	if len(args) == 4 {
		err = sendMail(args[0], args[1], args[2], args[3])
	} else {
		err = sendMail(args[0], args[1], args[2], "")
	}
	if err != nil {
		panic(err)
	}
}
