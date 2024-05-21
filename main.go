package main

import (
	"io"
	"net/http"
	"time"
)

func main() {
	brasilApiChan := make(chan string)
	viaCepChan := make(chan string)

	cep := "01001000"
	go brasilApi(cep, brasilApiChan)
	go viaCep(cep, viaCepChan)

	select {
	case res := <-brasilApiChan:
		println("Brasil API:", res)
	case res := <-viaCepChan:
		println("Via Cep:", res)
	case <-time.After(1 * time.Second):
		println("Timeout")
	}
}

func brasilApi(cep string, ch chan string) error {
	url := "https://brasilapi.com.br/api/cep/v1/" + cep
	if err := request(url, ch); err != nil {
		return err
	}
	return nil
}

func viaCep(cep string, ch chan string) error {
	url := "http://viacep.com.br/ws/" + cep + "/json"
	if err := request(url, ch); err != nil {
		return err
	}
	return nil
}

func request(url string, ch chan string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	ch <- string(body)
	return nil
}
