package goics_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	goics "github.com/jordic/goics"
)

func TestDateFieldFormat(t *testing.T) {
	result := "DTEND;VALUE=DATE:20140406\r\n"
	ti := time.Date(2014, time.April, 06, 0, 0, 0, 0, time.UTC)
	str := goics.WriteDateField("DTEND", ti)
	if result != str {
		t.Error("Expected", result, "Result", str)
	}
}

func TestDateTimeFieldFormat(t *testing.T) {
	result := "X-MYDATETIME;VALUE=DATE-TIME:20120901T130000\r\n"
	ti := time.Date(2012, time.September, 01, 13, 0, 0, 0, time.UTC)
	str := goics.WriteDateTimeField("X-MYDATETIME", ti)
	if result != str {
		t.Error("Expected", result, "Result", str)
	}
}

func TestDateTimeFormat(t *testing.T) {
	result := "DTSTART:19980119T070000Z\r\n"
	ti := time.Date(1998, time.January, 19, 07, 0, 0, 0, time.UTC)
	str := goics.WriteDateTime("DTSTART", ti)
	if result != str {
		t.Error("Expected", result, "Result", str)
	}
}

var shortLine string = `asdf defined is a test\n\r`

func TestLineWriter(t *testing.T) {
	
	w := &bytes.Buffer{}

	result := &bytes.Buffer{}
	fmt.Fprintf(result, shortLine)

	encoder := goics.NewWriter(w)
	encoder.WriteLine(shortLine)

	res := bytes.Compare(w.Bytes(), result.Bytes())

	if res != 0 {
		t.Errorf("%s!=%s", w, result)
	}

}

var longLine string = `As returned by NewWriter, a Writer writes records terminated by thisisat test that is expanded in multi lines` + goics.CRLF

func TestLineWriterLongLine(t *testing.T) {
	
	w := &bytes.Buffer{}

	result := &bytes.Buffer{}
	fmt.Fprintf(result, "As returned by NewWriter, a Writer writes records terminated by thisisat ")
	fmt.Fprintf(result, goics.CRLFSP)
	fmt.Fprintf(result, "test that is expanded in multi lines")
	fmt.Fprintf(result, goics.CRLF)

	encoder := goics.NewWriter(w)
	encoder.WriteLine(longLine)

	res := bytes.Compare(w.Bytes(), result.Bytes())

	if res != 0 {
		t.Errorf("%s!=%s %s", w, result, res)
	}
}

func Test2ongLineWriter(t *testing.T) {
	goics.LineSize = 10
	
	w := &bytes.Buffer{}

	result := &bytes.Buffer{}
	fmt.Fprintf(result, "12345678")
	fmt.Fprintf(result, goics.CRLF)
	fmt.Fprintf(result, " 2345678")
	fmt.Fprintf(result, goics.CRLF)
	fmt.Fprintf(result, " 2345678")

	var str string = `1234567823456782345678`
	encoder := goics.NewWriter(w)
	encoder.WriteLine(str)

	res := bytes.Compare(w.Bytes(), result.Bytes())

	if res != 0 {
		t.Errorf("%s!=%s %s", w, result, res)
	}

}


func TestContentEvent(t *testing.T) {
	goics.LineSize = 75
	c := &goics.Event{
		Start: time.Date(2014, time.April,   04,  0, 0, 0, 0, time.UTC),
		End: time.Date(2014, time.April,   06,  0, 0, 0, 0, time.UTC),
		Uid: "-kpd6p8pqal11-74iythu9giqs@xxx.com",
		Location: "Apartamento xxxx pax en Centro",
		Description: "test",
	}
	w := &bytes.Buffer{}
	writer := goics.NewWriter(w)
	result := &bytes.Buffer{}
	
	fmt.Fprintf(result, "BEGIN:VEVENT" + goics.CRLF)
	fmt.Fprintf(result, "DTEND;VALUE=DATE:20140406" + goics.CRLF)
	fmt.Fprintf(result, "DTSTART;VALUE=DATE:20140404" + goics.CRLF)
	fmt.Fprintf(result, "UID:-kpd6p8pqal11-74iythu9giqs@xxx.com" + goics.CRLF)
	fmt.Fprintf(result, "DESCRIPTION:test" + goics.CRLF)
	fmt.Fprintf(result, "LOCATION:Apartamento xxxx pax en Centro" + goics.CRLF)
	fmt.Fprintf(result, "END:VEVENT" + goics.CRLF)
	
	goics.WriteEvent(c, writer)
	
	res := bytes.Compare(w.Bytes(), result.Bytes())

	if res != 0 {
		t.Errorf("%s!=%s %s", w, result, res)
	}

	
}