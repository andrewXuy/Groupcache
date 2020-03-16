package cache

import (
	"cache/singleflight"
	"fmt"
	"log"
	"sync"
	pb "cache/cachepb"
)

type Group struct {
	name      string
	maincache cache
	// A interface fetch data
	fetcher Fetcher
	peers   PeerPicker
	// Make sure each key is only fetched once
	loader 	*singleflight.Group
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// Create new instance of Group
func NewGroup(name string, cacheBytes64 int64, fetcher Fetcher) *Group {
	if fetcher == nil {
		panic("Cant have nil Fetcher")
	}
	mu.Lock()
	defer mu.Unlock()
	group := &Group{
		name:      name,
		maincache: cache{cacheBytes64: cacheBytes64},
		fetcher:   fetcher,
		loader: &singleflight.Group{},
	}
	groups[name] = group
	return group
}

// Return specified group using name
func GetGroup(name string) *Group {
	mu.Lock()
	defer mu.Unlock()
	g := groups[name]
	return g

}

// Implement Group level get value
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("Valid Key is required")
	}
	if v, ok := g.maincache.get(key); ok {
		log.Println("[Group] hit")
		return v, nil
	}
	return g.load(key)
}

// Set the peer for Group
func (g *Group) RegisterPeers(peer PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeers call more than once")
	}
	g.peers = peer
}

// Load from remote peer
func (g *Group) getFromPeer(peer peerGetter, key string) (ByteView, error) {
	req := &pb.Request{
		Group:g.name,
		Key:key,
	}
	res := &pb.Response{}
	err := peer.Get(req, res)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{res.Value}, nil
}

// Right now it is just fetch data from local
func (g *Group) load(key string) (value ByteView, err error) {
		viewi, err := g.loader.Do(key, func() (i interface{}, err error) {
			if g.peers != nil {
			if peer, ok := g.peers.PickPeer(key); ok {
				if value, err := g.getFromPeer(peer, key); err == nil {
					return value, nil
				}
				log.Println("[Group] failed to fetch from peer ", err)
			}

		}
		return g.getLocal(key)
	})
	if err == nil {
		return viewi.(ByteView), nil
	}
	return

}

// Fetch data locally
func (g *Group) getLocal(key string) (ByteView, error) {
	bytes, err := g.fetcher.Fetch(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: clone(bytes)}
	g.addToCache(key, value)
	return value, nil

}
func (g *Group) addToCache(key string, data ByteView) {
	g.maincache.add(key, data)
}

// Fetch data from other side if data not in the cache
type Fetcher interface {
	Fetch(key string) ([]byte, error)
}

type FetcherFunc func(key string) ([]byte, error)

// Implement a function using Fetcher interface
func (f FetcherFunc) Fetch(key string) ([]byte, error) {
	return f(key)
}
