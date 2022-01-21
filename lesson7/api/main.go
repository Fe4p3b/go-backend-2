package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/mediocregopher/radix/v3"
	"gocloud.dev/pubsub"
	_ "gocloud.dev/pubsub/kafkapubsub"
)

var (
	connFunc = func(network, addr string) (radix.Conn, error) {
		return radix.Dial(network, addr,
			radix.DialTimeout(time.Minute),
		)
	}
	t *pubsub.Topic
	s *radix.Pool
)

func main() {
	port := 8080

	r := chi.NewRouter()
	r.Post("/rate", PostRateHandler)
	r.Get("/total", GetTotalHandler)

	wg := new(sync.WaitGroup)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(port int) {
			defer wg.Done()
			addr := fmt.Sprintf(":%d", port)

			log.Printf("starting server at port - %d", port)
			if err := http.ListenAndServe(addr, r); err != http.ErrServerClosed {
				log.Println(fmt.Errorf("error on listen and serve: %v", err))
				return
			}

		}(port)
		port++
	}
	wg.Wait()
}

func GetTotalHandler(w http.ResponseWriter, r *http.Request) {
	var rates []string
	err := storage().Do(radix.Cmd(&rates, "LRANGE", "result", "0", "10"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(rates) == 0 {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	var sum int
	for _, rate := range rates {
		v, err := strconv.Atoi(rate)
		if err != nil {
			continue
		}
		sum += v
	}
	result := float64(sum) / float64(len(rates))
	log.Printf("result - %.2f", result)
	_, _ = w.Write([]byte(fmt.Sprintf("%.2f", result)))
}

func PostRateHandler(w http.ResponseWriter, r *http.Request) {
	rate := r.FormValue("rate")
	log.Printf("rate - %s", rate)
	if _, err := strconv.Atoi(rate); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := topic().Send(context.Background(), &pubsub.Message{
		Body: []byte(rate),
	})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(rate))
}

func topic() *pubsub.Topic {
	if t != nil {
		return t
	}
	var err error
	t, err = pubsub.OpenTopic(context.Background(), "kafka://rates")
	if err != nil {
		panic(err)
	}
	return t
}

func storage() *radix.Pool {
	if s != nil {
		return s
	}
	var err error
	s, err = radix.NewPool("tcp", "redis:6379", 1, radix.PoolConnFunc(connFunc))
	if err != nil {
		panic(err)
	}
	return s
}
