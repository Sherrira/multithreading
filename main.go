package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func fetch(url string, ch chan<- string, apiName chan<- string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	ch <- string(body)
	apiName <- url
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Path[len("/cep/"):]

	ch := make(chan string)
	apiName := make(chan string)

	go fetch("http://viacep.com.br/ws/"+cep+"/json/", ch, apiName)
	go fetch("https://brasilapi.com.br/api/cep/v1/"+cep, ch, apiName)

	select {
	case res := <-ch:
		fmt.Fprintf(w, "Fastest response from: %s\nData: %s\n", <-apiName, res)
	case <-time.After(1 * time.Second):
		fmt.Fprintln(w, "Timeout")
	}
}

func main() {
	http.HandleFunc("/cep/", handleRequest)
	http.ListenAndServe(":8080", nil)
}
