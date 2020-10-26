package worker

import (

	"github.com/fedorkolmykow/avitoexp/pkg/api"
	log "github.com/sirupsen/logrus"
)

type parser interface{
	ParsePrices(notices []api.Notice) (ChangedNotices []api.Notice, err error)
}

type dbClient interface{
	SelectSubscriptions(notices []api.Notice)(subs []api.Subscription, err error)
	UpdateNoticesPrice(notices []api.Notice) (err error)
	SelectAllActiveNotices() (notices []api.Notice, err error)
}

type smtpClient interface{
	SendMailsWithNewPrices(subs []api.Subscription)
}

type worker struct{
	db dbClient
	smtp smtpClient
	par parser
}

type Worker interface {
	SendChanges() (err error)
}

func (w *worker) SendChanges() (err error){
	notices, err  := w.db.SelectAllActiveNotices()
	if err != nil{
		return
	}
	log.Trace(notices)
	notices, err = w.par.ParsePrices(notices)
	err = w.db.UpdateNoticesPrice(notices)
	subs, err := w.db.SelectSubscriptions(notices)
	log.Trace(subs)
	w.smtp.SendMailsWithNewPrices(subs)
	return
}

func NewWorker(db dbClient, smtp smtpClient, par parser) Worker{
	w := &worker{
		db: db,
		smtp: smtp,
		par: par,
	}
	return w
}