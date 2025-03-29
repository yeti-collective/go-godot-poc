package objects

import "maps"

import "sync"

// A generic, thread-safe map of objects with auto-incrementing IDs.
// Generic type T allows us to use this for any type we want (e.g. users)
type SharedCollection[T any] struct {
	objectsMap map[uint64]T
	nextId     uint64

	// Mutex is short for "mutual exclusion" -- it allows you to lock the code between mapMux.Lock()
	// and mapMux.Unlock() so that it can only be run by one goroutine at a time. This is what we
	// will use below in order to ensure that only one thread can make a new id at a time. If we didn't do
	// this, then if two people logged in close enough together, they might end up with the same id and
	// appear as the same person.
	// Check this out for a good summary: https://www.geeksforgeeks.org/mutex-in-golang-with-examples/
	mapMux sync.Mutex
}

// Constructor
func NewSharedCollection[T any](capacity ...int) *SharedCollection[T] {
	var newObjMap map[uint64]T

	// Optionally set a capacity
	if len(capacity) > 0 {
		newObjMap = make(map[uint64]T, capacity[0])
	} else {
		newObjMap = make(map[uint64]T)
	}

	return &SharedCollection[T]{
		objectsMap: newObjMap,
		nextId:     1,
	}
}

// Add an object to the map with the given ID (if provided) or the next available ID.
// Returns the ID of the object added.
func (s *SharedCollection[T]) Add(obj T, id ...uint64) uint64 {

	// This locks the thread, so this id process will only occur on one goroutine at a time
	s.mapMux.Lock()

	// Defer in Go will call this at the close of the function. This is cool because it will call
	// if the function just runs out, in the case of any return, or if there's an error thrown. Extra
	// useful right here because otherwise we would fail to unlock the mutex, which would effectively
	// turn our multithreaded multiplayer engine into a single-threaded system where every other
	// goroutine for every other player would permanently lose access to the instance!
	defer s.mapMux.Unlock()

	// Set this item to the next sequential id ...
	thisId := s.nextId

	// ... unless we manually pass in an id to set.
	if len(id) > 0 {
		thisId = id[0]
	}

	// Link it to the map
	s.objectsMap[thisId] = obj

	// Increment the next id to prep for the next client connection
	s.nextId++

	// And return (which will trigger the defer, and unlock the mutex).
	return thisId
}

// Remove removes an object from the map by ID, if it exists.
func (s *SharedCollection[T]) Remove(id uint64) {
	s.mapMux.Lock()
	defer s.mapMux.Unlock()

	// We mutex lock this operation as well to prevent multiple goroutines trying to mutate the map
	// at the same time, functionally preventing race conditions
	delete(s.objectsMap, id)
}

// Call the callback function for each object in the map.
func (s *SharedCollection[T]) ForEach(callback func(uint64, T)) {

	// Create a local copy while holding the lock. We do this again to avoid race conditions -- we don't
	// want to be doing logic on the loop and then have another goroutine mutate the list mid-iteration!
	s.mapMux.Lock()
	localCopy := make(map[uint64]T, len(s.objectsMap))
	maps.Copy(localCopy, s.objectsMap)
	s.mapMux.Unlock()

	// Iterate over the local copy without holding the lock.
	for id, obj := range localCopy {
		callback(id, obj)
	}
}

// Get the object with the given ID, if it exists, otherwise nil.
// Also returns a boolean indicating whether the object was found.
func (s *SharedCollection[T]) Get(id uint64) (T, bool) {
	s.mapMux.Lock()
	defer s.mapMux.Unlock()

	obj, ok := s.objectsMap[id]
	return obj, ok
}

// Get the approximate number of objects in the map.
// The reason this is approximate is because we don't lock the map to get the length.
// Locking the map incurs a performance hit, so we only do that if we actually need to for data security --
// map length doesn't really "matter" in the same way as a direct reference or mutation does,
// so there's no real need to lock it before accessing.
func (s *SharedCollection[T]) Len() int {
	return len(s.objectsMap)
}
