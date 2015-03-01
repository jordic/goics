package goics

import (
	"strings"
	"testing"
	"time"
	
)


var event_t4 string = `BEGIN:VCALENDAR
PRODID;X-RICAL-TZSOURCE=TZINFO:-//Airbnb Inc//Hosting Calendar 0.8.8//EN
CALSCALE:GREGORIAN
VERSION:2.0
BEGIN:VEVENT
DTEND;VALUE=DATE:20140406
DTSTART;VALUE=DATE:20140404
UID:-kpd6p8pqal11-74iythu9giqs@xxx.com
DESCRIPTION:CHECKIN:  04/04/2014\nCHECKOUT: 06/04/2014\nNIGHTS:   2\nPHON
 E:    \nEMAIL:    (no se ha facilitado ningún correo e
 lectrónico)\nPROPERTY: Apartamento Test 6-8 pax en Centro\n
SUMMARY:Clarisse (AZ8TDA)
LOCATION:Apartamento Test 6-8 pax en Centro
END:VEVENT
END:VCALENDAR`

func TestVEventCharProperties(t *testing.T) {
	a := NewDecoder(strings.NewReader(event_t4))
	err := a.Decode()
	if err != nil {
		t.Errorf("Error decoding properties %s", err)
	}
	if len(a.Calendar.Events) != 1 {
		t.Error("Wrong number of events decoded", len(a.Calendar.Events))
	}
	ev := a.Calendar.Events[0]
	if ev.Location != "Apartamento Test 6-8 pax en Centro" {
		t.Error("Wrong location param")
	}
	if strings.HasPrefix(ev.Description, "CHECKIN") != true {
		t.Error("Wrong description param")
	}
	if strings.Contains(ev.Description, "6-8 pax en Centro") != true {
		t.Error("Wrong description param", ev.Description)
	}
	if ev.Summary != "Clarisse (AZ8TDA)" {
		t.Error("Wrong summary param", ev.Summary)
	}
	if ev.Uid != "-kpd6p8pqal11-74iythu9giqs@xxx.com" {
		t.Error("Wrong uid", ev.Uid)
	}
}


var dates_tested = []string{
	"DTEND;VALUE=DATE:20140406",
	"DTSTART;TZID=Europe/Paris:20140116T120000",
}

func get_timezone(zone string) *time.Location {
	t, _ := time.LoadLocation(zone)
	return t
}


var times_expected =[]time.Time{
	time.Date(2014, time.April,   06,  0, 0, 0, 0, time.UTC),
	time.Date(2014, time.January, 16, 12, 0, 0, 0, get_timezone("Europe/Paris")),
}


func TestDateDecode(t *testing.T) {

	for i, d :=	range dates_tested {
		
		node := DecodeLine(d)
		res, err := dateDecode(node)
		if err != nil {
			t.Errorf("Error decoding time %s", err)
		}
		if res.String() != times_expected[i].String() {
			t.Errorf("Error parsing time %s expected %s", res, times_expected[i])
		}
	
	}

}