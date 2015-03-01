package goics_test

import (
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

var test2 =`BEGIN:VCALENDAR
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
	if a.Calendar.Extra["CALSCALE"] != "GREGORIAN" {
		t.Error("No extra keys for calendar decoded")
	}
	if a.Calendar.Extra["VERSION"] != "2.0" {
		t.Error("No extra keys for calendar decoded")
	}
}

var test3 =`BEGIN:VCALENDAR
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


var testlonglines =`BEGIN:VCALENDAR
PRODID;X-RICAL-TZSOURCE=TZINFO:-//test//EN
CALSCALE:GREGORIANGREGORIANGREGORIANGREGORIANGREGORIANGREGORIANGREGORIAN
 GREGORIANGREGORIAN
VERSION:2.0
END:VCALENDAR
`
func TestParseLongLines(t *testing.T) {
	a := goics.NewDecoder(strings.NewReader(testlonglines))
	_ = a.Decode()
	str := a.Calendar.Extra["CALSCALE"]
	if len(str) != 81 {
		t.Errorf("Multiline test failed %s", len(a.Calendar.Extra["CALSCALE"]) )
	}
	if strings.Contains("str", " ") {
		t.Error("Not handling correct begining of line")
	}
	
}

var testlonglinestab =`BEGIN:VCALENDAR
PRODID;X-RICAL-TZSOURCE=TZINFO:-//test//EN
CALSCALE:GREGORIANGREGORIANGREGORIANGREGORIANGREGORIANGREGORIANGREGORIAN
	GREGORIANGREGORIAN
VERSION:2.0
END:VCALENDAR
`
func TestParseLongLinesTab(t *testing.T) {
	a := goics.NewDecoder(strings.NewReader(testlonglinestab))
	_ = a.Decode()
	str := a.Calendar.Extra["CALSCALE"]
	if len(str) != 81 {
		t.Errorf("Multiline tab field test failed %s", len(str) )
	}
	if strings.Contains("str", "\t") {
		t.Error("Not handling correct begining of line")
	}
	
}

var testlonglinestab3 =`BEGIN:VCALENDAR
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
	str := a.Calendar.Extra["CALSCALE"]
	if len(str) != 151 {
		t.Errorf("Multiline (3lines) tab field test failed %s", len(str) )
	}
	if strings.Contains("str", "\t") {
		t.Error("Not handling correct begining of line")
	}
	
}

var testevent =`BEGIN:VCALENDAR
BEGIN:VEVENT
DTEND;VALUE=DATE:20140506
DTSTART;VALUE=DATE:20140501
UID:-kpd6p8pqal11-n66f1wk1tw76@airbnb.com
DESCRIPTION:CHECKIN:  01/05/2014\nCHECKOUT: 06/05/2014\nNIGHTS:   5\nPHON
 E:    \nEMAIL:    (no se ha facilitado ningún correo electrónico)\nPRO
 PERTY: Apartamento Muji 6-8 pax en Centro\n
SUMMARY:Luigi Carta (FYSPZN)
LOCATION:Apartamento Muji 6-8 pax en Centro
END:VEVENT
END:VCALENDAR
`
func TestVEvent(t *testing.T) {
	a := goics.NewDecoder(strings.NewReader(testevent))
	err := a.Decode()
	if err != nil {
		t.Errorf("Error decoding %s", err)
	}
	if len(a.Calendar.Events)!= 1 {
		t.Error("Not decoding events", len(a.Calendar.Events))
	}

}