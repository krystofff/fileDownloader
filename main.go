package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	pb "gopkg.in/cheggaaa/pb.v1"
)

const RefreshRate = time.Millisecond * 100

// almacena nro de bytes escritos y barra de progreso
type WriteCounter struct {
	n   int
	bar *pb.ProgressBar
}

// configuracion de barra de progreso
func NewWriteCounter(total int) *WriteCounter {
	b := pb.New(total)
	b.SetRefreshRate(RefreshRate)
	b.ShowTimeLeft = true
	b.ShowSpeed = true
	b.SetUnits(pb.U_BYTES)

	return &WriteCounter{
		bar: b,
	}
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	wc.n += len(p)
	wc.bar.Set(wc.n)
	return wc.n, nil
}

func (wc *WriteCounter) Start() {
	wc.bar.Start()
}

func (wc *WriteCounter) Finish() {
	wc.bar.Finish()
}

func main() {

	// la url se debe pasar como argumento para ejecución desde terminal
	if len(os.Args) < 2 {
		fmt.Println("Debes pasar una url como argumento para iniciar la descarga.")
		os.Exit(1)
	}

	// lee URL y descarga datos
	fileURL := os.Args[1]
	err := DownloadFile(fileURL, path.Base(fileURL))
	if err != nil {
		panic(err)
	}
}

// Descarga url en archivo local
// io.TeeReader reporta el progreso de la descarga
func DownloadFile(url string, filename string) error {
	// crea archivo temporal para no sobreescribir alguno con el mismo nombre
	out, err := os.Create(filename + ".tmp")
	if err != nil {
		return err
	}
	defer out.Close()

	// Obtiene datos de URL
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fsize, _ := strconv.Atoi(resp.Header.Get("Content-Length"))

	// Crea WriteCounter y le asigna el tamaño especificado por el Content-Length
	counter := NewWriteCounter(fsize)
	counter.Start()

	// io.TeeReader lee body de la respuesta y escribe en un nuevo WriteCounter
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		return err
	}

	counter.Finish()

	err = os.Rename(filename+".tmp", filename)
	if err != nil {
		return err
	}

	return nil
}
