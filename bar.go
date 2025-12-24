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
	NumBars int
	bars    []*Bar
	update  chan bool
}

type Bar struct {
	Title       string
	Size        int
	id          int
	currentSize int
	content     string
	update      chan bool
}

func calculateProgress(current, max, width int) (pc, blks int) {
	pc = int(math.Floor((float64(current) / float64(max)) * float64(100)))
	blks = int(math.Floor(float64(width) * (float64(pc) / float64(100))))
	return
}

func NewMBar() *MBar {
	return &MBar{update: make(chan bool)}
}

func (MBar) newBar(title string, size int, update chan bool) *Bar {
	return &Bar{Title: title, Size: size, id: 0, currentSize: 0, update: update}
}

func (b *Bar) genBar(n int) error {
	b.currentSize = b.currentSize + n
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	barSize := w / 2
	percent, barProg := calculateProgress(b.currentSize, b.Size, barSize)
	b.content = fmt.Sprintf("\033[K%s [%s%s] [%s/%s | %d%%]", b.Title, strings.Repeat("#", barProg), strings.Repeat("-", barSize-barProg), humanize.Bytes(uint64(b.currentSize)), humanize.Bytes(uint64(b.Size)), percent)
	return err
}

func (b *Bar) Write(p []byte) (n int, err error) {
	n = len(p)
	err = b.genBar(n)
	b.update <- true
	return
}

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

func (m *MBar) Add(title string, size int) *Bar {
	bar := m.newBar(title, size, m.update)
	m.bars = append(m.bars, bar)
	m.NumBars++
	return bar
}

func (m *MBar) Finish(msg string) {
	var buf string
	for _, b := range m.bars {
		buf += b.content
		buf += "\n"
	}
	buf += msg + "\n"
	fmt.Fprint(os.Stdout, buf)
}
