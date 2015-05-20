package main

import (
	"encoding/xml"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"net/http"
	"time"
)

func feed(domains []string) {
	log.Info(domains)
	var subs []Subscription
	for _, domain := range domains {
		subs = append(subs, Subscribe(Fetch(domain)))
	}
	mainFeed := Merge(subs...)
	updates := mainFeed.Updates()
	select {
	case item := <-updates:
		log.Info(item.Title)
	case <-time.After(time.Second * 10):
		mainFeed.Close()
	}
}

type Feed chan *[]Item

type Fetcher interface {
	Fetch() (items []Item, next time.Time, err error)
}

type pollFetcher struct {
	rssUrl string
}

func (p *pollFetcher) Fetch() ([]Item, time.Time, error) {
	next := time.Now().Add(10 * time.Second)
	response, err := http.Get(p.rssUrl)
	if err != nil {
		return nil, next, err
	}
	d := xml.NewDecoder(response.Body)
	var items []Item
	dErr := d.Decode(&items)
	log.Debug(len(items))
	return items, next, dErr
}

func Fetch(domain string) Fetcher {
	return &pollFetcher{"https://" + domain}
}

type Item struct {
	Title, Channel, GUID string
}

type Subscription interface {
	Updates() <-chan Item
	Close() []error
}

type concreteSub struct {
	feed       chan Item
	fetchers   []Fetcher
	closing    chan chan []error
	errors     []error
	maxPending int
}

type fetchResult struct {
	fetched []Item
	next    time.Time
	err     error
}

func (cs *concreteSub) Updates() <-chan Item {
	if cs.feed == nil {
		cs.feed = make(chan Item)
	}
	return cs.feed
}

func (cs *concreteSub) Close() []error {
	if cs.closing == nil {
		cs.closing = make(chan chan []error)
	}
	closeChan := make(chan []error)
	cs.closing <- closeChan
	close(cs.feed)
	item, ok := <-cs.feed
	if ok == true {
		return []error{fmt.Errorf("Couldn't close feed, read %v", item)}
	}
	return <-closeChan
}

func (cs *concreteSub) loop() {
	for _, f := range cs.fetchers {
		var next time.Time
		var pending []Item
		var err error
		var fetchDone chan fetchResult
		guids := make(map[string]bool)
		go func() {
			for {
				var feed chan Item
				var first Item
				if len(pending) > 0 {
					first = pending[0]
					feed = cs.feed
				}
				var delay time.Duration
				if now := time.Now(); next.After(now) {
					delay = next.Sub(time.Now())
				}
				var startFetch <-chan time.Time
				if fetchDone == nil && len(pending) < cs.maxPending {
					startFetch = time.After(delay)
				}
				select {
				case closingChan := <-cs.closing:
					closingChan <- cs.errors
					return
				case <-startFetch:
					var items []Item
					fetchDone = make(chan fetchResult, 1)
					go func() {
						items, next, err = f.Fetch()
						fetchDone <- fetchResult{items, next, err}
					}()

				case result := <-fetchDone:
					fetchDone = nil
					if result.err != nil {
						log.Errorln(err)
						cs.errors = append(cs.errors, err)
						next = time.Now().Add(10 * time.Second)
						break
					}
					next = result.next
					for _, item := range result.fetched {
						if !guids[item.GUID] {
							pending = append(pending, item)
							guids[item.GUID] = true
						}
					}
				case feed <- first:
					pending = pending[1:]
				}
			}
		}()
	}
}

func Subscribe(fetcher Fetcher) Subscription {
	c := &concreteSub{
		feed:       make(chan Item),
		fetchers:   []Fetcher{fetcher},
		maxPending: 1000,
	}
	go c.loop()
	return c
}

func Merge(subscriptions ...Subscription) Subscription {
	mergedSub := &concreteSub{
		feed:     make(chan Item),
		fetchers: []Fetcher{},
	}
	for _, sub := range subscriptions {
		convert := sub.(*concreteSub)
		mergedSub.fetchers = append(mergedSub.fetchers, convert.fetchers...)
		sub.Close()
	}
	mergedSub.loop()
	return mergedSub
}
