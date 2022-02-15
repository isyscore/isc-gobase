package isc

import "sort"

type T any
type C comparable

// A Stream is a stream that can be used to do stream processing.
type Stream struct {
	source <-chan any
}

type ComparableStream struct {
	source <-chan C
}

// Just converts the given arbitrary items to a Stream.
func Just(items ...T) Stream {
	source := make(chan any, len(items))
	for _, e := range items {
		source <- e
	}
	close(source)
	return Range(source)
}

// AllMatch returns whether all elements of this stream match the provided predicate.
// May not evaluate the predicate on all elements if not necessary for determining the result.
// If the stream is empty then true is returned and the predicate is not evaluated.
func (s Stream) AllMatch(predicate func(T) bool) bool {
	for item := range s.source {
		if !predicate(item) {
			// make sure the former goroutine not block, and current func returns fast.
			go drain(s.source)
			return false
		}
	}
	return true
}

// AnyMatch returns whether any elements of this stream match the provided predicate.
// May not evaluate the predicate on all elements if not necessary for determining the result.
// If the stream is empty then false is returned and the predicate is not evaluated.
func (s Stream) AnyMatch(predicate func(T) bool) bool {
	for item := range s.source {
		if predicate(item) {
			// make sure the former goroutine not block, and current func returns fast.
			go drain(s.source)
			return true
		}
	}
	return false
}

// NoneMatch returns whether all elements of this stream don't match the provided predicate.
// May not evaluate the predicate on all elements if not necessary for determining the result.
// If the stream is empty then true is returned and the predicate is not evaluated.
func (s Stream) NoneMatch(predicate func(T) bool) bool {
	for item := range s.source {
		if predicate(item) {
			go drain(s.source)
			return false
		}
	}
	return true
}

func (s Stream) Map(fn func(T) any) Stream {
	source := make(chan any)
	for item := range s.source {
		source <- fn(item)
	}
	return Range(source)
}

// Sort sorts the items from the underlying source.
func (s Stream) Sort(less func(T, T) bool) Stream {
	var items []T
	for item := range s.source {
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		return less(items[i], items[j])
	})
	return Just(items)
}

// Distinct removes the duplicated items base on the given keySelector.
func (s Stream) Distinct(keySelector func(T) C) Stream {
	source := make(chan any)
	m := make(map[C]T)
	for item := range s.source {
		m[keySelector(item)] = item
	}
	for _, v := range m {
		source <- v
	}
	return Range(source)
}

// Foreach seals the Stream with the fn on each item, no successive operations.
func (s Stream) Foreach(fn func(T)) {
	for item := range s.source {
		fn(item)
	}
}

//First returns the first element,channel is FIFO,so first goroutine will get head element
func (s Stream) First(valueSelector func(any) bool) Stream {
	source := make(chan any)
	go func() {
		for item := range s.source {
			if valueSelector(item) {
				source <- item
				break
			}
		}
	}()
	close(source)
	return Range(source)
}

// Last returns the last item, or nil if no items.
func (s Stream) Last(valueSelector func(any) bool) Stream {
	source := make(chan any)
	var lastValue any
	for lastValue = range s.source {
	}
	source <- lastValue
	close(source)
	return Range(source)
}

func New(source <-chan any) Stream {
	return Stream{
		source: source,
	}
}
func Range(source <-chan any) Stream {
	return Stream{
		source: source,
	}
}

// drain drains the given channel.
func drain(ch <-chan any) {
	for range ch {
	}
}

// Done waits all upstreaming operations to be done.
func (s Stream) Done() {
	drain(s.source)
}
