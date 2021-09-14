package main

import "fmt"

type Node struct {
    value string
    next *Node
    prev *Node
}

var m = map[string]*Node{}
var head *Node
var tail *Node
const max = 5

func add(value string) {
    if val, ok := m[value]; ok {

        if val == head {
            return
        }

        val.prev.next = val.next
        if val.next != nil {
            val.next.prev = val.prev
        } else {
            tail = val.prev
        }

        val.next = head
        val.prev = nil
        head = val

    } else {
        node := new(Node)
        node.value = value
        node.next = head
        if head != nil {
            head.prev = node
        }
        if tail == nil {
            tail = node
        }
	    head = node
        m[value] = head

        if len(m) > max {
            tail.prev.next = nil
            tail = tail.prev
            delete(m, value)
        }
    }
}

func printCache() {
    curr := head
    for curr != nil {
        fmt.Printf("Node %p, value %v, next: %p, prev: %p\n", curr, curr.value, curr.next, curr.prev)
        curr = curr.next
    }
    fmt.Printf("\nHead %p, value %v, next: %p, prev: %p\n", head, head.value, head.next, head.prev)
    fmt.Printf("Tail %p, value %v, next: %p, prev: %p\n", tail, tail.value, tail.next, tail.prev)

}

func main() {
    add("dog")
    add("cat")
    add("elephant")
    add("hohipopotamus")
    add("giraffe")
    add("giraffe")

    printCache()
}