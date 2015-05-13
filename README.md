## Golang Concurrency Workshop

Use git checkout tags/exo-{exercice number} for solution and instructions

### Exercice 1 : Pingpong

Implement a simple program with two player functions sending a ball to each other.
Every time the function receives a ball, log the number of hits.

You need : 
- A Ball struct
- A player function
- A table 

Finish the game with a timeout and reading a last time from the table.

### Exercice 2 : Simulate a Deadlock

Add an option to your program to simluate the pingpong game failing with a deadlock.

### Exercice 3 : Bootstrap the design of an rss feed reader

Model public and private structs and interfaces for :
- A feed of items
- A Fetcher that can fetch and return a list of items, the next time it'll fetch, any errors
- A concrete fetcher that fetches now
- A Fetch taking a domain and returning a Fetcher
- A Subscription that can be updated for items, and closed
- A concrete subscription
- A Subscribe function that returns a subscription
- A merge function that takes a list of subs and returns a new sub

be minimalistic, and use the zero initialization pattern (everything should work out of the box after an &struct{}{})

