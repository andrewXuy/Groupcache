package cache

import (
	"fmt"
	"cache/consistenthash"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	pb "cache/cachepb"
	"github.com/golang/protobuf/proto"
)

const (
	defaultBase      = "/_cache/"
	defaultreplicate = 50
)

type HTTPPool struct {
	node     string
	basePath string
	mu sync.Mutex
	// Using to choose key node
	peers *consistenthash.Map
	// <key Node, httpGetter>
	httpGetter map[string]*httpGetter // key e.g. "http://10.0.0.1:8008"
}
type httpGetter struct {
	// Remote node addr
	baseURL string
}
// Creat node for each http client
func (h *httpGetter) Get(in *pb.Request, out *pb.Response) error {
	u := fmt.Sprintf("%v%v/%v", h.baseURL, url.QueryEscape(in.GetGroup()), url.QueryEscape(in.GetKey()))
	res, err := http.Get(u)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server return: %v", res.Status)
	}
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("reading body : %v", err)
	}
	if err = proto.Unmarshal(bytes,out); err != nil {
		return fmt.Errorf("decoding response body %v",err)
	}
	return  nil
}

func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = consistenthash.New(defaultreplicate,nil)
	// Ravel peers as stream into hashmap
	p.peers.Add(peers...)
	p.httpGetter = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		p.httpGetter[peer] = &httpGetter{baseURL:peer + p.basePath}
	}

}
// Pick a peer according to its key
func (p *HTTPPool) PickPeer(key string) (peerGetter,  bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer:= p.peers.Get(key); peer!="" && peer != p.node {
		p.Log("Pick peer %s\n", peer)
		return p.httpGetter[peer], true
	}
	return nil , false
}

func NewPool(node string) *HTTPPool {
	return &HTTPPool{
		node:     node,
		basePath: defaultBase,
	}
}


// Log info of server name implementing Log interface
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s ", p.node, fmt.Sprintf(format, v...))

}

func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTP pool servering unexpected path: " + r.URL.Path)
	}
	p.Log("%s  %s", r.Method, r.URL.Path)
	// <basePath>/<groupName>/<key> all of them are needed!
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	groupName := parts[0]
	key := parts[1]

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no group exist: "+groupName, http.StatusNotFound)
		return
	}
	data, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	body ,err := proto.Marshal(&pb.Response{Value:data.ByteSlice()})
	if err != nil {
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(body)

}
