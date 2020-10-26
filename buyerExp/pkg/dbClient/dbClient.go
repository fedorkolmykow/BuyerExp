package dbClient

import (
	"database/sql"
	"os"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"

	api "github.com/fedorkolmykow/avitoexp/pkg/api"
)

const(
	SelectUserExistMail = `SELECT EXISTS(SELECT user_id FROM Users WHERE mail=$1) ;`
	SelectUserExistHash = `SELECT EXISTS(SELECT user_id FROM Users WHERE hash=$1) ;`
	SelectNoticeExist = `SELECT EXISTS(SELECT notice_id FROM Notices WHERE url=$1) ;`
	SelectConfirmed = `SELECT confirmed From Users WHERE mail=$1;`
	InsertUser = `INSERT INTO Users (mail, confirmed, hash) VALUES ($1, $2, $3);`
	InsertNotice = `INSERT INTO Notices (url, price) VALUES ($1, $2) RETURNING notice_id;`
	InsertSubscription = `INSERT INTO Subscription (notice_id, user_id) VALUES ($1, $2);`
	DeleteUser = `DELETE FROM Users WHERE mail=$1;`
	SelectUser = `SELECT user_id From Users WHERE mail=$1;`
	SelectNotice = `SELECT notice_id From Notices WHERE url=$1;`
	UpdateConfirm = `Update users SET confirmed = $1 WHERE hash = $2;`
	SelectAllActiveNotices = `SELECT DISTINCT url, price FROM Notices 
								INNER JOIN Subscription 
								ON Notices.notice_id = Subscription.notice_id
								INNER JOIN Users
								ON Subscription.user_id = Users.user_id
								WHERE Users.confirmed = true;`
	SelectSubscribedUsers = `SELECT Users.mail FROM Users
								INNER JOIN Subscription
								ON Subscription.user_id = Users.user_id
								INNER JOIN Notices 
								ON Notices.notice_id = Subscription.notice_id
								WHERE Notices.url = $1;`
	UpdatePrice = `Update Notices SET price = $1 WHERE url = $2;`
)

type DbClient interface{
	CheckUser(mail string) (registered bool, err error)
	AddUser(user *api.User) (err error)
	AddSubscription(Req *api.SubscribeReq) (Resp *api.SubscribeResp, err error)
	ConfirmMail(user *api.User)(Resp *api.ConfirmResp, err error)
	SelectSubscriptions(notices []api.Notice)(subs []api.Subscription, err error)
	UpdateNoticesPrice(notices []api.Notice) (err error)
	SelectAllActiveNotices() (notices []api.Notice, err error)
	Shutdown() error
}

type dbClient struct{
    db *sqlx.DB
}

func (d *dbClient) SelectAllActiveNotices() (notices []api.Notice, err error) {
	err = d.db.Select(&notices, SelectAllActiveNotices)
	return
}

func (d *dbClient) SelectSubscriptions(notices []api.Notice)(subs []api.Subscription, err error){
	subs = make([]api.Subscription, len(notices))
	tx, err := d.db.Beginx()
	if err != nil{
		return
	}
	for _, n := range notices{
		var rows *sql.Rows
		rows, err = tx.Query(SelectSubscribedUsers, n.URL)
		if err != nil{
			err = rollAndErr(tx, err)
			return
		}
		for rows.Next() {
			user := api.User{}
			err = rows.Scan(&user.Mail)
			if err != nil{
				err = rollAndErr(tx, err)
				return
			}
			sub := api.Subscription{
				Notice: n,
				User:   user,
			}
			subs = append(subs, sub)
		}
	}
	err= tx.Commit()
	return
}

func (d *dbClient) UpdateNoticesPrice(notices []api.Notice) (err error){
	tx, err := d.db.Beginx()
	if err != nil{
		return
	}
	for _, n := range notices{
		_, err = tx.Exec(UpdatePrice, n.Price, n.URL)
		if err != nil{
			err = rollAndErr(tx, err)
			return
		}
	}
	err= tx.Commit()
	return
}

func (d *dbClient) CheckUser(mail string) (registered bool, err error) {
	tx, err := d.db.Beginx()
	if err != nil{
		return
	}
	err = tx.QueryRow(SelectUserExistMail, mail).Scan(&registered)
	if err != nil{
		err = rollAndErr(tx, err)
		return
	}
	if registered{
		err = tx.QueryRow(SelectConfirmed, mail).Scan(&registered)
		log.Println(registered)
		if err != nil{
			err = rollAndErr(tx, err)
			return
		}
	}
	err = tx.Commit()
	return
}

func (d *dbClient) AddUser(user *api.User) (err error) {
	var exist bool
	tx, err := d.db.Beginx()
	if err != nil{
		return
	}
	err = tx.QueryRow(SelectUserExistMail, user.Mail).Scan(&exist)
	if err != nil{
		err = rollAndErr(tx, err)
		return
	}
	if exist{
		_, err = tx.Exec(DeleteUser, user.Mail)
		if err != nil{
			err = rollAndErr(tx, err)
			return
		}
	}
	_, err = tx.Exec(InsertUser, user.Mail, false, user.Hash)
	if err != nil{
		err = rollAndErr(tx, err)
		return
	}
	err = tx.Commit()
	return
}
func (d *dbClient) AddSubscription(Req *api.SubscribeReq) (Resp *api.SubscribeResp, err error) {
	var userId, noticeId int
	var exist bool
	tx, err := d.db.Beginx()
	if err != nil{
		return
	}
	err = tx.QueryRow(SelectUser, Req.Mail).Scan(&userId)
	if err != nil{
		err = rollAndErr(tx, err)
		return
	}
	err = tx.QueryRow(SelectNoticeExist, Req.NoticeURL).Scan(&exist)
	if exist{
		err = tx.QueryRow(SelectNotice, Req.NoticeURL).Scan(&noticeId)
	} else{
		err = tx.QueryRow(InsertNotice, Req.NoticeURL, 0).Scan(&noticeId)
	}
	if err != nil{
		err = rollAndErr(tx, err)
		return
	}
	_, err = tx.Exec(InsertSubscription, noticeId, userId)
	if err != nil{
		err = rollAndErr(tx, err)
		return
	}
	err = tx.Commit()
	Resp = &api.SubscribeResp{
		Id: userId,
	}
	return
}

func (d *dbClient) 	ConfirmMail(user *api.User)(Resp *api.ConfirmResp, err error) {
	var exist bool
	Resp = &api.ConfirmResp{}
	tx, err := d.db.Beginx()
	if err != nil{
		return
	}
	err = tx.QueryRow(SelectUserExistHash, user.Hash).Scan(&exist)
	if err != nil{
		err = rollAndErr(tx, err)
		return
	}
	if exist{
		_, err = tx.Exec(UpdateConfirm, true, user.Hash)
		if err != nil{
			err = rollAndErr(tx, err)
			return
		}
		Resp.Message = "Success!"
	} else{
		Resp.Message = "Link too old. Try subscribe again."
	}
	err = tx.Commit()
	return
}

func (d *dbClient) Shutdown() error{
	return d.db.Close()
}

func rollAndErr(tx *sqlx.Tx, err error) error{
	errRoll := tx.Rollback()
	if errRoll != nil{
		return errRoll
	}
	return err
}

func NewDbClient() DbClient{
	db, err := sqlx.Connect("pgx", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatal(err)
	}
	//db.SetMaxIdleConns(n int)
	//db.SetMaxOpenConns(n int)
	return &dbClient{db: db}
}

