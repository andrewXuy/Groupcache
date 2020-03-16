package main

import (
	"cache"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var db = map[string]string{
	"Google":   "666",
	"FaceBook": "255",
	"Amazon":   "532",
}

func creatGroup() *cache.Group {
	return cache.NewGroup("rank", 1024, cache.FetcherFunc(
		func(key string) ([]byte, error) {
			log.Printf("Search key in Database", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exists", key)
		}))
}

// Start server port default as 8001-8003
func startServer(addr string, addrs []string, g *cache.Group) {
	peers := cache.NewPool(addr)
	peers.Set(addrs...)
	g.RegisterPeers(peers)
	log.Println("group runing at: ", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))

}

func startAPIserver(apiAddr string, g *cache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			data, err := g.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(data.ByteSlice())
		}))
	log.Println("fronted server is runing at ", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}
func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "cache server port")
	flag.BoolVar(&api, "api", false, "start api server?")
	flag.Parse()
	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}
	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}
	g := creatGroup()
	if api {
		go startAPIserver(apiAddr, g)
	}
	startServer(addrMap[port], []string(addrs), g)

}
