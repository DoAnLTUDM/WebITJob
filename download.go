package main

import (
    "io"
    "log"
    "net/http"
    "os"
    "reflect"
)

var bytesToKiloBytes = 1024.0

// PassThru code originally from
// http://stackoverflow.com/a/22422650/613575
type PassThru struct {
    io.Reader
    curr  int64
    total float64
}

func (pt *PassThru) Read(p []byte) (int, error) {
    n, err := pt.Reader.Read(p)
    pt.curr += int64(n)

    // last read will have EOF err
    if err == nil || (err == io.EOF && n > 0) {
        log.Println("Downloading image")
    }

    return n, err
}

func download(url string, locat string, file_name string) {
    log.Println(locat)
    if _, err := os.Stat(locat); os.IsNotExist(err){
        direc := os.Mkdir(locat, os.ModePerm)
        log.Print(reflect.TypeOf(direc), direc)
    }
    if _, err := os.Stat(locat+file_name); os.IsNotExist(err) {
        resp, _ := http.Get(url)
        defer resp.Body.Close()

        out, _ := os.Create(locat+file_name)
        defer out.Close()

        src := &PassThru{Reader: resp.Body, total: float64(resp.ContentLength)}

        size, err := io.Copy(out, src)
        if err != nil {
            log.Println(err)
            return
        }

        log.Printf("Download success %s (%.1f KB)\n", file_name, float64(size)/bytesToKiloBytes)
    }
}
