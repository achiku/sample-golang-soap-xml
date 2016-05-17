package gosoap

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSOAPServerBySOAPAction(t *testing.T) {
	mux := NewSOAPMux()
	ts := httptest.NewServer(mux)
	defer ts.Close()
	url := ts.URL + "/dispatch/soapaction"
	actions := []string{"processA", "processB"}
	for _, action := range actions {
		req, err := http.NewRequest("POST", url, nil)
		req.Header.Set("soapAction", action)
		if err != nil {
			t.Fatal(err)
		}
		c := &http.Client{}
		resp, err := c.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%+v", resp)
	}
}
