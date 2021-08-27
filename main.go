package gumball

import (
	"fmt"
	"sync"
	"time"
)

type Dispenser struct {
	ttl         time.Duration
	expiration  time.Time
	item        Gumball
	refreshFunc RefreshFunc
	sync.Mutex
}

type Gumball interface{}
type RefreshFunc func() (Gumball, error)

func NewDispenser(ttl time.Duration, refreshFunc RefreshFunc) *Dispenser {
	return &Dispenser{
		ttl:         ttl,
		refreshFunc: refreshFunc,
	}
}

func (d *Dispenser) Dispense() (Gumball, error) {
	d.Lock()
	defer d.Unlock()

	now := time.Now()
	if now.After(d.expiration) {
		fmt.Printf("Cache invalid. ttl=%v, expiration=%v... refreshing!\n", d.ttl, d.expiration.Format(time.Stamp))
		var err error
		d.item, err = d.refreshFunc()
		if err != nil {
			return nil, err
		}
		d.expiration = now.Add(d.ttl)
	}
	return d.item, nil
}
