package user

import (
	"bytes"
	"demo-platform/conf"
	"demo-platform/model/db"
	"demo-platform/services/user"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignIn(t *testing.T) {
	_ = db.SetupDatabase()
	url := "http://127.0.0.1:8000/v1/usr/sign_in"
	contentType := "application/json;charset=utf-8"

	form := &user.Login{
		Account: "asuis",
		Passwd:  "Demo127117",
	}

	b ,err := json.Marshal(form)
	if err != nil {
		log.Println("json format error:", err)
		return
	}

	body := bytes.NewBuffer(b)


	router := conf.SetupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", contentType)
	router.ServeHTTP(w, req)


	log.Println("content:", w.Body.String())
}

func TestRegister(t *testing.T) {
	_ = db.SetupDatabase()
	url := "http://127.0.0.1:8000/v1/usr/sign_up"
	contentType := "application/json;charset=utf-8"

	form := &user.Register{
		Account: "asuis",
		Type:    "0",
		Passwd:  "Demo127117",
		IP:      "0.0.0.0",
	}

	b ,err := json.Marshal(form)
	if err != nil {
		log.Println("json format error:", err)
		return
	}

	body := bytes.NewBuffer(b)


	router := conf.SetupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", contentType)
	router.ServeHTTP(w, req)


	log.Println("content:", w.Body.String())
}

func TestSignedToken(t *testing.T) {

	token, err := user.SignedToken(&db.User{
		Id:   0,
		Salt: rand.Uint64(),
		LoginName: "asuis",
	})

	if err != nil {
		return
	}

	log.Println("content:", string(token))

}

func TestAuthMiddleware(t *testing.T) {
	_ = db.SetupDatabase()
	url := "http://127.0.0.1:8000/v1/admin/auth"
	contentType := "application/json;charset=utf-8"

	router := conf.SetupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhYyI6MSwibiI6ImFzdWlzIn0.Qgu1NRBvgnVoSe4w16x-TwcrKDCBFsqt8c-qmSFyZ14")
	router.ServeHTTP(w, req)


	log.Println("content:", w.Body.String())
}

func TestParseToken(t *testing.T)  {
	sign, _:= user.ParseToken("Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhYyI6MSwibiI6ImFzdWlzIn0.Qgu1NRBvgnVoSe4w16x-TwcrKDCBFsqt8c-qmSFyZ14")
	log.Println(sign)

}