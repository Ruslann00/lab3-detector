package main

import (
	"log"
	"net/http"

	_ "net/http/pprof"

	"lab3-detector/internal/processor"
)

func main() {
	go func() {
		log.Println("Pprof server started on http://localhost:6060/debug/pprof/")
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	log.Println("Image Metadata Processor started...")

	processor.RunWorkerPool(5)
}
