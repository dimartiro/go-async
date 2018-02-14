package main

import (
	"fmt"
	"reflect"
	"sync"
	"time"
)

func main() {
	a := func() string {
		time.Sleep(1 * time.Second)
		return "a"
	}

	b := func() string {
		time.Sleep(3 * time.Second)
		return "b"
	}

	c := func() string {
		time.Sleep(2 * time.Second)
		return "c"
	}

	calcTime(func() {
		any([]interface{}{a, b, c}, func(v interface{}) {
			fmt.Println(v)
		})
	})

	calcTime(func() {
		all([]interface{}{a, b, c}, func(v []interface{}) {
			for i, val := range v {
				fmt.Printf("%d - %s \n", i, val)
			}
		})
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

func any(actions []interface{}, callback func(v interface{})) {
	resp := make(chan interface{}, len(actions))

	for _, act := range actions {
		go func(act interface{}, resp chan<- interface{}) {
			val := reflect.ValueOf(act)
			resp <- val.Call(nil)
		}(act, resp)
	}

	callback(<-resp)
}

func all(actions []interface{}, callback func(v []interface{})) {
	resp := make(chan interface{}, len(actions))

	for _, act := range actions {
		go func(act interface{}, resp chan<- interface{}) {
			val := reflect.ValueOf(act)
			resp <- val.Call(nil)
		}(act, resp)
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
