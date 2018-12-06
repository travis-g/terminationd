package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const validInstanceAction = `{"action":"terminate","time":"2017-09-18T08:22:00Z"}`

func initTestServer(path string, resp string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI != path {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		w.Write([]byte(resp))
	}))
}

func TestIsTermiating(t *testing.T) {
	server := initTestServer(
		"/latest/meta-data/spot/instance-action",
		validInstanceAction,
	)
	defer server.Close()

	instanceAction, err := GetInstanceAction(server.URL + "/latest/meta-data/spot/instance-action")
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}

	t.Logf("data %v", instanceAction)

	if e, a := true, instanceAction.IsTerminating(); e != a {
		t.Errorf("expect %v, got %v", e, a)
	}
}
