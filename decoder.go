package goics

import (
	"bufio"
	"errors"
	//"fmt"
	"io"
	"strings"
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
	if d.nextFn != nil {
		d.err = VParseEndCalendar
		return d.err
	}
	return d.err
}

func (d *decoder) Lines() int {
	return d.line
}

func (d *decoder) CurrentLine() string {
	return d.current
}

// Advances a new token in the decoder
// And calls the next stateFunc
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
	if strings.Contains(d.current, keySep) {
		key, val := getKeyVal(d.current)
		if key == vBegin {
			if val == vCalendar {
				d.Calendar = &Calendar{}
				d.nextFn = decodeInsideCal
				d.next()
				return
			}
		}
	}
	d.nextFn = decodeInit
	d.next()
}

func decodeInsideCal(d *decoder) {
	if strings.Contains(d.current, keySep) {
		key, val := getKeyVal(d.current)
		if key == vBegin && val == vEvent {
			d.currentEvent = &Event{}
			d.nextFn = decodeInsideEvent
			d.next()
			return
		}
		if key == vEnd && val == vCalendar {
			d.nextFn = nil
			d.next()
			return
		}

		if key != "" && val != "" {
			if d.Calendar.Extra == nil {
				d.Calendar.Extra = make(map[string]string)
			}
			d.Calendar.Extra[key] = val
			d.next()
			return
		}

	}
	d.nextFn = decodeInsideCal
	d.next()
}

func decodeInsideEvent(d *decoder) {

	node := DecodeLine(d.current)
	if node.Key == vEnd && node.Val == vEvent {
		d.nextFn = decodeInsideCal
		d.Calendar.Events = append(d.Calendar.Events, d.currentEvent)
		d.next()
		return
	}

	switch {
	case node.Key == "UID":
		d.currentEvent.Uid = node.Val
	case node.Key == "DESCRIPTION":
		d.currentEvent.Description = node.Val
	case node.Key == "SUMMARY":
		d.currentEvent.Summary = node.Val
	case node.Key == "LOCATION":
		d.currentEvent.Location = node.Val

	}

	d.next()
}
