package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hello world")
	ch := make(chan int)
	test(ch)
	time.Sleep(time.Second * 5)
	ch <- 3
	time.Sleep(time.Second * 5)
}

func test(ch chan int) {
	<-ch
	fmt.Println("I am alive")
}
