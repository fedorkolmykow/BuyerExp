package httpServer
/*
import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	m "github.com/fedorkolmykow/avitoauto/pkg/api"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const(
	saveURL = iota
	redirect
)

type correctService struct{
}

type errorService struct{
}

type TestCase struct {
	Vars   map[string]string
	Req     []byte
	Resp    string
	Status  int
	S       server
	Handle  int
}

func TestHandles(t *testing.T){
	cases := []TestCase{
		{
			Vars:        map[string]string{},
			Req:          []byte(`{"url":"200","custom_key":"My First"}`),
			Resp:         `{"key":"1w34r"}`,
			Status:       http.StatusOK,
			S:            server{svc: &correctService{}},
			Handle:       saveURL,
		},
		{
			Vars:        map[string]string{},
			Req:          []byte(`Here is error`),
			Resp:         ``,
			Status:       http.StatusBadRequest,
			S:            server{svc: &correctService{}},
			Handle:       saveURL,
		},
		{
			Vars:        map[string]string{},
			Req:          []byte(`{"url":"200","custom_key":"My First"}`),
			Resp:         ``,
			Status:       http.StatusInternalServerError,
			S:            server{svc: &errorService{}},
			Handle:       saveURL,
		},
		{
			Vars:        map[string]string{"key":"1w34r"},
			Req:          []byte(``),
			Resp:         ``,
			Status:       http.StatusSeeOther,
			S:            server{svc: &correctService{}},
			Handle:       redirect,
		},
		{
			Vars:        map[string]string{"key":"1w34r"},
			Req:          []byte(``),
			Resp:         ``,
			Status:       http.StatusInternalServerError,
			S:            server{svc: &errorService{}},
			Handle:       redirect,
		},
	}
	log.SetLevel(log.FatalLevel)
	for num, c := range cases{
		req := httptest.NewRequest(
			"NotImportant",
			"http://localhost",
			bytes.NewBuffer(c.Req),
		)
		req = mux.SetURLVars(req, c.Vars)
		w := httptest.NewRecorder()
		switch c.Handle {
		case saveURL:     c.S.HandleSaveURL(w, req)
		case redirect:    c.S.HandleRedirect(w, req)
	}

		if w.Result().StatusCode != c.Status{
			t.Errorf("[%d] unexpected status: %d, expected: %d",num, w.Result().StatusCode,  c.Status)
		}
		if c.Status == http.StatusOK{
			if c.Resp != w.Body.String(){
				t.Errorf("[%d] unexpected result:\n%s\nexpected:\n%s ", num, w.Body.String(), c.Resp)
			}
		}
	}
}

//correctService
func (s *correctService)  SaveURL(Req *m.SaveURLReq) (Resp *m.SaveURLResp, err error){
	return &m.SaveURLResp{Key: "1w34r"}, nil
}

func (s *correctService) GetURL(Req *m.RedirectReq) (Resp *m.RedirectResp, err error){
	return &m.RedirectResp{OriginalURL: "localhost/"}, nil
}

//errorService
func (s *errorService) SaveURL(Req *m.SaveURLReq) (Resp *m.SaveURLResp, err error) {
	return nil, errors.New("test error")
}

func (s *errorService) GetURL(Req *m.RedirectReq) (Resp *m.RedirectResp, err error){
	return nil, errors.New("test error")
}*/