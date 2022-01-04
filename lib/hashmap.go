package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

var defaultTTL = 86400    // seconds
var defaultMaxSize = 1000 // items in map

// HashMap is a thread safe implementation that keeps deleting items after
type HashMap struct {
	IpDb    *IpDb
	MaxSize int
	TTL     int
	Data    map[string]MarcoPoloMessage
	L       sync.Mutex
}

// MarcoPoloMessage is the json message structure of incomming messages
type MarcoPoloMessage struct {
	Ip        string            `json:"ip"`
	Src       string            `json:"src,omitempty"`
	Lat       float32           `json:"lat,omitempty"`
	Long      float32           `json:"long,omitempty"`
	Tags      map[string]string `json:"tags,omitempty"`
	Timestamp int64             `json:"-"`
}

// GeoJson message structure
type GeoJson struct {
	Type       string            `json:"type"`
	Geometry   Geometry          `json:"geometry"`
	Properties map[string]string `json:"properties"`
}

// Geometry object in GeoJson
type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float32 `json:"coordinates"`
}

// NewHashMapDefault returns a new thread safe HashMap with default config
func NewHashMapDefault(ipdb *IpDb) *HashMap {
	return NewHashMap(defaultTTL, defaultMaxSize, ipdb)
}

// NewHashMap returns a new thread safe HashMap
func NewHashMap(ttl int, size int, ipdb *IpDb) *HashMap {
	if ttl == 0 {
		ttl = defaultTTL
	}
	if size == 0 {
		size = defaultMaxSize
	}
	hashMap := HashMap{IpDb: ipdb, MaxSize: size, TTL: ttl, Data: make(map[string]MarcoPoloMessage)}
	log.Println(fmt.Sprintf("HashMap initialized with [MaxSize: %d, TTL: %d]", size, ttl))
	return &hashMap
}

// ParseJSONString parses a Json string to type MarcoPolo
func ParseJSONString(message string) (*MarcoPoloMessage, string, error) {
	res := MarcoPoloMessage{Src: "default"}
	err := json.Unmarshal([]byte(message), &res)
	src := res.Src
	return &res, src, err
}

// Insert puts key into the map and removes it after TTL
func (x *HashMap) Insert(message *MarcoPoloMessage) {
	x.L.Lock()
	defer x.L.Unlock()

	if len(x.Data) >= x.MaxSize {
		// TODO: improve this logic
		log.Println("Insertion failed, Map too full.")
		return
	}
	err := enrichMessage(message, x)
	if err != nil {
		log.Println("Message not enriched", err)
		if message.Lat == 0 && message.Long == 0 {
			return
		}
	}
	x.Data[message.Ip] = *message
	x.RemoveAfterTTL(message.Ip, x.TTL)
	// fmt.Println("Insert SUCCESS")
}

// InsertString for inserting json strings
func (x *HashMap) InsertString(json string) {
	// log.Println("Inserting " + json)
	message, _, err := ParseJSONString(json)
	if err != nil {
		log.Println("Insertion failed" + err.Error())
		return
	}
	x.Insert(message)
}

func enrichMessage(message *MarcoPoloMessage, x *HashMap) error {
	val, err := x.IpDb.Lookup(message.Ip)
	if err != nil {
		return err
	}
	if val.Latitude == 0 && val.Longitude == 0 {
		return errors.New("Invalid/Private IP: " + message.Ip)
	}
	message.Lat = val.Latitude
	message.Long = val.Longitude
	if message.Tags == nil {
		message.Tags = make(map[string]string)
	}
	message.Tags["city"] = val.City
	message.Tags["country"] = val.Country_short
	// log.Println("Enriched data point", message)
	return nil
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
	// fmt.Println("Removing " + id)
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

// FlushGeoJson flushes current map in a GeoJson format
func (x *HashMap) FlushGeoJson() string {
	arr := x.GetGeoJson()
	ret, _ := json.Marshal(arr)
	return string(ret)
}

// GetGeoJson fetches list of current geoJsonObjects
func (x *HashMap) GetGeoJson() []*GeoJson {
	x.L.Lock()
	defer x.L.Unlock()
	arr := []*GeoJson{}
	for _, v := range x.Data {
		geoJson := getGeoJsonObject(v)
		arr = append(arr, geoJson)
	}
	return arr
}

// getGeoJsonObject constructs a GeoJson object
func getGeoJsonObject(v MarcoPoloMessage) *GeoJson {
	geoJson := GeoJson{}
	geoJson.Type = "Feature"
	geoJson.Geometry.Type = "Point"
	geoJson.Geometry.Coordinates = make([]float32, 2)
	geoJson.Geometry.Coordinates[0] = v.Long
	geoJson.Geometry.Coordinates[1] = v.Lat
	if v.Tags != nil {
		geoJson.Properties = v.Tags
	} else {
		geoJson.Properties = make(map[string]string)
	}
	geoJson.Properties["ip"] = v.Ip
	return &geoJson
}
