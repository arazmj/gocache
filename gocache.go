package main

import (
	"GoCache/LRUCache"
	"GoCache/Stats"
	"flag"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"strconv"
)

var cache LRUCache.LRUCache

func main() {
	capacity := flag.Int("capacity", 100,
		"how big the cache will be, the old values will be evicted")
	port := flag.Int("port", 8080, "the server port number.")
	flag.Parse()

	stats := Stats.NewStats()
	cache = LRUCache.NewCache(*capacity, stats)
	router := mux.NewRouter()
	router.HandleFunc("/cache/{key}", GetHandler).Methods(http.MethodGet)
	router.HandleFunc("/cache/{key}/{value}", PutHandler).Methods(http.MethodPost)
	router.Handle("/metrics", promhttp.Handler())
	log.Printf("GoCache started listening on %d port\n", *port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), router))
}

func PutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	value := vars["value"]
	cache.Put(key, value)
	w.WriteHeader(http.StatusOK)
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	if value, ok := cache.Get(key); ok {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(value))
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
