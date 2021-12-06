package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Vishnukvsvk/NearestMirrorSite/mirrors"
)

type response struct {
	FastestUrl string        `json:"fastest_url"`
	Latency    time.Duration `json:"latency"`
}

func main() {
	http.HandleFunc("/fastest-mirror", func(w http.ResponseWriter, r *http.Request) {
		response := findFastest(mirrors.MirrorList)
		respJSON, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.Write(respJSON)
	})
	port := ":8080"
	server := &http.Server{
		Addr:           port,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Printf("Starting server on port %s", port)
	log.Fatal(server.ListenAndServe())
}

func findFastest(urls []string) response {
	urlChan := make(chan string)
	latencyChan := make(chan time.Duration)

	for _, url := range urls {
		mirrorUrl := url
		go func() {
			start := time.Now()
			_, err := http.Get(mirrorUrl + "/README")
			latency := time.Since(start) / time.Millisecond
			if err == nil {
				urlChan <- mirrorUrl
				latencyChan <- latency
			}
		}()
	}

	return response{<-urlChan, <-latencyChan}
}

//The smart logic here is, whenever a goroutine receives a response, it writes data into two channels with the URL and latency information respectively. Upon receiving data, the two channels make the response struct and return from the findFastest function. When that function is returned, all goroutines spawned from that are stopped from whatever they are doing. So, we will have the shortest URL in urlChan and the smallest latency in latencyChan.
