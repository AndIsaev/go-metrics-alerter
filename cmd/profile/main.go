package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

const (
	limit        = 20
	randoMetrics = 100000
	postReq      = "http://127.0.0.1:8080/"
	path         = "cmd/profile/bench.txt"
)

type Req struct {
	url   string
	body  string
	rType string
}

type Resp struct {
	url  string
	code int
	body string
}

func responder(chResp chan *Resp, chErr chan error) {
	for {
		select {
		case resp := <-chResp:
			if resp.code != 201 && resp.code != 307 {
				log.Printf("url %s, wrong code: %d", resp.url, resp.code)
			}
		case err := <-chErr:
			log.Println(err)
		}
	}
}

func requester(wg *sync.WaitGroup, chReq chan *Req, chResp chan *Resp, chErr chan error) {
	for req := range chReq {
		resp, err := request(req)
		if err != nil {
			chErr <- err
		} else {
			chResp <- resp
		}
		wg.Done()
	}
}

func request(req *Req) (*Resp, error) {
	var resp *http.Response
	var err error
	switch req.rType {
	case "get":
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
		resp, err = client.Get(req.url)
		defer resp.Body.Close()

	case "post":
		resp, err = http.Post(req.url, "text/plain; charset=utf-8", strings.NewReader(req.body))
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("%s request error: %v", req.url, err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%s request error: %v", req.url, err)
	}

	dataResp := Resp{
		url:  req.url,
		code: resp.StatusCode,
		body: string(body),
	}

	return &dataResp, nil
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func main() {
	myDir, _ := os.Getwd()
	dir := fmt.Sprintf("%v/%v", myDir, path)
	ok, _ := exists(dir)

	if !ok {
		metrics := generateRandomMetrics(randoMetrics)
		err := writeToFile(dir, metrics)
		if err != nil {
			log.Fatal("Ошибка при записи файла:", err)
		} else {
			log.Println("Файл успешно создан")
		}
	} else {
		log.Println("файл уже был создан")
	}

	urls := make([]string, 0, randoMetrics)
	chReq := make(chan *Req)
	chResp := make(chan *Resp)
	chErr := make(chan error)
	var wg sync.WaitGroup

	file, err := os.OpenFile(dir, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Разбиваем строку по запятой
		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			fmt.Println("Неправильный формат строки:", line)
			continue
		}
		metricType := parts[0]
		name := parts[1]
		value := parts[2]
		result := fmt.Sprintf("update/%v/%v/%v", metricType, name, value)

		urls = append(urls, result)
	}

	file.Close()

	for i := 0; i < limit; i++ {
		go func() {
			requester(&wg, chReq, chResp, chErr)
		}()
	}
	go responder(chResp, chErr)

	for _, url := range urls {
		_, err := request(&Req{url: postReq + url, rType: "post", body: ""})
		if err != nil {
			log.Printf("url %s request error: %v", url, err)
			continue
		}
	}

	wg.Wait()
	log.Println("Нагрузка произведена")
}
