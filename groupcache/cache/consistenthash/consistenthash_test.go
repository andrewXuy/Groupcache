package consistenthash

import (
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	hash := New(3, func(data []byte) uint32 {
		i, _ := strconv.Atoi(string(data))
		return uint32(i)
	})
	hash.Add("6","4","2")
	testCast := map[string]string{
		"2":"2",
		"11":"2",
		"23":"4",
		"27":"2",
	}

	for k,v := range testCast {
		if hash.Get(k) != v {
			t.Errorf("The hash value should be %s but get %s\n", v, hash.Get(k))
		}
		hash.Add("8")
		testCast["27"] = "8"
		for k, v := range testCast {
			if hash.Get(k) != v {
				t.Errorf("The hash value should be %s but get %s\n", v, hash.Get(k))
			}
		}
	}
}

