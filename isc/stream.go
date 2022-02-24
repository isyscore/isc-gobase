package isc

import (
	"sort"
)

// A Stream is a stream that can be used to do stream processing.
type Stream[T any] struct {
	source <-chan T
}

// StreamJust converts the given arbitrary items to a Stream.
func StreamJust[T any](items ...T) Stream[T] {
	source := make(chan T, len(items))
	for _, e := range items {
		source <- e
	}
	close(source)
	return StreamRange(source)
}

// AllMatch returns whether all elements of this stream match the provided predicate.
// May not evaluate the predicate on all elements if not necessary for determining the result.
// If the stream is empty then true is returned and the predicate is not evaluated.
func (s Stream[T]) AllMatch(predicate func(T) bool) bool {
	for item := range s.source {
		if !predicate(item) {
			// make sure the former goroutine not block, and current func returns fast.
			go StreamDrain(s.source)
			return false
		}
	}
	return true
}

// AnyMatch returns whether any elements of this stream match the provided predicate.
// May not evaluate the predicate on all elements if not necessary for determining the result.
// If the stream is empty then false is returned and the predicate is not evaluated.
func (s Stream[T]) AnyMatch(predicate func(T) bool) bool {
	for item := range s.source {
		if predicate(item) {
			// make sure the former goroutine not block, and current func returns fast.
			go StreamDrain(s.source)
			return true
		}
	}
	return false
}

// NoneMatch returns whether all elements of this stream don't match the provided predicate.
// May not evaluate the predicate on all elements if not necessary for determining the result.
// If the stream is empty then true is returned and the predicate is not evaluated.
func (s Stream[T]) NoneMatch(predicate func(T) bool) bool {
	for item := range s.source {
		if predicate(item) {
			go StreamDrain(s.source)
			return false
		}
	}
	return true
}

func (s Stream[T]) Map(fn func(T) T) Stream[T] {
	source := make(chan T)
	for item := range s.source {
		source <- fn(item)
	}
	return StreamRange(source)
}

// Sort sorts the items from the underlying source.
func (s Stream[T]) Sort(less func(T, T) bool) Stream[T] {
	var items []T
	for item := range s.source {
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		return less(items[i], items[j])
	})
	return StreamJust(items...)
}

// Distinct removes the duplicated items base on the given keySelector.
func (s Stream[T]) Distinct(keySelector func(T) T) Stream[T] {
	source := make(chan T)
	m := make(map[any]T)
	for item := range s.source {
		m[keySelector(item)] = item
	}
	for _, v := range m {
		source <- v
	}
	return StreamRange(source)
}

// ForEach seals the Stream with the fn on each item, no successive operations.
func (s Stream[T]) ForEach(fn func(T)) {
	for item := range s.source {
		fn(item)
	}
}

//FirsVal returns the first element,channel is FIFO,so first goroutine will get head element or nil
func (s Stream[T]) FirsVal() any {
	for item := range s.source {
		go StreamDrain(s.source)
		return item
	}
	return nil
}

//First returns the first element,channel is FIFO,so first goroutine will get head element
func (s Stream[T]) First(valueSelector func(T) bool) Stream[T] {
	source := make(chan T)
	go func() {
		for item := range s.source {
			if valueSelector(item) {
				source <- item
				break
			}
		}
	}()
	close(source)
	return StreamRange(source)
}

// LastVal returns the last item, or nil if no items.
func (s Stream[T]) LastVal() (item T) {
	for item = range s.source {
	}
	return
}

// Last returns the last item, or nil if no items.
func (s Stream[T]) Last(valueSelector func(any) bool) Stream[T] {
	source := make(chan T)
	var lastValue T
	for item := range s.source {
		if valueSelector(item) {
			lastValue = item
		}
	}
	source <- lastValue
	close(source)
	return StreamRange(source)
}

//Filter Returns a list containing only elements matching the given predicate.
func (s Stream[T]) Filter(predicate func(T) bool) Stream[T] {
	source := make(chan T)
	for item := range s.source {
		if predicate(item) {
			source <- item
		}
	}
	defer close(source)
	return StreamRange(source)
}

func StreamRange[T any](source <-chan T) Stream[T] {
	return Stream[T]{
		source: source,
	}
}

// StreamDrain drains the given channel.
func StreamDrain[T any](ch <-chan T) {
	for range ch {
	}
}

// Done waits all upstreaming operations to be done.
func (s Stream[T]) Done() {
	StreamDrain(s.source)
}
