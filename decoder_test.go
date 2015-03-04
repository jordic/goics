package goics_test

import (
	"os"
	"strings"
	"testing"

	goics "github.com/jordic/goics"
)

func TestTesting(t *testing.T) {
	if 1 != 1 {
		t.Error("Error setting up testing")
	}
}

var source = "asdf\nasdf\nasdf"

func TestEndOfFile(t *testing.T) {
	a := goics.NewDecoder(strings.NewReader(source))
	err := a.Decode()
	if err != goics.VCalendarNotFound {
		t.Errorf("Decode filed, decode raised %s", err)
	}
	if a.Lines() != 3 {
		t.Errorf("Decode should advance to %s", a.Lines())
	}
	if a.Calendar != nil {
		t.Errorf("No calendar in file")
	}
}

var test2 = `BEGIN:VCALENDAR
PRODID;X-RICAL-TZSOURCE=TZINFO:-//test//EN
CALSCALE:GREGORIAN
VERSION:2.0
END:VCALENDAR

`

func TestInsideCalendar(t *testing.T) {
	a := goics.NewDecoder(strings.NewReader(test2))
	err := a.Decode()
	if err != nil {
		t.Errorf("Failed %s", err)
	}
	if a.Calendar.Calscale != "GREGORIAN" {
		t.Error("No extra keys for calendar decoded")
	}
	if a.Calendar.Version != "2.0" {
		t.Error("No extra keys for calendar decoded")
	}
}

var test3 = `BEGIN:VCALENDAR
PRODID;X-RICAL-TZSOURCE=TZINFO:-//test//EN
CALSCALE:GREGORIAN
VERSION:2.`

func TestDetectIncompleteCalendar(t *testing.T) {
	a := goics.NewDecoder(strings.NewReader(test3))
	err := a.Decode()
	if err != goics.VParseEndCalendar {
		t.Error("Test failed")
	}

}

var testlonglines = `BEGIN:VCALENDAR
PRODID;X-RICAL-TZSOURCE=TZINFO:-//test//EN
CALSCALE:GREGORIANGREGORIANGREGORIANGREGORIANGREGORIANGREGORIANGREGORIAN
 GREGORIANGREGORIAN
VERSION:2.0
END:VCALENDAR
`

func TestParseLongLines(t *testing.T) {
	a := goics.NewDecoder(strings.NewReader(testlonglines))
	_ = a.Decode()
	str := a.Calendar.Calscale
	if len(str) != 81 {
		t.Errorf("Multiline test failed %s", len(a.Calendar.Params["CALSCALE"]))
	}
	if strings.Contains("str", " ") {
		t.Error("Not handling correct begining of line")
	}

}

var testlonglinestab = `BEGIN:VCALENDAR
PRODID;X-RICAL-TZSOURCE=TZINFO:-//test//EN
CALSCALE:GREGORIANGREGORIANGREGORIANGREGORIANGREGORIANGREGORIANGREGORIAN
	GREGORIANGREGORIAN
VERSION:2.0
END:VCALENDAR
`

func TestParseLongLinesTab(t *testing.T) {
	a := goics.NewDecoder(strings.NewReader(testlonglinestab))
	_ = a.Decode()
	str := a.Calendar.Calscale
	if len(str) != 81 {
		t.Errorf("Multiline tab field test failed %s", len(str))
	}
	if strings.Contains("str", "\t") {
		t.Error("Not handling correct begining of line")
	}

}

var testlonglinestab3 = `BEGIN:VCALENDAR
PRODID;X-RICAL-TZSOURCE=TZINFO:-//test//EN
CALSCALE:GREGORIANGREGORIANGREGORIANGREGORIANGREGORIANGREGORIANGREGORIAN
	GREGORIANGREGORIANGREGORIANGREGORIANGREGORIANGREGORIANGREGORIANGREGORIANGG
 GRESESERSERSER
VERSION:2.0
END:VCALENDAR
`

func TestParseLongLinesMultilinethree(t *testing.T) {
	a := goics.NewDecoder(strings.NewReader(testlonglinestab3))
	_ = a.Decode()
	str := a.Calendar.Calscale
	if len(str) != 151 {
		t.Errorf("Multiline (3lines) tab field test failed %s", len(str))
	}
	if strings.Contains("str", "\t") {
		t.Error("Not handling correct begining of line")
	}

}

var testevent = `BEGIN:VCALENDAR
BEGIN:VEVENT
DTEND;VALUE=DATE:20140506
DTSTART;VALUE=DATE:20140501
UID:-kpd6p8pqal11-n66f1wk1tw76@xxxx.com
DESCRIPTION:CHECKIN:  01/05/2014\nCHECKOUT: 06/05/2014\nNIGHTS:   5\nPHON
 E:    \nEMAIL:    (no se ha facilitado ningún correo electrónico)\nPRO
 PERTY: Apartamento xxx 6-8 pax en Centro\n
SUMMARY:Luigi Carta (FYSPZN)
LOCATION:Apartamento xxx 6-8 pax en Centro
END:VEVENT
END:VCALENDAR
`

func TestVEvent(t *testing.T) {
	a := goics.NewDecoder(strings.NewReader(testevent))
	err := a.Decode()
	if err != nil {
		t.Errorf("Error decoding %s", err)
	}
	if len(a.Calendar.Events) != 1 {
		t.Error("Not decoding events", len(a.Calendar.Events))
	}

}

var valarmCt = `BEGIN:VCALENDAR
BEGIN:VEVENT
STATUS:CONFIRMED
CREATED:20131205T115046Z
UID:1ar5d7dlf0ddpcih9jum017tr4@google.com
DTEND;VALUE=DATE:20140111
TRANSP:OPAQUE
SUMMARY:PASTILLA Cu cs
DTSTART;VALUE=DATE:20140110
DTSTAMP:20131205T115046Z
LAST-MODIFIED:20131205T115046Z
SEQUENCE:0
DESCRIPTION:
BEGIN:VALARM
X-WR-ALARMUID:E283310A-82B3-47CF-A598-FD36634B21F3
UID:E283310A-82B3-47CF-A598-FD36634B21F3
TRIGGER:-PT15H
X-APPLE-DEFAULT-ALARM:TRUE
ATTACH;VALUE=URI:Basso
ACTION:AUDIO
END:VALARM
END:VEVENT
END:VCALENDAR`

func TestNotParsingValarm(t *testing.T) {
	a := goics.NewDecoder(strings.NewReader(valarmCt))
	err := a.Decode()
	if err != nil {
		t.Errorf("Error decoding %s", err)
	}
}

var tParseError = `BEGIN:VCALENDAR
BEGIN:VEVENT
DTEND;VALUE=DATE:20140230a
END:VEVENT
END:VCALENDAR`

func TestParserError(t *testing.T) {

	a := goics.NewDecoder(strings.NewReader(tParseError))
	err := a.Decode()
	if err == nil {
		t.Error("Should return a parsing error")
	}
	if a.Lines() != 3 {
		t.Error("Wrong line error reported")
	}

}

func TestReadingRealFile(t *testing.T) {

	file, err := os.Open("fixtures/test.ics")
	if err != nil {
		t.Error("Can't read file")
	}
	defer file.Close()

	cal := goics.NewDecoder(file)
	err = cal.Decode()
	if err != nil {
		t.Error("Cant decode a complete file")
	}
	
	if len(cal.Calendar.Events) != 28 {
		t.Errorf("Wrong number of events detected %s", len(cal.Calendar.Events))
	}

	if cal.Calendar.Events[0].Summary != "Clarisse De  (AZfTDA)" {
		t.Errorf("Wrong summary")
	}

}
