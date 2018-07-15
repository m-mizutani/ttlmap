package ttlmap

import (
	"bytes"
)

// Map is a main data structure of ttlmap.
type Map struct {
	table   map[hashValue]*elementBucket
	frames  []timeFrame
	current tick
}

type tick uint
type elementKey []byte

// New is a constructor of Map
func New(maxTick tick) *Map {
	m := new(Map)
	m.table = make(map[hashValue]*elementBucket)
	m.frames = make([]timeFrame, maxTick+1)
	return m
}

func (x *Map) tickToPtr(t tick) tick {
	return t % tick(len(x.frames))
}

// Set method puts value into map table with Time To Live
func (x *Map) Set(key []byte, value interface{}, ttl tick) error {
	maxTick := tick(len(x.frames))
	if ttl >= maxTick {
		return &ErrOverMaxTick{max: maxTick, arg: ttl}
	}

	hv := fnvHash(key)
	bucket := x.table[hv]
	if bucket == nil {
		bucket = newEleementBucket()
		x.table[hv] = bucket
	}

	if elem := bucket.search(key); elem != nil {
		return &ErrDuplicatedKey{key}
	}

	elem := newElement(key, value, ttl)
	bucket.insert(elem)

	ptr := x.tickToPtr(x.current + ttl)
	x.frames[ptr].addElement(elem)

	return nil
}

// Get method retrieves value from table.
func (x *Map) Get(key []byte) (value interface{}) {
	hv := fnvHash(key)
	bucket := x.table[hv]
	if bucket == nil {
		return
	}

	elem := bucket.search(key)
	if elem == nil {
		return
	}

	value = elem.value
	return
}

// Prune method updates current tick and remove elements that are expired.
func (x *Map) Prune(update tick) (values []interface{}) {
	for t := x.current; t < x.current+update; t++ {
		p := x.tickToPtr(t)
		elements := x.frames[p].purge()

		for _, elem := range elements {
			values = append(values, elem.value)
		}
	}

	x.current += update
	return
}

// -----------------------------------
// Internal data structures
//
type timeFrame struct {
	root element
}

func (x *timeFrame) addElement(e *element) {
	e.timeoutLink = x.root.timeoutLink
	x.root.timeoutLink = e
}

func (x *timeFrame) purge() (elements []*element) {
	for elem := x.root.timeoutLink; elem != nil; elem = elem.timeoutLink {
		elements = append(elements, elem)
		elem.detach()
	}

	return
}

type elementBucket struct {
	root element
}

func newEleementBucket() *elementBucket {
	b := new(elementBucket)
	return b
}

func (x *elementBucket) insert(e *element) {
	x.root.attach(e)
}

func (x *elementBucket) search(key elementKey) *element {
	for elem := x.root.next; elem != nil; elem = elem.next {
		if elem.matchKey(key) {
			return elem
		}
	}

	return nil
}

type element struct {
	next, prev  *element
	timeoutLink *element
	value       interface{}
	key         elementKey
	ttl         tick
	last        tick
}

func newElement(key []byte, value interface{}, ttl tick) *element {
	elem := &element{value: value, key: key, ttl: ttl}
	return elem
}

func (x *element) attach(e *element) {
	next := x.next
	x.next, e.prev = e, x

	if next != nil {
		next.prev, e.next = e, next
	}
}

func (x *element) detach() {
	if x.prev != nil {
		x.prev.next = x.next
	}
	if x.next != nil {
		x.next.prev = x.prev
	}

	x.prev, x.next = nil, nil
}

func (x *element) matchKey(key elementKey) bool {
	return bytes.Equal(x.key, key)
}

func (x *element) equals(e *element) bool {
	return x.matchKey(e.key)
}
