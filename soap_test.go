package gosoap

import (
	"testing"

	"github.com/achiku/xml"
)

func TestXMLNameSpace(t *testing.T) {
	p := Name{First: "Akira", Last: "Chiku"}
	buf1, err := xml.MarshalIndent(p, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("\n%s", string(buf1))
}

func TestMarshalSOAPWithNameSpace(t *testing.T) {
	v := SOAPEnvelope{
		Header: &SOAPHeader{
			Content: &Auth{
				UserID: "user",
				Pass:   "pass",
			},
		},
		Body: SOAPBody{
			Content: &Person{
				ID:  1,
				Age: 22,
				Name: &Name{
					Last:  "Mogami",
					First: "Moga",
				},
			},
		},
	}
	buf1, err := xml.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("\n%s", string(buf1))
}

func TestUnmarshalSOAPWtihNameSpace(t *testing.T) {
	testxml := []byte(`
    <soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope" xmlns:ns="http://example.com/ns">
      <soapenv:Body>
        <ns:Person>
          <ns:Id>1</ns:Id>
          <ns:Name>
            <ns:First>Moga</ns:First>
            <ns:Last>Mogami</ns:Last>
          </ns:Name>
          <ns:Age>22</ns:Age>
        </ns:Person>
      </soapenv:Body>
    </soapenv:Envelope>
	`)
	v := SOAPEnvelope{
		Body: SOAPBody{
			Content: &Person{},
		},
	}

	if err := xml.Unmarshal(testxml, &v); err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", v.XMLName)
	t.Logf("%+v", v.Body.Content)
}

func TestAbstructResponse(t *testing.T) {
	v := SOAPEnvelope{
		Body: SOAPBody{
			Content: &ConcreteResponse{
				AbstractResponse: &AbstractResponse{
					Code:   "success",
					Detail: "detail",
				},
				AdditionalMessage: "additional message",
			},
		},
	}
	buf1, err := xml.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("\n%s", string(buf1))
}

func TestUnmarshalAbstructResponse(t *testing.T) {
	testxml := []byte(`
    <Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope">
      <Body>
        <ConcreteResponse xmlns="http://example.com/ns">
          <Code>success</Code>
          <Detail>detail</Detail>
          <AdditionalMessage>additional message</AdditionalMessage>
        </ConcreteResponse>
      </Body>
    </Envelope>
	`)
	v := SOAPEnvelope{
		Body: SOAPBody{
			Content: &ConcreteResponse{},
		},
	}

	if err := xml.Unmarshal(testxml, &v); err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", v.XMLName)

	c := v.Body.Content.(*ConcreteResponse)
	t.Logf("%s", c.AdditionalMessage)
	t.Logf("%+v", c.Detail)
	t.Logf("%+v", c.Code)
}
