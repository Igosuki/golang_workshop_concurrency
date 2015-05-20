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

### Exercice 4 : Naive poll loop

Useful libs : 
- Logrus : github.com/Sirupsen/logrus

Put ping pong and rss in their own files, create a run script with 4 go procs that runs the program including all the source files (just using the main package).

Create a feed function that can launch subs for a list of domains and close them eventually.
You may use response request patterns in the form of chan chan struct{}

- Implement a simple Fetch with the http package and xml decoder
- Start implementing a naive update loop on the sub struct that manages the fetching. This loop will be called in a routine.

The loop :
1. Fetch
2. Send items on the subscription feed
3. Exit when close is called and report errors

### Several bugs are present in this loop

Race conditions when closing. (use go run -race)
Sleeps may keep the loop running. (can you guess why ?)
It may block on the feed.

### Exercice 5 : Select -> The answer to concurrent events 

The 3 features of the loop should each become a statement, which can only be push or receive channel statements.

Reminder : 
1. A close order was made, the loop should clean up and return.
2. It's time to fetch
3. Something is read from the subscription feed

#### Step 1 The Close channel
You are sending a close command to subscription loop routines. Use a request response structure for the close field (I send a request and read from the response).

N.B: if you are using multiple fetchers routines for a single sub (aka merge) you may want to close that channel after all routines have closed.

functions : make, close

#### Step 2 Fetching after a delay
Since you can't just fetch and push into the feed, you'll have to create an intermediary array of pending items waiting to be read from the feed.

You also have to store the next time you will fetch, and determine whether or not it has expired, in which case you have to fetch immediatly.

packages : time
functions : append

### Step 3 Sending items down the feed

Read the first pending item down the feed, then splice it from pending items.

** Yet another crash ** 

Jeez, the cake was a lie, I still get locked out...

### Exercice 6 : Fix select statements 

Techniques : nil channels, a nil channel is just considered a non passing case.

** Solution **
Only push in the feed if there is at least one pending item

### Exercice 7 : Duplicates

Fetch may return duplicates, maintain a register of a unique id and check that for all items

### Exercice 8 : Unbound pending items

Manage memory by having a max amount of items in the queue (you can pop it as fifo, you can use an int).





