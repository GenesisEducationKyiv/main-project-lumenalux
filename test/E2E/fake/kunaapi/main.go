package main

import (
	"fmt"
	"log"
	"net/http"
)

var (
	Response = `[["btcuah",1227057,394.7022,1231,5.547443,521,0.04,1225381,0.86715,1242872,1212000]]`
	Port     = ":8082"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, Response)
	})

	fmt.Printf("Serving on localhost%s", Port)
	err := http.ListenAndServe(Port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
