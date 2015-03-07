package goics

import (
	
)

// Package goics implements an ical encoder and decoder.
// First release will include decode and encoder of Event types
// Will try to add more features as needed.

// ICalDecoder is the really important part of the decoder lib
// 
type ICalConsumer interface {
	ConsumeICal(d *Calendar, err error) error
}

type Calendar struct {
	Data map[string]*IcsNode
	Events []*Event
}
 
type Event struct {
	Data map[string]*IcsNode
	Alarms []*map[string]*IcsNode
}

// http://www.kanzaki.com/docs/ical/vevent.html
