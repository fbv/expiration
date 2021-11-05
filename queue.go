package expiration

import (
	"log"
	"time"
)

type Queue struct {
	done     chan bool
	toAdd    chan *Entry
	toRemove chan string
	items    []*Entry
}

type Entry struct {
	ID     string
	Expire time.Time
}

func NewQueue(expired chan string) *Queue {
	q := &Queue{
		done:     make(chan bool),
		toAdd:    make(chan *Entry),
		toRemove: make(chan string),
		items:    make([]*Entry, 0),
	}
	go q.loop(expired)
	return q
}

func (q *Queue) loop(expired chan string) {
	log.Println("start expiration loop")
	var next *Entry
	var waitTime time.Duration
	var timer *time.Timer
	for {
		next = nil
		for _, d := range q.items {
			if next == nil || next.Expire.After(d.Expire) {
				next = d
			}
		}
		waitTime = time.Duration(time.Hour)
		if next != nil {
			waitTime = time.Until(next.Expire)
			if waitTime <= 0 {
				q.remove(next.ID)
				expired <- next.ID
				continue
			}
		}
		timer = time.NewTimer(waitTime)
		log.Printf("wait %v", waitTime)
		select {
		case <-timer.C:
			if next != nil {
				q.remove(next.ID)
				expired <- next.ID
			}
		case <-q.done:
			log.Println("end expiration loop")
			timer.Stop()
			return
		case e := <-q.toAdd:
			timer.Stop()
			q.add(e)
		case id := <-q.toRemove:
			timer.Stop()
			q.remove(id)
		}
	}
}

func (q *Queue) Add(id string, expire time.Time) {
	q.toAdd <- &Entry{
		ID:     id,
		Expire: expire,
	}
}

func (q *Queue) add(e *Entry) {
	found := false
	for _, d := range q.items {
		if d.ID == e.ID {
			d.Expire = e.Expire
			found = true
			break
		}
	}
	if !found {
		q.items = append(q.items, e)
	}
}

func (q *Queue) Remove(id string) {
	q.toRemove <- id
}

func (q *Queue) remove(id string) {
	for i, d := range q.items {
		if d.ID == id {
			j := len(q.items) - 1
			q.items[i] = q.items[j]
			q.items[j] = nil
			q.items = q.items[:j]
			break
		}
	}
}

func (q *Queue) Len() int {
	return len(q.items)
}

func (q *Queue) Close() {
	q.done <- true
}
