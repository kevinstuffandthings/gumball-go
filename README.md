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

We set up a new dispenser for that operation. These dispensers are best saved as class variables, so the dispenser itself is a singleton.

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
