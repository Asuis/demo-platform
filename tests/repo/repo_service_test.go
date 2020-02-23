package repo

import (
	"bytes"
	"demo-platform/conf"
	"demo-platform/model/db"
	"demo-platform/services/repo"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateRepoService(t *testing.T) {
	_ = db.SetupDatabase()

	url := "http://127.0.0.1:8000/v1/repo/create"

	contentType := "application/json;charset=utf-8"

	form := &gin.H{
		"Name": "asuis/test3",
		"Description": "test3 repository",
		"IsPrivate":true,
		"UseCustomAvatar":true,
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
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhYyI6MSwibiI6ImFzdWlzIn0.Qgu1NRBvgnVoSe4w16x-TwcrKDCBFsqt8c-qmSFyZ14")
	router.ServeHTTP(w, req)
	
	log.Println("content:", w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)


	//query
	re, err := repo.Info("asuis/test3", &db.User{
		Id: 1,
	})
	assert.Equal(t, err, nil)
	assert.Equal(t, re.Name, "asuis/test3")
	r, err:= repo.SearchDir("asuis/test3", "/")
	assert.Equal(t, err, nil)

	for _, value := range *r {
		log.Printf("ID: %s, Name: %s", value.ID, value.Name())
	}


}

func TestSearchRepoDirService(t *testing.T) {
	conf.SetupRouter()
	_ = db.SetupDatabase()

	url := "http://127.0.0.1:8000/v1/repo/file/asuis/test3/"

	router := conf.SetupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhYyI6MSwibiI6ImFzdWlzIn0.Qgu1NRBvgnVoSe4w16x-TwcrKDCBFsqt8c-qmSFyZ14")

	router.ServeHTTP(w, req)

	log.Println("content:", w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)

}

func TestDeleteRepoService(t *testing.T) {
	_ = db.SetupDatabase()

	url := "http://127.0.0.1:8000/v1/repo/asuis/test3"

	router := conf.SetupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("DELETE", url, nil)

	router.ServeHTTP(w, req)

	log.Println("content:", w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)


	//query
	re, err := repo.Info("asuis/test3", &db.User{
		Id: 1,
	})
	assert.Equal(t, err, nil)
	assert.Equal(t, re, nil)
	r, err:= repo.SearchDir("asuis/test3", "/")
	assert.Equal(t, err, nil)
	assert.Equal(t, r, nil)
}

func TestGetRepoInfoService(t *testing.T) {

}

func TestGetBranchRepoInfoService(t *testing.T) {

}

func TestGetCommitRepoInfoService(t *testing.T) {

}

func TestGetFileInfoRepoService(t *testing.T) {

}

func TestGetDirInfoRepoService(t *testing.T) {

}