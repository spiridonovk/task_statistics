package api

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"test_task/engine"
	"time"
)

var isLetter = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString

//GetStatistics print statistics in keep-alive con
func GetStatistics(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Server does not support Flusher!",
			http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	ch := make(chan struct{}, 2)
	ch <- struct{}{}
	engine.AddNewClient(ch)
	//waiting for data
	for {
		select {
		case <-ch:
			if len(engine.Stat.Statistics) > 0 {

				engine.Stat.Mu.Lock()
				for n, k := range engine.Stat.Statistics {
					time.Sleep(20 * time.Millisecond) // Simulating work...

					_, err := fmt.Fprintf(w, "%s: %d \n", n, k)
					if err != nil {
						log.Println("Fprintf err:", err.Error())
					}
					flusher.Flush()
				}
				_, err := fmt.Fprintf(w, "=======\n")
				if err != nil {
					log.Println("Fprintf err:", err.Error())
				}
				flusher.Flush()
				engine.Stat.Mu.Unlock()
			}
		}
	}
}

//SetCommand check and save new cmd
func SetCommand(w http.ResponseWriter, r *http.Request) {
	x := r.URL.Path
	x = strings.TrimPrefix(x, "/command/")
	if !isLetter(x) {
		w.WriteHeader(http.StatusBadRequest)
		_, err := fmt.Fprintf(w, "cmd not correct: use only letters")
		if err != nil {
			log.Println("Fprintf err:", err.Error())
		}
		return
	}
	engine.AddToStat(x)
	w.WriteHeader(http.StatusOK)

}
