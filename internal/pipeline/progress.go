package pipeline

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// Live-redraw cadence for the spinner+timer line. 100ms is fast enough that
// the elapsed counter feels live without flooding the terminal.
const spinnerInterval = 100 * time.Millisecond

// labelColumn is the column at which the trailing duration is right-padded
// so the durations line up vertically across phases.
const labelColumn = 52

var spinnerFrames = []rune{'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'}

// phaseTimer prints a labelled, in-place updating timer for a phase of
// generation, modeled loosely on `docker pull` per-layer progress lines.
//
// When stdout is not a TTY it falls back to two plain lines (start/end), so
// piped/CI output stays grep-friendly. When enabled is false, all methods are
// no-ops — used in --debug mode where streaming `make` output would fight
// the spinner for the cursor.
type phaseTimer struct {
	label   string
	w       io.Writer
	enabled bool
	tty     bool
	start   time.Time
	stop    chan struct{}
	done    chan struct{}
}

func newPhaseTimer(label string, enabled bool) *phaseTimer {
	pt := &phaseTimer{
		label:   label,
		w:       os.Stdout,
		enabled: enabled,
	}
	if enabled {
		if stat, err := os.Stdout.Stat(); err == nil {
			pt.tty = (stat.Mode() & os.ModeCharDevice) != 0
		}
	}
	return pt
}

// Start prints the initial line and, on a TTY, launches a redraw goroutine
// that overwrites the line every spinnerInterval with the new elapsed time.
func (pt *phaseTimer) Start() {
	if !pt.enabled {
		return
	}
	pt.start = time.Now()
	if !pt.tty {
		fmt.Fprintf(pt.w, "%s…\n", pt.label)
		return
	}
	pt.stop = make(chan struct{})
	pt.done = make(chan struct{})
	fmt.Fprintf(pt.w, "%s %s… %s", string(spinnerFrames[0]), pt.label, formatDur(0))
	go pt.tick()
}

func (pt *phaseTimer) tick() {
	defer close(pt.done)
	t := time.NewTicker(spinnerInterval)
	defer t.Stop()
	frame := 0
	for {
		select {
		case <-pt.stop:
			return
		case <-t.C:
			frame = (frame + 1) % len(spinnerFrames)
			// \r returns to start of line; \x1b[K clears to end of line so
			// shorter redraws don't leave trailing characters from a prior
			// (longer) frame.
			fmt.Fprintf(pt.w, "\r\x1b[K%s %s… %s",
				string(spinnerFrames[frame]), pt.label, formatDur(time.Since(pt.start)))
		}
	}
}

// Stop finalizes the line. ok controls the leading glyph (✓ on success,
// ✗ on failure).
func (pt *phaseTimer) Stop(ok bool) {
	if !pt.enabled {
		return
	}
	elapsed := time.Since(pt.start)
	if !pt.tty {
		verb := "done"
		if !ok {
			verb = "failed"
		}
		fmt.Fprintf(pt.w, "%s %s in %s\n", pt.label, verb, formatDur(elapsed))
		return
	}
	close(pt.stop)
	<-pt.done
	glyph := "✓"
	if !ok {
		glyph = "✗"
	}
	pad := max(labelColumn-len(pt.label), 1)
	fmt.Fprintf(pt.w, "\r\x1b[K%s %s%s%s\n",
		glyph, pt.label, strings.Repeat(" ", pad), formatDur(elapsed))
}

// formatDur produces a compact human-readable elapsed time:
//
//	d < 1m  → "12.4s"
//	d ≥ 1m  → "1m23s"
func formatDur(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	m := int(d / time.Minute)
	s := int((d - time.Duration(m)*time.Minute) / time.Second)
	return fmt.Sprintf("%dm%02ds", m, s)
}
