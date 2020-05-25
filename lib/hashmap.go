package lib

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

var defaultTTL = 5         // seconds
var defaultMaxSize = 10000 // items in map

// HashMap is a thread safe implementation that keeps deleting items after
type HashMap struct {
	MaxSize int
	TTL     int
	Data    map[string]MarcoPoloMessage
	L       sync.Mutex
}

// MarcoPoloMessage is the json message structure of incomming messages
type MarcoPoloMessage struct {
	ID        string            `json:"id"`
	Lat       float32           `json:"lat,omitempty"`
	Long      float32           `json:"long,omitempty"`
	Tags      map[string]string `json:"tags,omitempty"`
	Timestamp int64             `json:"-"`
}

// NewHashMapDefault returns a new thread safe HashMap with default config
func NewHashMapDefault() *HashMap {
	return NewHashMap(defaultTTL, defaultMaxSize)
}

// NewHashMap returns a new thread safe HashMap
func NewHashMap(ttl int, size int) *HashMap {
	if ttl == 0 {
		ttl = defaultTTL
	}
	if size == 0 {
		size = defaultMaxSize
	}
	hashMap := HashMap{MaxSize: size, TTL: ttl, Data: make(map[string]MarcoPoloMessage)}
	return &hashMap
}

// parseJSONString parses a Json string to type MarcoPolo
func parseJSONString(message string) (*MarcoPoloMessage, error) {
	res := MarcoPoloMessage{}
	err := json.Unmarshal([]byte(message), &res)
	return &res, err
}

// Insert puts key into the map and removes it after TTL
func (x *HashMap) Insert(json string) {
	x.L.Lock()
	defer x.L.Unlock()
	fmt.Println("Inserting " + json)
	message, err := parseJSONString(json)
	if err != nil {
		fmt.Println("Insertion failed" + err.Error())
	}
	if len(x.Data) >= x.MaxSize {
		// TODO: improve this logic
		fmt.Println("Insertion failed, Map too full.")
		return
	}
	x.Data[message.ID] = *message
	x.RemoveAfterTTL(message.ID, x.TTL)
	fmt.Println("Insert SUCCESS")
}

// RemoveAfterTTL removes the given id from HashMap after the given TTL
func (x *HashMap) RemoveAfterTTL(id string, ttl int) {
	timer := time.NewTimer(time.Duration(ttl) * time.Second)
	go func() {
		<-timer.C
		x.Remove(id)
	}()
}

// Remove removes the given id from HashMap
func (x *HashMap) Remove(id string) {
	x.L.Lock()
	defer x.L.Unlock()
	fmt.Println("Removing " + id)
	_, ok := x.Data[id]
	if ok {
		delete(x.Data, id)
	}
}

// FlushMap flushes current map to a string
func (x *HashMap) FlushMap() string {
	x.L.Lock()
	defer x.L.Unlock()
	ret, _ := json.Marshal(x.Data)
	return string(ret)
}
