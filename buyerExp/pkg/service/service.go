package service

import (
	"encoding/base64"
	"crypto/sha1"
	"time"

	api "github.com/fedorkolmykow/avitoexp/pkg/api"

)

type Service interface {
	Subscribe(Req *api.SubscribeReq) (Resp *api.SubscribeResp, err error)
	Confirm(user *api.ConfirmReq) (Resp *api.ConfirmResp, err error)
}

type dbClient interface{
	CheckUser(mail string) (registered bool, err error)
	AddUser(user *api.User) (err error)
	AddSubscription(Req *api.SubscribeReq) (Resp *api.SubscribeResp, err error)
	ConfirmMail(user *api.User)(Resp *api.ConfirmResp, err error)
}

type smtpClient interface {
	SendConfirmationMail(user *api.User) (err error)
}

type service struct{
	db dbClient
	smtp smtpClient
}

func (s *service) Subscribe(Req *api.SubscribeReq) (Resp *api.SubscribeResp, err error){
	registered, err:= s.db.CheckUser(Req.Mail)
	if !registered{
		str:= Req.Mail + time.Now().Format(time.RFC822)
		hasher := sha1.New()
		hasher.Write([]byte(str))
		user := &api.User{
			Mail: Req.Mail,
			Hash: hasher.Sum(nil),
		}
		err = s.db.AddUser(user)
		if err != nil{
			return
		}
		err = s.smtp.SendConfirmationMail(user)
		if err != nil{
			return
		}
	}
	Resp, err = s.db.AddSubscription(Req)
	return
}

func (s *service) Confirm(Req *api.ConfirmReq) (Resp *api.ConfirmResp, err error){
	hash, err := base64.URLEncoding.DecodeString(Req.Hash)
	if err != nil{
		return
	}
	user := &api.User{
		Hash: hash,
	}
	return 	s.db.ConfirmMail(user)
}

func NewService(db dbClient, smtp smtpClient) Service{
    svc := &service{
    	db: db,
    	smtp: smtp,
	}
    return svc
}

