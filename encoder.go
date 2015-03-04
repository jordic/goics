package goics

import (
	"io"
	"strings"
	"time"
)

const (
	CRLF   = "\r\n"
	CRLFSP = "\r\n "
)

type Writer struct {
	w    io.Writer
	opts map[string]func(string, time.Time) string
}


var defaultOptions = map[string]func(string, time.Time) string{
	"Start":        WriteDateField,
	"End":          WriteDateField,
	"LastModified": WriteDateTimeField,
	"Dtstamp":      WriteDateTimeField,
	"Created":      WriteDateTimeField,
}

var LineSize int = 75

// Creates a New writer with default options for date fields
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w:    w,
		opts: defaultOptions,
	}
}

// Creates a new Writer with custom options for Date Fields
func NewWriterOpts(w io.Writer, opts map[string]func(string, time.Time) string) *Writer {
	return &Writer{
		w:    w,
		opts: opts,
	}
}

// Writes a calendar to writer.
func WriteCalendar(c *Calendar, w *Writer) {
	w.WriteLine("BEGIN:VCALENDAR" + CRLF)
	if c.Prodid != "" {
		w.WriteLine(WriteStringField("Prodid", c.Prodid))
	} else {
		w.WriteLine(WriteStringField("Prodid", "-//tmpo.io/src/goics"))
	}
	if c.Calscale != "" {
		w.WriteLine(WriteStringField("Calscale", c.Calscale))
	} else {
		w.WriteLine(WriteStringField("Calscale", "GREGORIAN"))
	}
	if c.Version != "" {
		w.WriteLine(WriteStringField("Version", c.Version))
	} else {
		w.WriteLine(WriteStringField("Version", "2.0"))
	}
	if c.Uid != "" {
		w.WriteLine(WriteStringField("Uid", c.Uid))
	}
	for key, val := range c.Params {
		w.WriteLine(WriteStringField(key, val))
	}
	
	for _, ev := range c.Events {
		WriteEvent(ev, w)
	}
	
	w.WriteLine("END:VCALENDAR" + CRLF)
}



// Writes an event struct to a writer
func WriteEvent(ev *Event, w *Writer) {

	w.WriteLine("BEGIN:VEVENT" + CRLF)
	// Date fields
	t := time.Time{}
	
	if ev.End != t {
		w.WriteLine(w.opts["End"]("DTEND", ev.End))
	}
	if ev.Start != t {
		w.WriteLine(w.opts["Start"]("DTSTART", ev.Start))
	}
	if ev.LastModified != t {
		w.WriteLine(w.opts["LastModified"]("LAST-MODIFIED", ev.LastModified))
	}
	if ev.Dtstamp != t {
		w.WriteLine(w.opts["Dtstamp"]("DTSTAMP", ev.Dtstamp))
	}
	if ev.Created != t {
		w.WriteLine(w.opts["Created"]("CREATED", ev.Created))
	}
	// String fields
	if ev.Uid != "" {
		w.WriteLine(WriteStringField("Uid", ev.Uid))
	}
	if ev.Summary != "" {
		w.WriteLine(WriteStringField("Summary", ev.Summary))
	}
	if ev.Description != "" {
		w.WriteLine(WriteStringField("Description", ev.Description))
	}
	if ev.Location != "" {
		w.WriteLine(WriteStringField("Location", ev.Location))
	}
	if ev.Status != "" {
		w.WriteLine(WriteStringField("Status", ev.Status))
	}
	if ev.Transp != "" {
		w.WriteLine(WriteStringField("Transp", ev.Transp))
	}
	w.WriteLine("END:VEVENT" + CRLF)
}

func quoteString(s string) string {
	s = strings.Replace(s, "\\;", ";", -1)
	s = strings.Replace(s, "\\,", ",", -1)
	s = strings.Replace(s, "\\n", "\n", -1)
	s = strings.Replace(s, "\\\\", "\\", -1)
	return s
}

// Write a line in ics format max length = 75
// continuation lines start with a space.
func (w *Writer) WriteLine(s string) {
	if len(s) <= LineSize {
		io.WriteString(w.w, s)
		return
	}
	length := len(s)
	current := 0
	// LineSize -2 is CRLF
	shortLine := LineSize - 2
	// First line write from 0 to totalline - 2 ( must include CRLFS)
	io.WriteString(w.w, s[current:current+(shortLine)]+CRLFSP)
	current = shortLine
	// Rest of lines, we must include ^space at begining for marquing
	// continuation lines
	for (current + shortLine) <= length {
		io.WriteString(w.w, s[current:current+(shortLine-1)]+CRLFSP)
		current += shortLine - 1
	}
	// Also we need to write the reminder
	io.WriteString(w.w, s[current:length])
}

// Writes a Date in format: "DTEND;VALUE=DATE:20140406"
func WriteDateField(key string, val time.Time) string {
	return key + ";VALUE=DATE:" + val.Format("20060102") + CRLF
}

// "X-MYDATETIME;VALUE=DATE-TIME:20120901T130000"
func WriteDateTimeField(key string, val time.Time) string {
	return key + ";VALUE=DATE-TIME:" + val.Format("20060102T150405") + CRLF
}

// "DTSTART:19980119T070000Z"
func WriteDateTime(key string, val time.Time) string {
	return key + ":" + val.Format("20060102T150405Z") + CRLF
}

// Write a key field UID:asdfasdfаs@dfasdf.com
func WriteStringField(key string, val string) string {
	return strings.ToUpper(key) + ":" + quoteString(val) + CRLF
}
