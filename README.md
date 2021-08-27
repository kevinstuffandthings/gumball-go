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

    gumball "github.com/kevinstuffandthings/gumball-go"
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
Cache invalid. ttl=5s, expiration=Jan  1 00:00:00... refreshing!
08:27:36: got value from gumball: 60
08:27:36: got value from gumball: 60
08:27:37: got value from gumball: 60
Cache invalid. ttl=5s, expiration=Aug 27 08:27:38... refreshing!
08:27:41: got value from gumball: 69
08:27:41: got value from gumball: 69
08:27:42: got value from gumball: 69
08:27:43: got value from gumball: 69
Cache invalid. ttl=5s, expiration=Aug 27 08:27:43... refreshing!
08:27:47: got value from gumball: 8
08:27:47: got value from gumball: 8
08:27:48: got value from gumball: 8
08:27:49: got value from gumball: 8
Cache invalid. ttl=5s, expiration=Aug 27 08:27:49... refreshing!
08:27:53: got value from gumball: 21
08:27:53: got value from gumball: 21
08:27:54: got value from gumball: 21
08:27:55: got value from gumball: 21
Cache invalid. ttl=5s, expiration=Aug 27 08:27:55... refreshing!
08:27:59: got value from gumball: 34
08:27:59: got value from gumball: 34
08:28:00: got value from gumball: 34
08:28:01: got value from gumball: 34
Cache invalid. ttl=5s, expiration=Aug 27 08:28:01... refreshing!
08:28:05: got value from gumball: 52
```
