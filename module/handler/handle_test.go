package handler

import (
	"fmt"
	"github.com/stretchr/testify"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndex(t *testing.T) {
	expectedBody, _ := ioutil.ReadFile(TemplatesLocation["index"])
	handler := new(IndexHandler)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("http://localhost:8080")
	req, err := http.NewRequest("GET", url, nil)
	testify.Assert.Nil(t, err)

	handler.ServeHTTP(recorder, req)

	testify.Assert.Equal(t, expectedBody, recorder.Body.Bytes())
}
