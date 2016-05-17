package gosoap

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
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

func TestSOAPServerByRequestBody(t *testing.T) {
	mux := NewSOAPMux()
	ts := httptest.NewServer(mux)
	defer ts.Close()
	url := ts.URL + "/dispatch/soapbody"
	contents := []interface{}{
		ProcessARequest{RequestID: "request-a-id"},
		ProcessBRequest{RequestID: "request-b-id"},
	}

	for _, c := range contents {
		envelope := SOAPEnvelope{
			Body: SOAPBody{
				Content: c,
			},
		}
		buf, err := xml.MarshalIndent(envelope, "", "")
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(buf))
		if err != nil {
			t.Fatal(err)
		}
		c := &http.Client{}
		resp, err := c.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Response:\n%s", body)
	}
}
