package cache

import (
	"fmt"
	"log"
	"reflect"
"testing")

var db = map[string]string {
	"Harry" : "123",
	"Bob" : "23",
	"Google" : "100",
	"Leetcode" :"1",
}

func TestFetcher(t *testing.T) {
	var f Fetcher = FetcherFunc(func(key string) ([]byte,error){
		return []byte(key),nil
	})
	expect := []byte("key")
	if v, _ := f.Fetch("key"); !reflect.DeepEqual(v,expect) {
		t.Errorf("Fetch data failed")
	}
}
func TestFetch(t *testing.T) {
	load := make(map[string]int, len(db))
	g:=NewGroup("rank", 1024, FetcherFunc(
		func(key string) ([]byte,error){
			log.Println("[DB] search Key",key)
			if v,ok := db[key]; ok {
				// Initialize
				if _,ok := load[key]; !ok{
					load[key] = 0
				}
				load[key] ++
				return []byte(v),nil
			}
			return nil, fmt.Errorf("%s not exist",key)
		}))
	for key,value := range db {
		if temp,err := g.Get(key); err!=nil || temp.String() != value {
			t.Fatalf("Failed get Value of %s\n",key)
		}
		if _,err := g.Get(key); err!=nil || load[key]>1 {
			t.Fatalf("cache miss %s\n",key)
		}

	}
	if temp, err := g.Get("saodijasofvnov");err ==nil {
		t.Fatalf("Unkown value get valid value %s\n",temp)
	}
}

func TestGetGroup(t *testing.T) {
	NewGroup("rank", 1024, FetcherFunc(
		func(key string) (k []byte ,e error){
			return
		}))
	if g:= GetGroup("rank");g == nil || g.name != "rank" {
		t.Fatalf("rank group not exist\n")
	}
	if g:=GetGroup("rank12345");g!=nil{
		t.Fatalf("unset group exist")
	}

}