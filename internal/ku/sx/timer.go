package sx

import (
	"cmp"
	"fmt"
	"io"
	"slices"
	"time"
)

// Section identifies compiler/build section for storing
// and creating timings.
type Section uint32

const (
	SectionParse Section = 1 + iota
)

// Timer carries start time and section name in compact way.
type Timer struct {
	// stores timer start timestamp in unix nanoseconds
	// or total number of nanoseconds for collected timer
	val uint64

	name Section
}

func (t Timer) String() string {
	var name string
	switch t.name {
	case SectionParse:
		name = "parse"
	default:
		panic(fmt.Sprintf("invalid section (=%d)", t.name))
	}

	nano := t.val
	ms := nano / 1_000_000
	micro := (nano - ms*1_000_000) / 1000

	return fmt.Sprintf("%s: %d.%03d ms", name, ms, micro)
}

func (p *Pool) Start(name Section) Timer {
	if p.timers == nil {
		return Timer{}
	}

	return Timer{
		name: name,
		val:  uint64(time.Now().UnixNano()),
	}
}

func (p *Pool) Stop(t Timer) {
	if p.timers == nil || t.name == 0 {
		return
	}

	end := uint64(time.Now().UnixNano())
	if end <= t.val {
		return
	}

	p.timers[t.name] += time.Duration(end-t.val) * time.Nanosecond
}

func (p *Pool) Timers() []Timer {
	if len(p.timers) == 0 {
		return nil
	}

	timers := make([]Timer, 0, len(p.timers))
	for name, dur := range p.timers {
		timers = append(timers, Timer{
			val:  uint64(dur),
			name: name,
		})
	}
	slices.SortFunc(timers, func(a, b Timer) int {
		return cmp.Compare(a.name, b.name)
	})
	return timers
}

func (p *Pool) DumpTimers(out io.Writer) {
	timers := p.Timers()
	for _, t := range timers {
		fmt.Fprintf(out, "%s\n", t)
	}
}
