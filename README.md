<div align="center">
<pre>
         ____     ______  ____       
 /'\_/`\/\  _`\  /\  _  \/\  _`\     
/\      \ \ \L\ \\ \ \L\ \ \ \L\ \   
\ \ \__\ \ \  _ <'\ \  __ \ \ ,  /   
 \ \ \_/\ \ \ \L\ \\ \ \/\ \ \ \\ \  
  \ \_\\ \_\ \____/ \ \_\ \_\ \_\ \_\
   \/_/ \/_/\/___/   \/_/\/_/\/_/\/ /
<br>
A shitty multi-line bar for Go
<br>
<img alt="GitHub License" src="https://img.shields.io/github/license/ItzAfroBoy/mbar"> <img alt="GitHub tag (with filter)" src="https://img.shields.io/github/v/tag/ItzAfroBoy/mbar?label=version"> <a href="https://www.codefactor.io/repository/github/itzafroboy/mbar"><img src="https://www.codefactor.io/repository/github/itzafroboy/mbar/badge" alt="CodeFactor" /></a> <img alt="GitHub code size in bytes" src="https://img.shields.io/github/languages/code-size/ItzAfroBoy/mbar">
</pre>
</div>

## Installation

```shell
go get github.com/ItzAfroBoy/mbar@latest
```

## [Example](https://github.com/ItzAfroBoy/mbar/blob/main/bar_test.go)

```go
package main

import (
 "io"
 "net/http"
 "os"

 "github.com/ItzAfroBoy/mbar"
)

func main() {
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
```
