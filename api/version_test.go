package api

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVersion(t *testing.T) {
	wanted := "1.0.0"
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	Version(w, r)

	res := w.Result()
	defer res.Body.Close()
	got, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	if string(got) != "1.0.0" {
		t.Errorf("wanted %s, got %s", wanted, got)
	}
}
