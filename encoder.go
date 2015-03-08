package goics

import (
	"io"
	"time"
	"strings"
)

const (
	CRLF   = "\r\n"
	CRLFSP = "\r\n "
)



func NewComponent() *Component {
	return &Component{
		Elements:   make([]Componenter, 0),
		Properties: make(map[string]string),
	}
}

type Component struct {
	Tipo       string
	Elements   []Componenter
	Properties map[string]string
}

func (c *Component) Write(w *ICalEncode) {
	w.WriteLine("BEGIN:" + c.Tipo + CRLF )
	// Iterate over component properites
	for key, val := range c.Properties {
		w.WriteLine( WriteStringField(key, val) )
	
	}
	for _, xc := range c.Elements {
		xc.Write( w )	
	}
	
	w.WriteLine("END:" + c.Tipo + CRLF)
}

func (c *Component) SetType(t string) {
	c.Tipo = t
}

func (c *Component) AddComponent(cc Componenter) {
	c.Elements = append(c.Elements, cc)
}

func (c *Component) AddProperty(key string, val string) {
	c.Properties[key] = val
}


// ICalEncode is the real writer, that wraps every line, 
// in 75 chars length... Also gets the component from the emmiter
// and starts the iteration.

type ICalEncode struct {
	w io.Writer
}

func NewICalEncode(w io.Writer) *ICalEncode {
	return &ICalEncode{
		w: w,
	}
}

func (enc *ICalEncode) Encode(c ICalEmiter) {
	component := c.EmitICal()
	component.Write(enc)
}

var LineSize int = 75

// Write a line in ics format max length = 75
// continuation lines start with a space.
func(enc *ICalEncode)  WriteLine(s string) {
	if len(s) <= LineSize {
		io.WriteString(enc.w, s)
		return
	}
	length := len(s)
	current := 0
	// LineSize -2 is CRLF
	shortLine := LineSize - 2
	// First line write from 0 to totalline - 2 ( must include CRLFS)
	io.WriteString(enc.w, s[current:current+(shortLine)]+CRLFSP)
	current = shortLine
	// Rest of lines, we must include ^space at begining for marquing
	// continuation lines
	for (current + shortLine) <= length {
		io.WriteString(enc.w, s[current:current+(shortLine-1)]+CRLFSP)
		current += shortLine - 1
	}
	// Also we need to write the reminder
	io.WriteString(enc.w, s[current:length])
}




// Writes a Date in format: "DTEND;VALUE=DATE:20140406"
func FormatDateField(key string, val time.Time) (string, string) {
	return key + ";VALUE=DATE", val.Format("20060102")
}

// "X-MYDATETIME;VALUE=DATE-TIME:20120901T130000"
func FormatDateTimeField(key string, val time.Time) (string, string) {
	return key + ";VALUE=DATE-TIME", val.Format("20060102T150405")
}

// "DTSTART:19980119T070000Z"
func FormatDateTime(key string, val time.Time) (string, string) {
	return key, val.Format("20060102T150405Z") 
}

// Write a key field UID:asdfasdfаs@dfasdf.com
func WriteStringField(key string, val string) string {
	return strings.ToUpper(key) + ":" + quoteString(val) + CRLF
}

func quoteString(s string) string {
	s = strings.Replace(s, "\\;", ";", -1)
	s = strings.Replace(s, "\\,", ",", -1)
	s = strings.Replace(s, "\\n", "\n", -1)
	s = strings.Replace(s, "\\\\", "\\", -1)
	return s
}