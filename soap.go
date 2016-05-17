package gosoap

import "github.com/achiku/xml"

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

// Auth authorization header
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

// AbstractResponse struct
type AbstractResponse struct {
	Code   string `xml:"Code,omitempty"`
	Detail string `xml:"Detail,omitempty"`
}

// ConcreteResponse struct
type ConcreteResponse struct {
	*AbstractResponse
	XMLName           xml.Name `xml:"http://example.com/ns ConcreteResponse"`
	AdditionalMessage string   `xml:"AdditionalMessage,omitempty"`
}
