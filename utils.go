package asyncLogger

import "sync"

type linkedData struct {
	data interface{}
	next *linkedData
}

type LinkedList struct {
	count uint64
	blockOnEmpty bool
	cond *sync.Cond
	pHead *linkedData
	pTail *linkedData
}

func NewLinkedList(blockOnEmpty bool)  *LinkedList {
	return &LinkedList{
		count:0,
		blockOnEmpty:blockOnEmpty,
		cond:sync.NewCond(&sync.Mutex{}),
		pHead:nil,
		pTail:nil}
}

func (ll *LinkedList) Count() uint64  {
	return ll.count
}

func (ll *LinkedList) Add(dataToAdd interface{}) error {

	linkData := &linkedData{data:dataToAdd, next:nil}

	ll.cond.L.Lock()
	defer ll.cond.L.Unlock()

	if ll.pHead == nil{
		ll.pHead = linkData
		ll.pTail = linkData
	}else {
		ll.pTail.next = linkData
		ll.pTail = linkData
	}

	if ll.blockOnEmpty {
		ll.cond.Signal()
	}

	ll.count++
	return nil
}

func (ll *LinkedList) Take() interface{} {

	ll.cond.L.Lock()
	defer ll.cond.L.Unlock()

	for ll.pHead == nil {
		if ll.blockOnEmpty {
			ll.cond.Wait()
		}else {
			return nil
		}
	}

	tmp := ll.pHead
	ll.pHead = ll.pHead.next
	ll.count--

	return tmp.data

}


