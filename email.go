package main

import (
	"errors"
	"fmt"
	"net/smtp"
	"os"
	"strings"
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

	auth := NewLoginAuth(username, password)
	addr := "smtp-mail.outlook.com:587"

	headers := fmt.Sprintf("From: %s\r\nTo: %s\r\n", username, strings.Join(mailingList, ","))
	message = headers + message

	err := smtp.SendMail(
		addr,
		auth,
		username,
		mailingList,
		[]byte(message),
	)
	if err != nil {
		return err
	}
	return nil
}
