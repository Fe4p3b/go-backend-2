package main

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"

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
	sub *pubsub.Subscription
	s   *radix.Pool
)

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

func subscription() (*pubsub.Subscription, error) {
	if sub != nil {
		return sub, nil
	}
	var err error
	sub, err = pubsub.OpenSubscription(context.Background(), "kafka://process?topic=rates")
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func main() {
	wg := new(sync.WaitGroup)
	n := 10

	for i := 0; i < n; i++ {
		wg.Add(1)
		go Job(wg, i)
	}
	wg.Wait()
}

func Job(wg *sync.WaitGroup, job int) {
	defer wg.Done()
	log.Printf("job %d started", job)
	for {
		s, err := subscription()
		if err != nil {
			log.Printf("%d - error %v", job, err)
			time.Sleep(time.Second)
			continue
		}
		msg, err := s.Receive(context.Background())
		if err != nil {
			log.Printf("%d - error %v", job, err)
			time.Sleep(time.Second)
			continue
		}
		log.Printf("%d - message received - %v", job, msg.Body)

		err = storage().Do(radix.Cmd(nil, "LPUSH", "result", string(msg.Body)))
		if err != nil {
			log.Printf("%d - error %v", job, err)
		}
		if rand.Float64() < .05 {
			_ = storage().Do(radix.Cmd(nil, "LTRIM", "result", "0", "9"))
		}
		msg.Ack()
	}
}
