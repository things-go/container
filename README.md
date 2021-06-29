# container

container implements containers, currently the containers are not thread-safe.

[![GoDoc](https://godoc.org/github.com/things-go/container?status.svg)](https://godoc.org/github.com/things-go/container)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/things-go/container?tab=doc)
[![Build Status](https://www.travis-ci.com/things-go/container.svg?branch=main)](https://www.travis-ci.com/things-go/container)
[![codecov](https://codecov.io/gh/things-go/container/branch/main/graph/badge.svg)](https://codecov.io/gh/things-go/container)
![Action Status](https://github.com/things-go/container/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/things-go/container)](https://goreportcard.com/report/github.com/things-go/container)
[![License](https://img.shields.io/github/license/things-go/container)](https://github.com/things-go/container/raw/master/LICENSE)
[![Tag](https://img.shields.io/github/v/tag/things-go/container)](https://github.com/things-go/container/tags)

- **[How to use this repo](#how-to-use-this-package)**
- **[Containers](#Containers-Interface)**
    - [Stack](#stack)
        - stack use container/list.
        - quick stack use builtin slice.
    - [Queue](#queue)
        - queue use container/list
        - quick queue use builtin slice.
    - [PriorityQueue](#priorityqueue) use builtin slice with container/heap
    - [ArrayList](#arraylist) use builtin slice.
    - [LinkedList](#linkedlist) use container/list
    - [LinkedMap](#linkedMap) use container/list and builtin map.
- **[safe container](#safe-container)**
    - [fifo](#fifo) FIFO is a thread-safe Queue. in which (a) each accumulator is simply the most 
      recently provided object and (b) the collection of keys to process is a FIFO.
      > FIFO solves this use case:
      > * You want to process every object (exactly) once.
      > * You want to process the most recent version of the object when you process it.
      > * You do not want to process deleted objects, they should be removed from the queue.
      > * You do not want to periodically reprocess objects.
    - [heap](#heap) Heap is a thread-safe producer/consumer queue that implements a heap data structure.It can be used to implement priority queues and similar data structures.
- **[others](#others)**
    - [Comparator](#Comparator)
        - [Sort](#sort) sort with Comparator interface
        - [Heap](#heap) heap with Comparator interface
    
## Donation

if package help you a lot,you can support us by:

**Alipay**

![alipay](https://github.com/thinkgos/thinkgos/blob/master/asserts/alipay.jpg)

**WeChat Pay**

![wxpay](https://github.com/thinkgos/thinkgos/blob/master/asserts/wxpay.jpg)