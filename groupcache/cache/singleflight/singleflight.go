package singleflight


import "sync"

// call is an in-air or complete DO call
type call struct {
	wg sync.WaitGroup
	val interface{}
	err error
}

// Group represents a class of work and forms a namespace in which
// units of work can be executed with duplicate suppression
type Group struct {
	mu sync.Mutex
	m map[string]*call // Lazy initialized
}

// No matter how many times call to DO , the same key will only execute the fn once
func (g *Group) Do(key string, fn func()(interface{}, error))(interface{},error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()// if is execting, wait
		return c.val, c.err// request end, return result
	}
	c:= new(call)
	c.wg.Add(1)// add lock before requesting
	g.m[key] = c// show key is in handling
	g.mu.Unlock()

	c.val , c.err = fn() // call function
	c.wg.Done() // Request done

	g.mu.Lock()
	// update map
	delete(g.m,key)
	g.mu.Unlock()
	return c.val , c.err

}