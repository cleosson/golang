package main

import "fmt"

type FreqListNode struct {
    count  int
    first  *ItemNode
    last   *ItemNode
    prev   *FreqListNode
    next   *FreqListNode
}

type ItemNode struct {
    key      int
    value    string
    freqList *FreqListNode
    prev     *ItemNode
    next     *ItemNode
}

type LFU struct {
    mapItems     map[int]*ItemNode
    rootFreqList *FreqListNode
    capacity     int
}

func (node *FreqListNode) removeFreqListNode() {
    fmt.Printf("FreqListNode.removeFreqListNode - node: %p\n", node)

    prevNode := node.prev
    nextNode := node.next
    if node.next != nil {
        node.next.prev = prevNode
    }
    if node.prev != nil {
        node.prev.next = nextNode
    }
    node.next = nil
    node.prev = nil
}

func (node *FreqListNode) addFreqListNode(prevNode, nextNode *FreqListNode) {
    fmt.Printf("FreqListNode.addFreqListNode - node: %p, prevNode: %p, nextNode: %p\n", node, prevNode, nextNode)

    node.next = nextNode
    node.prev = prevNode

    if nextNode != nil {
        nextNode.prev = node
    }
    if prevNode != nil {
        prevNode.next = node
    }
}

func (node *ItemNode) removeItemNode() {
    fmt.Printf("ItemNode.removeItemNode - node: %p,\n", node)
    prevNode := node.prev
    nextNode := node.next
    if node.next != nil {
        node.next.prev = prevNode
    }
    if node.prev != nil {
        node.prev.next = nextNode
    }
    node.next = nil
    node.prev = nil
}

func (node *ItemNode) addItemNode(prevNode, nextNode *ItemNode) {
    fmt.Printf("ItemNode.addItemNode - node: %p, prevNode: %p, nextNode: %p\n", node, prevNode, nextNode)
    node.next = nextNode
    node.prev = prevNode

    if nextNode != nil {
        nextNode.prev = node
    }
    if prevNode != nil {
        prevNode.next = node
    }
}

func (lfu *LFU) getFreqList(count int) (*FreqListNode) {
    fmt.Printf("LFU.getFreqList - count: %+v, LFU: %+v\n", count, lfu)
    currFreqListNode := lfu.rootFreqList
    var prevFreqListNode *FreqListNode = nil
    for currFreqListNode != nil && currFreqListNode.count < count {
        prevFreqListNode = currFreqListNode
        currFreqListNode = currFreqListNode.next
    }

    if currFreqListNode != nil && currFreqListNode.count == count {
        fmt.Printf("LFU.getFreqList - LFU: %+v, currFreqListNode: %p, %+v\n", lfu, currFreqListNode, currFreqListNode)

        return currFreqListNode
    }

    newFreqListNode := new(FreqListNode)
    newFreqListNode.count = count
    newFreqListNode.addFreqListNode(prevFreqListNode, currFreqListNode)
    if lfu.rootFreqList == nil {
        lfu.rootFreqList = newFreqListNode
    }
    if prevFreqListNode == nil {
        lfu.rootFreqList = newFreqListNode
    }

    fmt.Printf("LFU.getFreqList - LFU: %+v, newFreqListNode: %p, %+v\n", lfu, newFreqListNode, newFreqListNode)

    return newFreqListNode
}

func (flNode *FreqListNode) addItemNode(itemNode *ItemNode) {
    fmt.Printf("FreqListNode.addItemNode - FreqList: %v, itemNode: %p, %+v\n", flNode, itemNode, itemNode)
    itemNode.addItemNode(flNode.last, nil)
    flNode.last = itemNode
    if flNode.first == nil {
        flNode.first = itemNode
    }
}

func (flNode *FreqListNode) removeItemNode(itemNode *ItemNode) {
    fmt.Printf("FreqListNode.removeItemNode - FreqList: %+v, itemNode: %p, %+v\n", flNode, itemNode, itemNode)
    if flNode.last == itemNode {
        flNode.last = itemNode.prev
    }
    if flNode.first == itemNode {
        flNode.first = itemNode.next
    }

    itemNode.removeItemNode()
    fmt.Printf("FreqListNode.removeItemNode - FreqList: %+v\n", flNode)
}

