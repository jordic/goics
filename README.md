
### just another go icalendar file encoder/decoder ics format

Alpha status

Features implemented:

- Parsing and writing of vevent, not completly.. 
- No recursive events, 
- And no alarms inside events


Follows:
http://tools.ietf.org/html/rfc5545


TODO
Reflect:
https://gist.github.com/drewolson/4771479

Tests from:
https://github.com/libical/libical

Decoder` that you create by passing it an `io.Reader`, similar to your setup
Joe Shaw [9:37 PM]
and then you `err := dec.Decode(&foo)`, where `foo` for `json` is `interface{}` but 

also curious why `WriteCalendar` and `WriteEvent` aren't methods of `*Writer`

Joe Shaw [8:10 PM]
seems like you'd want `WriteLine` and a lot of the helper functions to be unexported

i think returning an unexported struct from an exported function is considered a no-no, which you do with `func NewDecoder(r io.Reader) *decoder`