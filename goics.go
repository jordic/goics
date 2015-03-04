package goics

import (
	"time"
)

// Package goics implements an ical encoder and decoder.
// First release will include decode and encoder of Event types
// Will try to add more features as needed.

type Calendar struct {
	Uid    string
	Events []*Event
	Calscale string
	Version string
	Prodid string
	Params  map[string]string
}

// http://www.kanzaki.com/docs/ical/vevent.html
type Event struct {
	Start        time.Time
	End          time.Time
	LastModified time.Time
	Dtstamp      time.Time
	Created      time.Time
	Uid          string
	Summary      string
	Description  string
	Location     string
	Status       string
	Transp       string
	Params       map[string]string
	Alarms []*Alarm
}

// http://www.kanzaki.com/docs/ical/valarm.html
// @todo, not implemented
type Alarm struct {

}
