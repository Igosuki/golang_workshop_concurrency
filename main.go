package main

import (
	"time"
)

func main() {
	table := make(chan *Ball)
	go player("Ping !", table)
	go player("Pong !", table)
	table <- &Ball{}
	<-time.After(time.Second * 5)
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
