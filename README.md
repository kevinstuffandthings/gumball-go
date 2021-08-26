# Gumball (Golang edition)

A Go translation of the Ruby [Gumball](https://github.com/kevinstuffandthings/gumball) gem.

_I wrote this as a way to learn Go. Both my learning and this package are a work-in-progress. Please be extremely careful if you consider
using (or even looking) at this code!_

## Usage
Let's say we have some expensive operation we need to utilize the value of. We need to refresh it occasionally, but if we sometimes get a
slightly-stale copy, that's ok.

```go
func expensiveOperation() (int, error) {
    rand.Seed(time.Now().UnixNano())
    time.Sleep(5 * time.Second)
    return rand.Intn(100), nil
}
```

We set up a new dispenser for that operation, and initialize it with a TTL and refresh function:

```go
dispenser = gumball.NewDispenser(300*time.Second, func() (gumball.Gumball, error) {
    value, err := expensiveOperation()
    return gumball.Gumball(value), err
}

gb, err := dispenser.Dispense() // this will take a while...
value, ok := gb.(int) // check that this went ok, of course!
fmt.Println("Value is", value)
```

Subsequent calls to `dispenser.Dispense()` within the specified TTL will return the same value, immediately.
Once the TTL has come and gone, the next call will kick off our `expensiveOperation`.

## Example
```go
package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/kevinstuffandthings/gumball-go"
)

func expensiveOperation() (int, error) {
	rand.Seed(time.Now().UnixNano())
	time.Sleep(3 * time.Second)
	return rand.Intn(100), nil
}

func main() {
	dispenser := gumball.NewDispenser(5*time.Second, func() (gumball.Gumball, error) {
		return expensiveOperation()
	})

	iterations := 20
	ticker := time.NewTicker(1 * time.Second)
	for {
		<-ticker.C
		gb, _ := dispenser.Dispense() // of course you should check this...
		value, _ := gb.(int)          // ...and this!
		fmt.Printf("%v: got value from gumball: %d\n", time.Now().Format("15:04:05"), value)

		iterations--
		if iterations == 0 {
			break
		}
	}
}
```

#### Output
```
Cache invalid. ttl=5s, expiration=0001-01-01 00:00:00 +0000 UTC... refreshing!
18:33:07: got value from gumball: 90
18:33:07: got value from gumball: 90
18:33:08: got value from gumball: 90
18:33:09: got value from gumball: 90
Cache invalid. ttl=5s, expiration=2021-08-26 18:33:09.921185 -0400 EDT m=+6.005250071... refreshing!
18:33:13: got value from gumball: 42
18:33:13: got value from gumball: 42
18:33:14: got value from gumball: 42
18:33:15: got value from gumball: 42
Cache invalid. ttl=5s, expiration=2021-08-26 18:33:15.921125 -0400 EDT m=+12.005217875... refreshing!
18:33:19: got value from gumball: 41
18:33:19: got value from gumball: 41
18:33:20: got value from gumball: 41
18:33:21: got value from gumball: 41
Cache invalid. ttl=5s, expiration=2021-08-26 18:33:21.921112 -0400 EDT m=+18.005232767... refreshing!
18:33:25: got value from gumball: 98
18:33:25: got value from gumball: 98
18:33:26: got value from gumball: 98
Cache invalid. ttl=5s, expiration=2021-08-26 18:33:27.916175 -0400 EDT m=+24.000322463... refreshing!
18:33:30: got value from gumball: 7
18:33:30: got value from gumball: 7
18:33:31: got value from gumball: 7
18:33:32: got value from gumball: 7
Cache invalid. ttl=5s, expiration=2021-08-26 18:33:32.921068 -0400 EDT m=+29.005238501... refreshing!
18:33:36: got value from gumball: 73
```
