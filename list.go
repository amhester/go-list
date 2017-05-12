package list

import (
	"runtime"
	"sync"
)

var defaultParallelism = int(runtime.NumCPU())

//GenericLink is a container stucture for each item in the list.
type GenericLink struct {
	Value    interface{}
	Next     *GenericLink
	Previous *GenericLink
}

//List is the actual double-link(ish) list data structure
type List struct {
	root   *GenericLink
	last   *GenericLink
	Length int
}

//Push puts val on to the end of the list.
func (list *List) Push(val interface{}) {
	node := &GenericLink{Value: val}
	if list.root == nil {
		list.root = node
		list.last = node
		list.Length = 1
		return
	}
	node.Previous = list.last
	list.last.Next = node
	list.last = node
	list.Length++
}

//Pop removes and returns the last item on the list.
func (list *List) Pop() interface{} {
	if list.Length == 0 {
		return nil
	}
	if list.Length == 1 {
		node := list.last
		list.last = nil
		list.root = nil
		list.Length = 0
		return node.Value
	}
	node := list.last
	list.last = list.last.Previous
	list.last.Next = nil
	list.Length--
	node.Previous = nil
	return node.Value
}

//Queue adds val to the beginning of the list.
func (list *List) Queue(val interface{}) {
	node := &GenericLink{Value: val}
	if list.root == nil {
		list.root = node
		list.last = node
		list.Length = 1
		return
	}
	node.Next = list.root
	list.root.Previous = node
	list.root = node
	list.Length++
}

//Dequeue removes an item from the beginning of the list.
func (list *List) Dequeue() interface{} {
	if list.Length == 0 {
		return nil
	}
	if list.Length == 1 {
		node := list.root
		list.last = nil
		list.root = nil
		list.Length = 0
		return node.Value
	}
	node := list.root
	list.root = list.root.Next
	list.root.Previous = nil
	list.Length--
	node.Next = nil
	return node.Value
}

//Enumerate returns a slice representation of the list.
func (list *List) Enumerate() []interface{} {
	if list.Length == 0 {
		return []interface{}{}
	}
	res := make([]interface{}, list.Length)
	i := 0
	node := list.root
	for node != nil {
		res[i] = node.Value
		node = node.Next
		i++
	}
	return res
}

//Slice returns a new *List starting at start (inclusive) going to end (exclusive).
func (list *List) Slice(start, end int) *List {
	if start < 0 {
		start = 0
	}
	res := &List{}
	i := 0
	node := list.root
	for node != nil && i < end {
		if i >= start {
			res.Push(node.Value)
		}
		node = node.Next
		i++
	}
	return res
}

//ForEach iterates over list and calls fn with the value and index of each iteration.
func (list *List) ForEach(fn func(val interface{}, idx int)) {
	i := 0
	node := list.root
	for node != nil {
		fn(node.Value, i)
		node = node.Next
		i++
	}
}

//ParallelForEach is the same as ForEach, but will run fn concurrently with goroutines, but will wait for all routines to stop.
func (list *List) ParallelForEach(fn func(val interface{}, idx int), maxParallelism int) {
	wg := &sync.WaitGroup{}
	if maxParallelism < 1 {
		maxParallelism = defaultParallelism
	}
	sema := make(chan bool, maxParallelism)
	for i := 0; i < maxParallelism; i++ {
		sema <- true
	}
	i := 0
	node := list.root
	for node != nil {
		<-sema
		wg.Add(1)
		go func(wg_ *sync.WaitGroup, sema_ chan bool, val_ interface{}, idx_ int) {
			fn(val_, idx_)
			wg_.Done()
			sema_ <- true
		}(wg, sema, node.Value, i)
		node = node.Next
		i++
	}
	wg.Wait()
	close(sema)
}

//Map iterates over list calling fn with the value and index of each iteration, pushing the result onto a new list which is returned.
func (list *List) Map(fn func(val interface{}, idx int) interface{}) *List {
	res := &List{}
	i := 0
	node := list.root
	for node != nil {
		val := fn(node.Value, i)
		res.Push(val)
		node = node.Next
		i++
	}
	return res
}

//ParallelMap is the same as Map, but will run fn concurrently with goroutines, but will wait for all routines to stop. Order of returned list is not gaurunteed.
func (list *List) ParallelMap(fn func(val interface{}, idx int) interface{}, maxParallelism int) *List {
	wg := &sync.WaitGroup{}
	if maxParallelism < 1 {
		maxParallelism = defaultParallelism
	}
	sema := make(chan bool, maxParallelism)
	for i := 0; i < maxParallelism; i++ {
		sema <- true
	}
	c := make(chan interface{}, list.Length)
	res := &List{}
	i := 0
	node := list.root
	for node != nil {
		<-sema
		wg.Add(1)
		go func(wg_ *sync.WaitGroup, c_ chan interface{}, sema_ chan bool, val_ interface{}, idx_ int) {
			val := fn(val_, idx_)
			c <- val
			wg_.Done()
			sema_ <- true
		}(wg, c, sema, node.Value, i)
		node = node.Next
		i++
	}
	wg.Wait()
	close(c)
	for val := range c {
		res.Push(val)
	}
	return res
}

//Filter iterates over list calling fn with the value and index of each iteration and will push the value onto a new list if the result of fn == true.
func (list *List) Filter(fn func(val interface{}, idx int) bool) *List {
	res := &List{}
	i := 0
	node := list.root
	for node != nil {
		val := fn(node.Value, i)
		if val {
			res.Push(node.Value)
		}
		node = node.Next
		i++
	}
	return res
}

//Reduce iterates over list calling fn with the value and index of each iteration along with the accumulation variable, assigning accum the result of fn. The value of accum after iteration is returned.
func (list *List) Reduce(fn func(val interface{}, idx int, accum interface{}) interface{}, init interface{}) interface{} {
	res := init
	i := 0
	node := list.root
	for node != nil {
		res = fn(node.Value, i, res)
		node = node.Next
		i++
	}
	return res
}

//First returns the first n items on the list.
func (list *List) First(n int) *List {
	res := &List{}
	i := 0
	node := list.root
	for node != nil && i < n {
		res.Push(node.Value)
		node = node.Next
		i++
	}
	return res
}

//Last returns the last n items on the list.
func (list *List) Last(n int) *List {
	res := &List{}
	i := 0
	node := list.last
	for node != nil && i < n {
		res.Push(node.Value)
		node = node.Previous
		i++
	}
	return res
}

//ValueAt returns the value at the specified index.
func (list *List) ValueAt(index int) interface{} {
	i := 0
	node := list.root
	for node != nil {
		if i == index {
			return node.Value
		}
		node = node.Next
		i++
	}
	return nil
}

//RemoveAt removes an item at the specified index. Returns a true if an item was removed, false if not.
func (list *List) RemoveAt(index int) bool {
	i := 0
	node := list.root
	for node != nil {
		if i == index {
			if node.Next != nil {
				node.Next.Previous = node.Previous
			} else {
				list.last = node.Previous
			}
			if node.Previous != nil {
				node.Previous.Next = node.Next
			} else {
				list.root = node.Next
			}
			node.Next = nil
			node.Previous = nil
			node = nil
			return true
		}
		node = node.Next
		i++
	}
	return false
}

//Clear empties the list of all items.
func (list *List) Clear() {
	node := list.last
	for node != nil {
		prev := node.Previous
		node.Next = nil
		node.Previous = nil
		node = prev
	}
	list.Length = 0
}
