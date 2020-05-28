package server

import (
	"encoding/json"
	"github.com/lexbond13/api_core/config"
	"github.com/lexbond13/api_core/services/api/handler"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	router *gin.Engine
	writer *httptest.ResponseRecorder
)

func init() {
	var err error
	router, err = GetRouter(&config.Params{
		FrontendURL: "http://localhost",
		ApiURL: "http://localhost",
		OpenAuthMode: true,
	})
	if err != nil {
		panic(err)
	}

	writer = httptest.NewRecorder()
}

func TestIndexRoute(t *testing.T) {

	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(writer, req)

	if writer.Code != http.StatusOK {
		t.Fatal("wrong code", "expected", http.StatusOK, "got", writer.Code)
	}

	response := &handler.Response{}

	err := json.Unmarshal(writer.Body.Bytes(), &response)
	if err != nil {
		t.Error(err)
	}

	if !response.Status {
		t.Error("expected status", true, "got", response.Status)
	}
}
