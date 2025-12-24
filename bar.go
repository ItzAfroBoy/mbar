// A multi-progressbar solution.
// Supports I/O operations only at this point in
// development
package mbar

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
	"golang.org/x/term"
)

type MBar struct {
	// Number of bars
	NumBars int
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
	// ID of the bar
	id int
	// Current amount written so far
	currentSize int
	// The output for the bar
	content string
	// The channel to send write updates to
	update chan bool
}

// Calcualtes the progress for the bar
// Returns the percentage and how many blocks to fill depending
// on the width of the bar
func calculateProgress(current, max, width int) (pc, blks int) {
	pc = int(math.Floor((float64(current) / float64(max)) * float64(100)))
	blks = int(math.Floor(float64(width) * (float64(pc) / float64(100))))
	return
}

// Returns a new nulti-bar instance
func NewMBar() *MBar {
	return &MBar{update: make(chan bool)}
}

// Returns a new bar to attach to a multi-bar instance
func (MBar) newBar(title string, size int, update chan bool) *Bar {
	return &Bar{Title: title, Size: size, id: 0, currentSize: 0, update: update}
}

// Generates the bar for output to the terminal
func (b *Bar) genBar(n int) error {
	b.currentSize = b.currentSize + n
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	barSize := w / 2
	percent, barProg := calculateProgress(b.currentSize, b.Size, barSize)
	b.content = fmt.Sprintf("\033[K%s [%s%s] [%s/%s | %d%%]", b.Title, strings.Repeat("#", barProg), strings.Repeat("-", barSize-barProg), humanize.Bytes(uint64(b.currentSize)), humanize.Bytes(uint64(b.Size)), percent)
	return err
}

// io.writer implementation
func (b *Bar) Write(p []byte) (n int, err error) {
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
			var buf string
			if m.NumBars > 0 {
				for _, b := range m.bars {
					buf += b.content
					buf += "\n"
				}
				buf += fmt.Sprintf("\033[%dA\r", m.NumBars)
				fmt.Fprint(os.Stdout, buf)
			}
		}
	}()
}

// Add will add a new bar to the multi-bar instance.
func (m *MBar) Add(title string, size int) *Bar {
	bar := m.newBar(title, size, m.update)
	m.bars = append(m.bars, bar)
	m.NumBars++
	return bar
}

// Finish is the final function when all bars are done.
// It must be called to reset the cursor's
// position for output purposes.
func (m *MBar) Finish(msg string) {
	var buf string
	for _, b := range m.bars {
		buf += b.content
		buf += "\n"
	}
	buf += msg + "\n"
	fmt.Fprint(os.Stdout, buf)
}
