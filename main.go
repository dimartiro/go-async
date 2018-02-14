package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	a := func() interface{} {
		time.Sleep(1 * time.Second)
		return "a"
	}

	b := func() interface{} {
		time.Sleep(3 * time.Second)
		return "b"
	}

	c := func() interface{} {
		time.Sleep(2 * time.Second)
		return "c"
	}

	calcTime(func() {
		any(func(v interface{}) {
			fmt.Println(v)
		}, a, b, c)
	})

	calcTime(func() {
		all(func(v []interface{}) {
			for i, val := range v {
				fmt.Printf("%d - %s \n", i, val)
			}
		}, a, b, c)
	})

	calcTime(func() {
		eachParallel([]interface{}{1, 2, 3, 4}, func(t interface{}) {
			sleepTime, ok := t.(int)
			if ok {
				time.Sleep(time.Duration(sleepTime) * time.Second)
			}
		})
	})
}

func any(callback func(v interface{}), actions ...func() interface{}) {
	resp := make(chan interface{}, len(actions))

	for _, f := range actions {
		go func(f func() interface{}, resp chan<- interface{}) {
			resp <- f()
		}(f, resp)
	}

	callback(<-resp)
}

func all(callback func(v []interface{}), actions ...func() interface{}) {
	resp := make(chan interface{}, len(actions))

	for _, f := range actions {
		go func(f func() interface{}, resp chan<- interface{}) {
			resp <- f()
		}(f, resp)
	}

	finalResp := make([]interface{}, len(actions))
	for i := 0; i < len(actions); i++ {
		finalResp[i] = <-resp
	}

	callback(finalResp)
}

func eachParallel(elements []interface{}, function func(element interface{})) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(elements))

	defer waitGroup.Wait()

	for _, e := range elements {
		go func(e interface{}) {
			defer waitGroup.Done()
			function(e)
		}(e)
	}
}

func calcTime(c func()) {
	start := time.Now()
	c()
	elapsed := time.Since(start)
	fmt.Printf("Finish in %s \n", elapsed)
}
