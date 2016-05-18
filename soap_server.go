package gosoap

import (
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/achiku/xml"
)

// ProcessBRequest struct
type ProcessBRequest struct {
	XMLName   xml.Name `xml:"http://example.com/ns ProcessBRequest"`
	RequestID string   `xml:"RequestId"`
}

// ProcessARequest struct
type ProcessARequest struct {
	XMLName   xml.Name `xml:"http://example.com/ns ProcessARequest"`
	RequestID string   `xml:"RequestId"`
}

// ProcessAResponse struct
type ProcessAResponse struct {
	*AbstractResponse
	XMLName xml.Name `xml:"http://example.com/ns ProcessAResponse"`
	ID      string   `xml:"Id,omitifempty"`
	Process string   `xml:"Process,omitifempty"`
}

// ProcessBResponse struct
type ProcessBResponse struct {
	*AbstractResponse
	XMLName xml.Name `xml:"http://example.com/ns ProcessBResponse"`
	ID      string   `xml:"Id,omitifempty"`
	Process string   `xml:"Process,omitifempty"`
	Amount  string   `xml:"Amount,omitifempty"`
}

func processA() ProcessAResponse {
	return ProcessAResponse{
		AbstractResponse: &AbstractResponse{
			Code:   "200",
			Detail: "success",
		},
		ID:      "100",
		Process: "ProcessAResponse",
	}
}

func processB() ProcessBResponse {
	return ProcessBResponse{
		AbstractResponse: &AbstractResponse{
			Code:   "200",
			Detail: "success",
		},
		ID:      "100",
		Process: "ProcessBResponse",
		Amount:  "10000",
	}
}

func soapActionHandler(w http.ResponseWriter, r *http.Request) {
	soapAction := r.Header.Get("soapAction")
	var res interface{}
	switch soapAction {
	case "processA":
		res = processA()
	case "processB":
		res = processB()
	default:
		res = nil
	}
	v := SOAPEnvelope{
		Body: SOAPBody{
			Content: res,
		},
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/xml")
	x, err := xml.MarshalIndent(v, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(x)
	return
}

func soapBodyHandler(w http.ResponseWriter, r *http.Request) {
	rawbody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	a := regexp.MustCompile(`<ProcessARequest xmlns="http://example.com/ns">`)
	b := regexp.MustCompile(`<ProcessBRequest xmlns="http://example.com/ns">`)

	var res interface{}
	if a.MatchString(string(rawbody)) {
		res = processA()
	} else if b.MatchString(string(rawbody)) {
		res = processB()
	} else {
		res = nil
	}
	v := SOAPEnvelope{
		Body: SOAPBody{
			Content: res,
		},
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/xml")
	x, err := xml.MarshalIndent(v, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(x)
	return
}

// NewSOAPMux return SOAP server mux
func NewSOAPMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/dispatch/soapaction", soapActionHandler)
	mux.HandleFunc("/dispatch/soapbody", soapBodyHandler)
	return mux
}

// NewSOAPServer create i2c mock server
func NewSOAPServer(port string) *http.Server {
	mux := NewSOAPMux()
	server := &http.Server{
		Handler: mux,
		Addr:    "localhost:" + port,
	}
	return server
}
