package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	lock := flag.Bool("pingponglock", false, "set this to true to demonstrate a dead lock situation using channels")
	flag.Parse()

	table := make(chan *Ball)
	go player("Ping !", table)
	go player("Pong !", table)

	if !*lock {
		table <- &Ball{}
	}
	<-time.After(time.Second * 2)
	<-table

}

type Ball struct{ hits int }

func player(name string, table chan *Ball) {
	for {
		ball := <-table
		ball.hits++
		println(name, ball.hits)
		time.Sleep(time.Millisecond * 100)
		table <- ball
	}
}

type Feed chan *[]Item

type Fetcher interface {
	Fetch() (items []Item, next time.Time, err error)
}

type pollFetcher struct {
}

func (p *pollFetcher) Fetch() ([]Item, time.Time, error) {
	return []Item{}, time.Now(), nil
}

func Fetch(domain string) Fetcher {
	return &pollFetcher{}
}

type Item struct {
	Title, Channel, GUID string
}

type Subscription interface {
	Updates() chan Item
	Close() error
}

type concreteSub struct {
	feed chan Item
}

func (cs *concreteSub) Updates() chan Item {
	if cs.feed == nil {
		cs.feed = make(chan Item)
	}
	return cs.feed
}

func (cs *concreteSub) Close() error {
	close(cs.feed)
	item, ok := <-cs.feed
	if ok == true {
		return fmt.Errorf("Couldn't close feed, read %v", item)
	}
	return nil
}

func Subscribe(fetcher Fetcher) Subscription {
	return &concreteSub{}
}

func Merge(subscriptions ...Subscription) Subscription {
	mergedSub := &concreteSub{
		feed: make(chan Item),
	}
	for _, sub := range subscriptions {
		sub.Close()
	}
	return mergedSub
}
