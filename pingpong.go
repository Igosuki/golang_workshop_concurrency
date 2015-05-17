package main

import (
	"flag"
	"time"
)

func pingpong() {
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
