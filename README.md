# container

container implements containers, currently the containers are not thread-safe.

[![GoDoc](https://godoc.org/github.com/things-go/container?status.svg)](https://godoc.org/github.com/things-go/container)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/things-go/container?tab=doc)
[![codecov](https://codecov.io/gh/things-go/container/branch/main/graph/badge.svg)](https://codecov.io/gh/things-go/container)
![Action Status](https://github.com/things-go/container/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/things-go/container)](https://goreportcard.com/report/github.com/things-go/container)
[![License](https://img.shields.io/github/license/things-go/container)](https://github.com/things-go/container/raw/master/LICENSE)
[![Tag](https://img.shields.io/github/v/tag/things-go/container)](https://github.com/things-go/container/tags)

- Containers
  - Stack
    - stack use go/list.
    - quick stack use builtin slice.
  - Queue
    - queue use go/list
    - quick queue use builtin slice.
    - priority queue
  - PriorityQueue use builtin slice with container/heap
  - ArrayList use builtin slice.
  - LinkedList use go/list
  - LinkedMap use go/list and builtin map.
- safe container
  - fifo FIFO is a thread-safe Queue. in which (a) each accumulator is simply the most
    recently provided object and (b) the collection of keys to process is a FIFO.
    > FIFO solves this use case:
    > - You want to process every object (exactly) once.
    > - You want to process the most recent version of the object when you process it.
    > - You do not want to process deleted objects, they should be removed from the queue.
    > - You do not want to periodically reprocess objects.
- others
  - Comparator sort and heap with Comparable
  - go
    - list
    - heap
    - ring
