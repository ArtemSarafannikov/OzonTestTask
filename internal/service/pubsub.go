package service

import (
	"github.com/ArtemSarafannikov/OzonTestTask/internal/models"
	"sync"
)

type PubSub struct {
	mu          sync.RWMutex
	subscribers map[string][]chan *models.Comment
	publishCh   chan publishEvent
}

type publishEvent struct {
	postID  string
	comment *models.Comment
}

func NewPubSub() *PubSub {
	pb := &PubSub{
		subscribers: make(map[string][]chan *models.Comment),
		publishCh:   make(chan publishEvent, 50),
	}
	go pb.startPublisher()
	return pb
}

func (p *PubSub) Subscribe(postID string) (<-chan *models.Comment, func()) {
	ch := make(chan *models.Comment, 1)

	p.mu.Lock()
	p.subscribers[postID] = append(p.subscribers[postID], ch)
	p.mu.Unlock()

	unsubscribe := func() {
		p.mu.Lock()
		defer p.mu.Unlock()

		subs := p.subscribers[postID]
		for i, c := range subs {
			if c == ch {
				p.subscribers[postID] = append(subs[:i], subs[i+1:]...)
				break
			}
		}
		close(ch)
	}
	return ch, unsubscribe
}

func (p *PubSub) Publish(postID string, comment *models.Comment) {
	p.publishCh <- publishEvent{postID, comment}
}

func (p *PubSub) startPublisher() {
	for event := range p.publishCh {
		p.mu.RLock()
		subs := p.subscribers[event.postID]
		p.mu.RUnlock()

		for _, ch := range subs {
			select {
			case ch <- event.comment:
			default:
			}
		}
	}
}
