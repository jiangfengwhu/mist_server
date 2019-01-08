package main

import (
	"fmt"
	"log"
	"net/smtp"
)

func sendActiveMail(register string, link string) error {
	auth := smtp.PlainAuth(
		"",
		"jf13163291713@163.com",
		"saber931228",
		"smtp.163.com",
	)
	msg := fmt.Sprintf("To: %s\r\nFrom: %s\r\nSubject: %s\r\nContent-Type: %s; charset=UTF-8\r\n\r\n%s",
		register, "jf13163291713@163.com", "Mist激活邮件", "text/html", link)
	err := smtp.SendMail(
		"smtp.163.com:25",
		auth,
		"jf13163291713@163.com",
		[]string{register},
		[]byte(msg),
	)
	if err != nil {
		log.Println(err)
	}
	return err
}
