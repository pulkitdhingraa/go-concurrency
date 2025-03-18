package main

import (
	"context"
	"fmt"
	"sync"
)

// producer-consumer problem with sync.cond
// consumers wait until producers signal that something is ready

type Queue struct {
	data       []string
	maxSize    int
	mu         sync.Mutex
	canConsume *sync.Cond
	canProduce *sync.Cond
}

func NewQueue(maxSize int) *Queue {
	q := &Queue{
		data:    make([]string, 0),
		maxSize: maxSize,
	}
	q.canConsume = sync.NewCond(&q.mu) //Can the consumer take an item from the q?
	q.canProduce = sync.NewCond(&q.mu) //Can the producer add an item to the q?
	return q
}

func (q *Queue) Produce(ctx context.Context, item string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	for len(q.data) >= q.maxSize {
		if ctx.Err() != nil {
			q.mu.Unlock()
			return
		}
		q.canProduce.Wait()
	}

	q.data = append(q.data, item)
	fmt.Printf("Produced: %s\n", item)
	q.canConsume.Signal()
}

func (q *Queue) Consume(ctx context.Context, id int) {
	for {
		q.mu.Lock()
		for len(q.data) == 0 {
			if ctx.Err() != nil {
				q.mu.Unlock()
				fmt.Printf("Consumer %d exiting\n", id)
				return
			}
			q.canConsume.Wait()
		}
		item := q.data[0]
		q.data = q.data[1:]
		q.mu.Unlock()
		fmt.Printf("Consumer %d consumed %s\n", id, item)
		q.mu.Lock()
		q.canProduce.Signal()
		q.mu.Unlock()
	}
}

func main() {
	var producerWg sync.WaitGroup
	var consumerWg sync.WaitGroup
	q := NewQueue(5)

	ctx, cancel := context.WithCancel(context.Background())

	// Start Consumers
	for i := 0; i < 3; i++ {
		consumerWg.Add(1)
		go func(id int) {
			defer consumerWg.Done()
			q.Consume(ctx, id)
		}(i)
	}

	// Start Producers
	for i := 0; i < 3; i++ {
		producerWg.Add(1)
		go func(id int) {
			defer producerWg.Done()
			for j := 0; j < 4; j++ {
				if ctx.Err() != nil {
					return
				}
				q.Produce(ctx, fmt.Sprintf("P-%d Item %d", id, j))
			}
		}(i)
	}

	producerWg.Wait()
	cancel()
	consumerWg.Wait()
	fmt.Println("Main func exiting")
}
