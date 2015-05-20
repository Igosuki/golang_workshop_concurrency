package main

import (
	"encoding/xml"
	// "fmt"
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
	Close() error
}

type concreteSub struct {
	feed     chan Item
	fetchers []Fetcher
	closed   bool
	err      error
}

func (cs *concreteSub) Updates() <-chan Item {
	if cs.feed == nil {
		cs.feed = make(chan Item)
	}
	return cs.feed
}

func (cs *concreteSub) Close() error {
	cs.closed = true
	return cs.err
}

func (cs *concreteSub) loop() {
	for _, f := range cs.fetchers {
		go func() {
			for {
				items, next, err := f.Fetch()
				if err != nil {
					cs.err = err
					time.Sleep(10 * time.Second)
					continue
				}
				for _, item := range items {
					cs.feed <- item
				}
				if now := time.Now(); next.After(now) {
					time.Sleep(next.Sub(now))
				}
			}
		}()
	}
}

func Subscribe(fetcher Fetcher) Subscription {
	c := &concreteSub{
		feed:     make(chan Item),
		fetchers: []Fetcher{fetcher},
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
