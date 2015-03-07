package goics

import (
	
	"testing"
	"time"
	
)





var dates_tested = []string{
	"DTEND;VALUE=DATE:20140406",
	"DTSTART;TZID=Europe/Paris:20140116T120000",
	"X-MYDATETIME;VALUE=DATE-TIME:20120901T130000",
	//"RDATE:20131210Z",
	"DTSTART:19980119T070000Z",
}

func get_timezone(zone string) *time.Location {
	t, _ := time.LoadLocation(zone)
	return t
}


var times_expected =[]time.Time{
	time.Date(2014, time.April,   06,  0, 0, 0, 0, time.UTC),
	time.Date(2014, time.January, 16, 12, 0, 0, 0, get_timezone("Europe/Paris")),
	time.Date(2012, time.September,01,  13, 0, 0, 0, time.UTC),
	time.Date(1998, time.January, 19,  07, 0, 0, 0, time.UTC),
}


func TestDateDecode(t *testing.T) {

	for i, d :=	range dates_tested {
		
		node := DecodeLine(d)
		res, err := node.DateDecode()
		if err != nil {
			t.Errorf("Error decoding time %s", err)
		}
		if res.Equal(times_expected[i]) == false {
			t.Errorf("Error parsing time %s expected %s", res, times_expected[i])
		}
		if res.String() != times_expected[i].String() {
			t.Errorf("Error parsing time %s expected %s", res, times_expected[i])
		}
	
	}

}