package consistenthash

import ("hash/crc32"
"sort"
"strconv"
)

type Hash func(data []byte) uint32
type Map struct {
	hashfunc Hash
	// virtual node to handle imbalance question
	replicate int
	keys [] int
	hashMap map[int]string
}
// Initialization a HashMap with crc32 as default function
func New(replicate int, fn Hash) *Map {
	m := &Map{
		hashfunc:  fn,
		replicate: replicate,
		hashMap:   make(map[int]string),
	}
	if fn == nil {
		m.hashfunc = crc32.ChecksumIEEE
	}
	return m
}

// Add nodes from true node(From Ip addr and machine ID) to virtual node into Map
func (m *Map) Add(keys ...string) {
	for _ ,key:= range keys {
		for i:=0;i<m.replicate;i++ {
			hash := int(m.hashfunc([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	// sort the node in hash loop
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hash := int(m.hashfunc([]byte(key)))
	// Binary Search
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	return m.hashMap[m.keys[idx%len(m.keys)]]
}