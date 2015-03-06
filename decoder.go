package goics

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"
)

const (
	keySep    = ":"
	vBegin    = "BEGIN"
	vCalendar = "VCALENDAR"
	vEnd      = "END"
	vEvent    = "VEVENT"

	maxLineRead = 65
)

// Errors
var (
	VCalendarNotFound = errors.New("vCalendar not found")
	VParseEndCalendar = errors.New("wrong format END:VCALENDAR not Found")
)

type decoder struct {
	scanner      *bufio.Scanner
	err          error
	Calendar     *Calendar
	currentEvent *Event
	nextFn       stateFn
	prevFn       stateFn
	current      string
	buffered     string
	line         int
}

type stateFn func(*decoder)

func NewDecoder(r io.Reader) *decoder {
	d := &decoder{
		scanner: bufio.NewScanner(r),
		//calendar: &Calendar{},
		nextFn:   decodeInit,
		line:     0,
		buffered: "",
	}
	return d
}

func (d *decoder) Decode() (err error) {
	d.next()
	if d.Calendar == nil {
		d.err = VCalendarNotFound
		return d.err
	}
	// If theres no error but, nextFn is not reset
	// last element not closed
	if d.nextFn != nil && d.err == nil {
		d.err = VParseEndCalendar
		return d.err
	}
	return d.err
}

// Decode events should decode events found in buff
// into a slice of passed struct &Event
// Idealy struct should add mapping capabilities in tag strings
// like:
//
// type Event struct {
//		ID string `ics:"UID"`
//		Start time.Time `ics:"DTSTART"`
// }
// v must be a slice of []*Event
//
func (d *decoder) DecodeEvents(dest interface{}) error {
	err := d.Decode()
	if err != nil {
		return err
	}
	
	dt := reflect.ValueOf(dest).Elem()
	direct := reflect.Indirect( reflect.ValueOf(dest) )
	fmt.Println( dt.Kind() )
	if dt.Kind() == reflect.Slice {
		tipo := dt.Type().Elem()
		eventitem := reflect.New(tipo)
		item := eventitem.Elem()
		for _, ev := range d.Calendar.Events {
			for i:=0; i<tipo.NumField(); i++ {
				typeField := tipo.Field(i)
				if strings.ToLower(typeField.Name) == "dtstart" {
					item.Field(i).Set( reflect.ValueOf(ev.Start) )
				}
				fmt.Printf("Field Name %s, field type %s\n", typeField.Name, typeField)
			}
			direct.Set(reflect.Append(dt, item))
			fmt.Println(dt)
		}
	}
	/*to := reflect.Indirect(reflect.ValueOf(v))
	
		tipo := reflect.ValueOf(v).Type()
		fmt.Println(tipo)
		fmt.Println("Is a slice", to.Type().Elem() )
		if len(d.Calendar.Events) == 0 {
			return nil			
		}
		//
			for i:=0; i<tipo.NumField(); i++ {
				typeField := tipo.Field(i)
				fmt.Println("Field Name %s, field type %s", typeField.Name, typeField)
			}
		//}
	}*/
	return nil
}



// Lines processed. If Decoder reports an error.
// Error
func (d *decoder) Lines() int {
	return d.line
}

// Current Line content
func (d *decoder) CurrentLine() string {
	return d.current
}

// Advances a new line in the decoder
// And calls the next stateFunc
// checks if next line is continuation line
func (d *decoder) next() {
	// If there's not buffered line
	if d.buffered == "" {
		res := d.scanner.Scan()
		if true != res {
			d.err = d.scanner.Err()
			return
		}
		d.line++
		d.current = d.scanner.Text()
	} else {
		d.current = d.buffered
		d.buffered = ""
	}

	if len(d.current) > 65 {
		is_continuation := true
		for is_continuation == true {
			res := d.scanner.Scan()
			d.line++
			if true != res {
				d.err = d.scanner.Err()
				return
			}
			line := d.scanner.Text()
			if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
				d.current = d.current + line[1:]
			} else {
				// If is not a continuation line, buffer it, for the
				// next call.
				d.buffered = line
				is_continuation = false
			}
		}
	}

	if d.nextFn != nil {
		d.nextFn(d)
	}
}

func decodeInit(d *decoder) {
	node := DecodeLine(d.current)
	if node.Key == vBegin && node.Val == vCalendar {
		d.Calendar = &Calendar{
			Params: make(map[string]string),
		}
		d.prevFn = decodeInit
		d.nextFn = decodeInsideCal
		d.next()
		return
	}
	d.next()
}

func decodeInsideCal(d *decoder) {
	node := DecodeLine(d.current)
	switch {
	case node.Key == vBegin && node.Val == vEvent:
		d.currentEvent = &Event{
			Params: make(map[string]string),
		}
		d.nextFn = decodeInsideEvent
		d.prevFn = decodeInsideCal
	case node.Key == vEnd && node.Val == vCalendar:
		d.nextFn = nil
	case node.Key == "VERSION":
		d.Calendar.Version = node.Val
	case node.Key == "PRODID":
		d.Calendar.Prodid = node.Val
	case node.Key == "CALSCALE":
		d.Calendar.Calscale = node.Val
	default:
		d.Calendar.Params[node.Key] = node.Val
	}
	d.next()
}

func decodeInsideEvent(d *decoder) {

	node := DecodeLine(d.current)
	if node.Key == vEnd && node.Val == vEvent {
		// Come back to parent node
		d.nextFn = d.prevFn
		d.Calendar.Events = append(d.Calendar.Events, d.currentEvent)
		d.next()
		return
	}
	//@todo handle Valarm

	var err error
	var t time.Time
	switch {
	case node.Key == "UID":
		d.currentEvent.Uid = node.Val
	case node.Key == "DESCRIPTION":
		d.currentEvent.Description = node.Val
	case node.Key == "SUMMARY":
		d.currentEvent.Summary = node.Val
	case node.Key == "LOCATION":
		d.currentEvent.Location = node.Val
	case node.Key == "STATUS":
		d.currentEvent.Status = node.Val
	case node.Key == "TRANSP":
		d.currentEvent.Transp = node.Val
	// Date based
	case node.Key == "DTSTART":
		t, err = dateDecode(node)
		d.currentEvent.Start = t
	case node.Key == "DTEND":
		t, err = dateDecode(node)
		d.currentEvent.End = t
	case node.Key == "LAST-MODIFIED":
		t, err = dateDecode(node)
		d.currentEvent.LastModified = t
	case node.Key == "DTSTAMP":
		t, err = dateDecode(node)
		d.currentEvent.Dtstamp = t
	case node.Key == "CREATED":
		t, err = dateDecode(node)
		d.currentEvent.Created = t
	default:
		d.currentEvent.Params[node.Key] = node.Val
	}
	if err != nil {
		//@todo improve error notification, adding node info.. and line number
		d.err = err
	} else {
		d.next()
	}
}
