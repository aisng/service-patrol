package main

import (
	"errors"
	"fmt"
	"net/smtp"
	"os"
)

type LoginAuth struct {
	username, password string
}

func NewLoginAuth(username, password string) smtp.Auth {
	return &LoginAuth{username, password}
}

func (a *LoginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *LoginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("unkown fromServer")
		}
	}
	return nil, nil
}

func SendMail(mailingList []string, message string) error {
	username := os.Getenv("SPMAILUSERNAME")
	password := os.Getenv("SPMAILPASSWORD")
	// username := os.Getenv("SPMAILUSERNAMEG")
	// password := os.Getenv("SPMAILTOKEN")

	auth := NewLoginAuth(username, password)

	// client, err := smtp.Dial("smtp-mail.outlook.com:587")
	// if err != nil {
	// 	return err
	// }
	// client.Auth(auth)
	server := "smtp-mail.outlook.com:587"
	// server := "smtp.gmail.com:587"

	receiversList := "To: "
	for i, receiver := range mailingList {
		receiversList = receiversList + receiver
		if i != len(mailingList)-1 {
			receiversList += ", "
		}
	}
	// fmt.Println(receiversList)

	message = fmt.Sprintf("From: %s\r\n"+receiversList+"\r\n"+message, username)

	fmt.Println(message)
	err := smtp.SendMail(
		server,
		auth,
		username,
		mailingList,
		[]byte(message),
	)
	// fmt.Println(message)
	if err != nil {
		return err
	}
	return nil
}
