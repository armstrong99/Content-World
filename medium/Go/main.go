// Title: Implementing the Or-Pattern w/ Dynamic Select Block
// Developer: Ndukwe C. Armstrong
// Lang: Golang
// =========================
package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

// ============ UTILITY ============================================
func Get_List_Of_RecChans(n uint64) (c []<-chan int) {
	c = make([]<-chan int, 0, n)
	for i := range n {
		t := time.Duration(i+5) * time.Second // the least time is 5 seconds
		c = append(c, sig(t))
	}
	return c
}
func sig(t time.Duration) <-chan int {
	c := make(chan int)

	go func() {
		time.Sleep(t)
		c <- rand.Intn(10)
	}()

	return c
}

// ============== OR-Pattern w/ Dynamic Select Case ========
func Or(c ...<-chan int) <-chan int {
	done := make(chan int, 1)
	cases := make([]reflect.SelectCase, len(c))

	for i, ch := range c {
		cases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ch),
		}
	}

	go func() {
		chosenIndex, value, ok := reflect.Select(cases)

		// remove from case
		cases = append(cases[:chosenIndex], cases[chosenIndex+1:]...)

		if ok {
			done <- int(value.Int())
			return
		}
	}()

	return done
}

func main() {
	start := time.Now()

	channels := Get_List_Of_RecChans(5)

	report := <-Or(channels...)

	fmt.Printf("The total time spent is: %v\n", time.Since(start))
	fmt.Printf("And we go a value of: %v\n", report)
}

// package main

// import (
// 	"fmt"
// 	"time"
// )

// // or-channel pattern
// func or(channels ...<-chan interface{}) <-chan interface{} {
// 	switch len(channels) {
// 	case 0:
// 		return nil
// 	case 1:
// 		return channels[0]
// 	}

// 	orDone := make(chan interface{})
// 	go func() {
// 		defer close(orDone)
// 		switch len(channels) {
// 		case 2:
// 			select {
// 			case <-channels[0]:
// 			case <-channels[1]:
// 			}
// 		default:
// 			select {
// 			case <-channels[0]:
// 			case <-channels[1]:
// 			case <-channels[2]:
// 			case <-or(append(channels[3:], orDone)...): // recursion
// 			}
// 		}
// 	}()
// 	return orDone
// }

// func sig(after time.Duration) <-chan interface{} {
// 	c := make(chan interface{})
// 	go func() {
// 		defer close(c)
// 		time.Sleep(after)
// 	}()
// 	return c
// }

// func main() {
// 	start := time.Now()

// 	<-or(
// 		sig(2*time.Hour),
// 		sig(5*time.Minute),
// 		sig(1*time.Second),
// 		sig(1*time.Hour),
// 	)

// 	fmt.Printf("done after %v\n", time.Since(start))
// }
