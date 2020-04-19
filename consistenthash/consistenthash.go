package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

type ConsistentHash interface {
	Add(node string)
	Remove(node string)
	Get(key string) string
}

type consistentHashImpl struct {
	number    int
	hashSlice []uint32
	hashMap   map[uint32]string
	sync.RWMutex
}

func New(number int) ConsistentHash {
	consistentHashInstance := &consistentHashImpl{number: number, hashSlice: make([]uint32, 0), hashMap: make(map[uint32]string, 0)}
	return consistentHashInstance
}

func (c *consistentHashImpl) Add(node string) {
	c.Lock()
	defer c.Unlock()

	for i := 0; i < c.number; i++ {
		virtualNode := node + `_virtual_node_` + strconv.Itoa(i)
		hashData := crc32.ChecksumIEEE([]byte(virtualNode))
		c.hashSlice = append(c.hashSlice, hashData)
		sort.Slice(c.hashSlice, func(i int, j int) bool {
			return c.hashSlice[i] < c.hashSlice[j]
		})

		c.hashMap[hashData] = node
	}

	return
}
func (c *consistentHashImpl) Remove(node string) {
	c.Lock()
	defer c.Unlock()

	for i := 0; i < c.number; i++ {
		virtualNode := node + `_virtual_node_` + strconv.Itoa(i)
		hashData := crc32.ChecksumIEEE([]byte(virtualNode))

		remove(c.hashSlice, hashData)
		delete(c.hashMap, hashData)
	}
	return
}

func remove(slice []uint32, elems uint32) []uint32 {
	for i := range slice {
		if slice[i] == elems {
			slice = append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func (c *consistentHashImpl) Get(key string) string {
	c.RLock()
	defer c.RUnlock()

	hashData := crc32.ChecksumIEEE([]byte(key))
	pos := sort.Search(len(c.hashSlice), func(i int) bool {
		return c.hashSlice[i] >= hashData
	})

	if pos >= len(c.hashSlice) {
		pos = 0
	}

	return c.hashMap[c.hashSlice[pos]]
}
