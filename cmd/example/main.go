package main

import (
	"log"
	"time"

	"github.com/fbv/expiration"
)

func main() {

	out := make(chan string)

	q := expiration.NewQueue(out)
	defer q.Close()

	go func() {
		now := time.Now()
		q.Add("1", now.Add(8*time.Second))
		q.Add("2", now.Add(3*time.Second))
		q.Add("3", now.Add(2*time.Second))
		q.Add("4", now.Add(5*time.Second))
		q.Add("0", now.Add(time.Duration(-2)*time.Second))
		q.Add("5", now.Add(7*time.Second))
		q.Add("1", now.Add(1*time.Second))
		log.Printf("queue len: %d", q.Len())
		q.Remove("2")
		log.Printf("queue len: %d", q.Len())
	}()

	for {
		select {
		case id := <-out:
			log.Printf("expired: %s", id)
			log.Printf("queue len: %d", q.Len())
		case <-time.After(10 * time.Second):
			log.Printf("done")
			return
		}
	}
}
