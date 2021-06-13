package singleflight

import "sync"

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait() // if request is running, wait
		return c.val, c.err
	}
	c := new(call)
	c.wg.Add(1) // add wait group when request
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn() // request
	c.wg.Done()         // notify

	g.mu.Lock()
	delete(g.m, key) // update
	g.mu.Unlock()

	return c.val, c.err
}
