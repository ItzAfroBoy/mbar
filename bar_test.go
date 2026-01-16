package mbar_test

import (
	"io"
	"net/http"
	"os"

	"github.com/ItzAfroBoy/mbar"
)

func Example() {
	mb := mbar.NewMBar(mbar.Config{ShowTime: true, ShowSpeed: true, ShowSize: true})
	res, err := http.Get("https://ash-speed.hetzner.com/100MB.bin")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	bar := mb.Add("100MB.bin", int(res.ContentLength))
	file, err := os.Create("100MB.bin")
	if err != nil {
		panic(err)
	}
	io.Copy(io.MultiWriter(file, bar), res.Body)
	mb.Finish("Done")
}
