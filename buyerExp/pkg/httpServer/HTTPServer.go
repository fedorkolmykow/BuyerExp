package httpServer

import (
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"

	api "github.com/fedorkolmykow/avitoexp/pkg/api"
)

type service interface {
	Subscribe(Req *api.SubscribeReq) (Resp *api.SubscribeResp, err error)
	Confirm(Req *api.ConfirmReq) (Resp *api.ConfirmResp, err error)
}

type server struct {
	svc service
}

func (s *server) HandleConfirm(w http.ResponseWriter, r *http.Request){
	hash := r.FormValue("hash")
	if hash == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}
	req := &api.ConfirmReq{Hash: hash}
	resp, err := s.svc.Confirm(req)
	if err != nil{
		log.Warn(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	body, err := resp.MarshalJSON()
	if err != nil{
		log.Warn(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(body)
	if err != nil {
		log.Warn(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *server) HandleSubscribe(w http.ResponseWriter, r *http.Request){
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Warn(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req := &api.SubscribeReq{}
	err = req.UnmarshalJSON(body)
	if err != nil{
		log.Warn(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Trace("Received data: " + fmt.Sprintf("%+v", req))
	resp, err := s.svc.Subscribe(req)
	if err != nil{
		log.Warn(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	body, err = resp.MarshalJSON()
	if err != nil{
		log.Warn(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(body)
	if err != nil {
		log.Warn(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func NewHTTPServer(svc service) (httpServer *mux.Router) {
	router := mux.NewRouter()
    s := &server{svc: svc}
	router.HandleFunc(api.Subscribe, s.HandleSubscribe).
		Methods("POST")
	router.HandleFunc(api.Сonfirmation, s.HandleConfirm).
		Methods("GET")
	return router
}