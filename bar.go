// A multi-progressbar solution.
// Supports I/O operations only at this point in development
package mbar

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"golang.org/x/term"
)

type MBar struct {
	// Number of bars
	NumBars int
	// Config for the displaying of info to the bar
	Config Config
	// Bar instances
	bars []*Bar
	// Update channel used to update the bars
	update chan bool
}

type Bar struct {
	// Title of the bar
	Title string
	// The full size of the object being operated on
	Size int
	// Config for the displaying of info to the bar
	Config Config
	// ID of the bar
	id int
	// Current amount written so far
	currentSize int
	// Total writes
	writes int
	// Start time for I/O operations
	startTime time.Time
	// The output for the bar
	content string
	// The channel to send write updates to
	update chan bool
}

type Config struct {
	ShowTime  bool
	ShowSpeed bool
	ShowSize  bool
}

// Calcualtes the progress for display on the bar.
// Returns the percentage completed.
func calculateProgress(current, max int) (pc int) {
	pc = int(math.Floor((float64(current) / float64(max)) * float64(100)))
	return
}

// Calcualtes the progress for the bar.
// Returns the total bar blocks and how many to fill.
func calculateBarProgress(maxWidth, percent int, els ...int) (barSize, barProg int) {
	var strWidth int
	for _, v := range els {
		strWidth += v
	}
	barSize = maxWidth - strWidth
	barProg = int(math.Floor(float64(barSize) * (float64(percent) / float64(100))))
	return
}

// Calcualtes the progress for display on the bar.
// Returns a string containing the I/O speed and
// how long is left to complete the operation.
func calculateStats(bar *Bar) string {
	if bar.Config.ShowSpeed || bar.Config.ShowTime {
		diff := time.Since(bar.startTime)
		speed := float64(bar.currentSize) / diff.Seconds()
		eta := float64(bar.Size-bar.currentSize) / speed
		if bar.Config.ShowSpeed && bar.Config.ShowTime {
			return fmt.Sprintf(" [%s/s | %s]", humanize.Bytes(uint64(speed)), (time.Duration(eta) * time.Second).Round(time.Second).String())
		} else if bar.Config.ShowSpeed {
			return fmt.Sprintf(" [%s/s]", humanize.Bytes(uint64(speed)))
		}
		return fmt.Sprintf(" [%s]", (time.Duration(eta) * time.Second).Round(time.Second).String())
	}
	return ""
}

// Returns a new nulti-bar instance
func NewMBar(config Config) *MBar {
	return &MBar{Config: config, update: make(chan bool)}
}

// Returns a new bar to attach to a multi-bar instance
func (MBar) newBar(title string, size int, config Config, update chan bool) *Bar {
	return &Bar{Title: title, Size: size, Config: config, id: 0, currentSize: 0, writes: 0, update: update}
}

// Generates the bar for output to the terminal
func (b *Bar) genBar(n int) error {
	var suffix string
	if b.writes == 1 {
		b.startTime = time.Now()
	}
	b.currentSize = b.currentSize + n
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	percent := calculateProgress(b.currentSize, b.Size)
	stats := calculateStats(b)
	if b.Config.ShowSize {
		suffix = fmt.Sprintf(" [%s/%s | %d%%]", humanize.Bytes(uint64(b.currentSize)), humanize.Bytes(uint64(b.Size)), percent)
	} else {
		suffix = ""
	}
	barSize, barProg := calculateBarProgress(w, percent, len(stats), len(suffix), len(b.Title), 3)
	b.content = fmt.Sprintf("\033[K%s [%s%s]%s%s", b.Title, strings.Repeat("#", barProg), strings.Repeat("-", barSize-barProg), stats, suffix)

	return err
}

// io.writer implementation
func (b *Bar) Write(p []byte) (n int, err error) {
	b.writes++
	n = len(p)
	err = b.genBar(n)
	b.update <- true
	return
}

// Start will update the bar whenever the Write
// function is called. This must be called before any
// I/O operations have begun
//
// BUG(Brett): Fonts with ligitures enabled may have unexpected behaviour
func (m *MBar) Start() {
	go func() {
		for range m.update {
			var buf strings.Builder
			if m.NumBars > 0 {
				for _, b := range m.bars {
					buf.WriteString(b.content)
					buf.WriteString("\n")
				}
				buf.WriteString(fmt.Sprintf("\033[%dA\r", m.NumBars))
				fmt.Fprint(os.Stdout, buf.String())
			}
		}
	}()
}

// Add will add a new bar to the multi-bar instance.
func (m *MBar) Add(title string, size int) *Bar {
	bar := m.newBar(title, size, m.Config, m.update)
	m.bars = append(m.bars, bar)
	m.NumBars++
	return bar
}

// Finish is the final function when all bars are done.
// It must be called to reset the cursor's
// position for output purposes.
func (m *MBar) Finish(msg string) {
	var buf strings.Builder
	for _, b := range m.bars {
		buf.WriteString(b.content)
		buf.WriteString("\n")
	}
	buf.WriteString(msg + "\n")
	fmt.Fprint(os.Stdout, buf.String())
}
