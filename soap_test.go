package xmlnstest

import (
	"testing"

	"github.com/achiku/xml"
)

// SOAPEnvelope envelope
type SOAPEnvelope struct {
	XMLName xml.Name    `xml:"http://schemas.xmlsoap.org/soap/envelope Envelope"`
	Header  *SOAPHeader `xml:",omitempty"`
	Body    SOAPBody    `xml:",omitempty"`
}

// SOAPHeader header
type SOAPHeader struct {
	XMLName xml.Name    `xml:"http://schemas.xmlsoap.org/soap/envelope Header"`
	Content interface{} `xml:",omitempty"`
}

// SOAPBody body
type SOAPBody struct {
	XMLName xml.Name    `xml:"http://schemas.xmlsoap.org/soap/envelope Body"`
	Fault   *SOAPFault  `xml:",omitempty"`
	Content interface{} `xml:",omitempty"`
}

// SOAPFault fault
type SOAPFault struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope Fault"`
	Code    string   `xml:"faultcode,omitempty"`
	String  string   `xml:"faultstring,omitempty"`
	Actor   string   `xml:"faultactor,omitempty"`
	Detail  string   `xml:"detail,omitempty"`
}

type CommonResponse struct {
	Code   string `xml:"Code,omitempty"`
	Detail string `xml:"Detail,omitempty"`
}

type ConcreteResponse struct {
	*CommonResponse
	XMLName           xml.Name `xml:"http://example.com/ns ConcreteResponse"`
	AdditionalMessage string   `xml:"AdditionalMessage,omitempty"`
}

// UnmarshalXML unmarshal SOAPBody
func (b *SOAPBody) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if b.Content == nil {
		return xml.UnmarshalError("Content must be a pointer to a struct")
	}
	var (
		token    xml.Token
		err      error
		consumed bool
	)
Loop:
	for {
		if token, err = d.Token(); err != nil {
			return err
		}
		if token == nil {
			break
		}
		envelopeNameSpace := "http://schemas.xmlsoap.org/soap/envelope/"
		switch se := token.(type) {
		case xml.StartElement:
			if consumed {
				return xml.UnmarshalError(
					"Found multiple elements inside SOAP body; not wrapped-document/literal WS-I compliant")
			} else if se.Name.Space == envelopeNameSpace && se.Name.Local == "Fault" {
				b.Fault = &SOAPFault{}
				b.Content = nil
				err = d.DecodeElement(b.Fault, &se)
				if err != nil {
					return err
				}
				consumed = true
			} else {
				if err = d.DecodeElement(b.Content, &se); err != nil {
					return err
				}
				consumed = true
			}
		case xml.EndElement:
			break Loop
		}
	}
	return nil
}

// Name struct
type Name struct {
	XMLName xml.Name `xml:"http://example.com/ns Name"`
	First   string   `xml:"First,omitempty"`
	Last    string   `xml:"Last,omitempty"`
}

type Auth struct {
	XMLName xml.Name `xml:"http://example.com/ns Auth"`
	UserID  string   `xml:"UserID"`
	Pass    string   `xml:"Pass"`
}

// Person struct
type Person struct {
	XMLName xml.Name `xml:"http://example.com/ns Person"`
	ID      int      `xml:"Id,omitempty"`
	Name    *Name
	Age     int `xml:"Age,omitempty"`
}

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
				CommonResponse: &CommonResponse{
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
