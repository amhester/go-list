package list

import (
	"log"
	"testing"
	"time"
)

func TestListPush(t *testing.T) {
	list := &List{}
	list.Push(int(1))
	if list.Length != 1 {
		t.Log("List length not equal to 1.")
		t.Fail()
	}
	if list.root.Value.(int) != 1 {
		t.Log("Root value not equal to 1.")
		t.Fail()
	}
	list.Push(int(2))
	if list.Length != 2 {
		t.Log("List length not equal to 2.")
		t.Fail()
	}
	if list.root.Next.Value.(int) != 2 {
		t.Log("Root.Next value not equal to 2.")
		t.Fail()
	}
}

func TestListPop(t *testing.T) {
	list := &List{}
	list.Push(int(1))
	list.Push(int(2))
	if list.Length != 2 {
		t.Log("List length not equal to 2.")
		t.Fail()
	}
	val2 := list.Pop().(int)
	if list.Length != 1 {
		t.Log("List length not equal to 1.")
		t.Fail()
	}
	if val2 != 2 {
		t.Log("Val2 value not equal to 2.")
		t.Fail()
	}
	val1 := list.Pop().(int)
	if list.Length != 0 {
		t.Log("List length not equal to 0.")
		t.Fail()
	}
	if val1 != 1 {
		t.Log("Val1 value not equal to 1.")
		t.Fail()
	}
	val0 := list.Pop()
	if val0 != nil {
		t.Log("Val0 value not equal to nil.")
		t.Fail()
	}
}

func TestListQueue(t *testing.T) {
	list := &List{}
	list.Queue(int(1))
	if list.Length != 1 {
		t.Log("List length not equal to 1.")
		t.Fail()
	}
	if list.root.Value.(int) != 1 {
		t.Log("Root value not equal to 1.")
		t.Fail()
	}
	list.Queue(int(2))
	if list.Length != 2 {
		t.Log("List length not equal to 2.")
		t.Fail()
	}
	if list.root.Value.(int) != 2 {
		t.Log("Root value not equal to 2.")
		t.Fail()
	}
	if list.root.Next.Value.(int) != 1 {
		t.Log("Root.Next value not equal to 1.")
		t.Fail()
	}
}

func TestListDequeue(t *testing.T) {
	list := &List{}
	list.Queue(int(1))
	list.Queue(int(2))
	if list.Length != 2 {
		t.Log("List length not equal to 2.")
		t.Fail()
	}
	val2 := list.Dequeue().(int)
	if list.Length != 1 {
		t.Log("List length not equal to 1.")
		t.Fail()
	}
	if val2 != 2 {
		t.Log("Val2 value not equal to 2.")
		t.Fail()
	}
	val1 := list.Dequeue().(int)
	if list.Length != 0 {
		t.Log("List length not equal to 0.")
		t.Fail()
	}
	if val1 != 1 {
		t.Log("Val1 value not equal to 1.")
		t.Fail()
	}
	val0 := list.Dequeue()
	if val0 != nil {
		t.Log("Val0 value not equal to nil.")
		t.Fail()
	}
}

func TestListEnumerate(t *testing.T) {
	list := &List{}
	list.Push(1)
	list.Push(2)
	list.Push(3)
	ls := list.Enumerate()
	if ls == nil {
		t.Log("Enumeration is nil.")
		t.Fail()
	}
	if len(ls) != 3 {
		t.Log("Enumeration yielded slice of incorrect length.")
		t.Fail()
	}
}

func TestListSlice(t *testing.T) {
	list := &List{}
	list.Push(1)
	list.Push(2)
	list.Push(3)
	ls := list.Slice(1, 3)
	if ls == nil {
		t.Log("Slice is nil.")
		t.Fail()
	}
	if ls.Length != 2 {
		t.Log("Slice yielded list of incorrect length.")
		t.Fail()
	}
}

func TestListForEach(t *testing.T) {
	list := &List{}
	list.Push(1)
	list.Push(2)
	list.Push(3)
	list.ForEach(func(val interface{}, idx int) {
		if val.(int) != idx+1 {
			t.Log("Incorrect value at position in list.", val, idx+1)
			t.Fail()
		}
	})
}

func TestListParallelForEach(t *testing.T) {
	list := &List{}
	list.Push(1)
	list.Push(2)
	list.Push(3)
	list.ParallelForEach(func(val interface{}, idx int) {
		time.Sleep(time.Duration(1/(idx+1)) * time.Second * 3)
		log.Println(val)
		if val.(int) != idx+1 {
			t.Log("Incorrect value at position in list.", val, idx+1)
			t.Fail()
		}
	})
}

func TestListMap(t *testing.T) {
	list := &List{}
	list.Push(1)
	list.Push(2)
	list.Push(3)
	ls := list.Map(func(val interface{}, idx int) interface{} {
		return val.(int) * val.(int)
	})
	if ls == nil {
		t.Log("Mapped list is nil.")
		t.Fail()
	}
	if ls.Length != 3 {
		t.Log("Length of mapped list is incorrect.")
		t.Fail()
	}
	for i, item := range ls.Enumerate() {
		if item.(int) != ((i + 1) * (i + 1)) {
			t.Log("Incorrect mapped value.")
			t.Fail()
		}
	}
}

func TestListParallelMap(t *testing.T) {
	list := &List{}
	list.Push(1)
	list.Push(2)
	list.Push(3)
	ls := list.ParallelMap(func(val interface{}, idx int) interface{} {
		return val.(int) * val.(int)
	})
	if ls == nil {
		t.Log("Mapped list is nil.")
		t.Fail()
	}
	if ls.Length != 3 {
		t.Log("Length of mapped list is incorrect.")
		t.Fail()
	}
	accum := 0
	for _, val := range ls.Enumerate() {
		accum += val.(int)
	}
	if accum != 14 {
		t.Log("Incorrect checksum for parallel map.")
		t.Fail()
	}
}

func TestListFilter(t *testing.T) {
	list := &List{}
	list.Push(1)
	list.Push(2)
	list.Push(3)
	ls := list.Filter(func(val interface{}, idx int) bool {
		return idx != 1
	})
	if ls == nil {
		t.Log("Filtered list is nil.")
		t.Fail()
	}
	if ls.Length != 2 {
		t.Log("Length of filtered list is incorrect.")
		t.Fail()
	}
	accum := 0
	for _, val := range ls.Enumerate() {
		accum += val.(int)
	}
	if accum != 4 {
		t.Log("Incorrect checksum for filtered list.")
		t.Fail()
	}
}

func TestListReduce(t *testing.T) {
	list := &List{}
	list.Push(1)
	list.Push(2)
	list.Push(3)
	final := list.Reduce(func(val interface{}, idx int, accum interface{}) interface{} {
		return accum.(int) + val.(int)
	}, 0)
	if final != 6 {
		t.Log("Incorrect reduce value.")
		t.Fail()
	}
}

func TestListClear(t *testing.T) {
	list := &List{}
	list.Push(1)
	list.Push(2)
	list.Push(3)
	list.Clear()
	if list.Length != 0 {
		t.Log("List length is not 0.")
		t.Fail()
	}
}