func (lfu *LFU) get(key int) string {
    fmt.Printf("LFU.get - key: %+v\n", key)
    value := ""
    if itemNode, ok := lfu.mapItems[key]; ok {
        oldFreqList := itemNode.freqList
        oldFreqList.removeItemNode(itemNode)
        newFreqList := lfu.getFreqList(itemNode.freqList.count + 1)
        newFreqList.addItemNode(itemNode)
        itemNode.freqList = newFreqList
        if oldFreqList.first == nil {
            if lfu.rootFreqList == oldFreqList {
                lfu.rootFreqList = oldFreqList.next
            }
            oldFreqList.removeFreqListNode()
        }
        fmt.Printf("LFU.get - itemNode: %p, %+v\n", itemNode, itemNode)
        value = itemNode.value
    }

    return value
}

func (lfu *LFU) put(key int, value string) {
    fmt.Printf("LFU.put - key: %+v, value: %+v, len(lfu.mapItems): %v\n", key, value, len(lfu.mapItems))
    if len(lfu.mapItems) < lfu.capacity {
        if itemNode, ok := lfu.mapItems[key]; !ok {
            itemNode = new(ItemNode)
            itemNode.key = key
            itemNode.value = value
            freqList := lfu.getFreqList(1)
            freqList.addItemNode(itemNode)
            itemNode.freqList = freqList
            lfu.mapItems[key] = itemNode

            fmt.Printf("LFU.put - itemNode: %p, %+v, len(lfu.mapItems): %v\n", itemNode, itemNode, len(lfu.mapItems))
        }
    } else {
        itemNode := new(ItemNode)
        itemNode.key = key
        itemNode.value = value
        toBeRemoved := lfu.rootFreqList.last
        lfu.rootFreqList.removeItemNode(toBeRemoved)
        lfu.rootFreqList.addItemNode(itemNode)
        itemNode.freqList = lfu.rootFreqList
        delete(lfu.mapItems, toBeRemoved.key)
        lfu.mapItems[key] = itemNode

        fmt.Printf("LFU.put - itemNode: %p, %+v, len(lfu.mapItems): %v\n", itemNode, itemNode, len(lfu.mapItems))
    }
}

func (lfu *LFU) print() {
    fmt.Println("####")
    fmt.Printf("lfu %+v\n", lfu)
    fmt.Println("#")
    for key, value := range lfu.mapItems {
        fmt.Printf("lfu.mapItems - key: %+v, value: %p, %v+\n", key, value, value)
    }
    fmt.Println("#")
    currFreq := lfu.rootFreqList
    for currFreq != nil {
        fmt.Printf("currFreq: %p, count: %+v, first: %p, last: %p, prev: %p, next: %p\n", currFreq, currFreq.count, currFreq.first, currFreq.last, currFreq.prev, currFreq.next)
        currFreq = currFreq.next
    }
    fmt.Println("#")
    fmt.Printf("lfu.capacity: %+v\n", lfu.capacity)
}

func main() {
    test := LFU{}
    test.mapItems = make(map[int]*ItemNode)
    test.capacity = 4
    test.rootFreqList = nil

    fmt.Println("##################")
    fmt.Println("test.put(1, dog)")
    test.put(1, "dog")
    test.print()
    fmt.Println("##################")
    fmt.Println("test.put(2, cat)")
    test.put(2, "cat")
    test.print()
    fmt.Println("##################")
    fmt.Printf("test.get(1) %+v\n", test.get(1))
    test.print()
    fmt.Println("##################")
    fmt.Printf("test.get(2) %+v\n", test.get(2))
    test.print()
    fmt.Println("##################")
    fmt.Println("test.put(3, elephant)")
    test.put(3, "elephant")
    test.print()
    fmt.Println("##################")
    fmt.Println("test.put(4, giraffe)")
    test.put(4, "giraffe")
    test.print()
    fmt.Println("##################")
    fmt.Println("test.put(5, hipopotamus)")
    test.put(5, "hipopotamus")
    test.print()
    fmt.Println("##################")
    fmt.Printf("test.get(5) %+v\n", test.get(5))
    test.print()
    fmt.Println("##################")
    fmt.Printf("test.get(5) %+v\n", test.get(5))
    test.print()
    fmt.Println("##################")
    fmt.Printf("test.get(5) %+v\n", test.get(5))
    test.print()
    fmt.Println("##################")
    fmt.Printf("test.get(5) %+v\n", test.get(5))
    test.print()
    fmt.Println("##################")
    fmt.Printf("test.get(5) %+v\n", test.get(5))
    test.print()
    fmt.Println("##################")
    fmt.Printf("test.get(5) %+v\n", test.get(5))
    test.print()
    fmt.Println("##################")
    fmt.Printf("test.get(5) %+v\n", test.get(5))
    test.print()
    fmt.Println("##################")
    fmt.Printf("test.get(5) %+v\n", test.get(5))
    test.print()
    fmt.Println("##################")
    fmt.Printf("test.get(5) %+v\n", test.get(5))
    test.print()
}