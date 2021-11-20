package main

import (
	"fmt"
	"runtime"
)

var (
	const_pool_size   = 100000000
	const_routine_num = runtime.GOMAXPROCS(8)
)

func main() {
	done := make(chan struct{})
	var chs []chan int
	for i := 0; i < const_routine_num; i++ {
		chs = append(chs, make(chan int))
	}

	g := func(start int) {
		for i := start; ; i += const_routine_num {
			<-chs[start]
			fmt.Printf("协程 (%d): %d\n", start, i)
			chs[(start+1)%const_routine_num] <- i
			if i >= const_pool_size {
				done <- struct{}{}
			}
		}
	}

	for i := 0; i < const_routine_num; i++ {
		go g(i)
	}

	// 开始执行业务逻辑
	chs[0] <- 0

	// 等待结束
	<-done
}
