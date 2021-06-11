package main

import (
	"GoCache"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	GoCache.NewGroup("scores", 2<<10, GoCache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[DB] search key: " + key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := "localhost:9999"
	peers := GoCache.NewHTTPPool(addr)
	log.Println("Running at ", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
