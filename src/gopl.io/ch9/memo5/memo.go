// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 278.

// Package memo provides a concurrency-safe non-blocking memoization
// of a function.  Requests for different keys proceed in parallel.
// Concurrent requests for the same key block until the first completes.
// This implementation uses a monitor goroutine.
package memo

//!+Func

// Func is the type of the function to memoize.
type Func func(key string) (interface{}, error)

// A result is the result of calling a Func.
type result struct {
	value interface{}
	err   error
}

type entry struct {
	res   result
	ready chan struct{} // closed when res is ready
}

//!-Func

//!+get

// A request is a message requesting that the Func be applied to key.
type request struct {
	key      string
	response chan<- result // the client wants a single result
}

type Memo struct{ requests chan request }

// New returns a memoization of f.  Clients must subsequently call Close.
func New(f Func) *Memo {
	memo := &Memo{requests: make(chan request)}
	go memo.server(f)
	return memo
}

func (memo *Memo) Get(key string) (interface{}, error) {
	response := make(chan result)
	memo.requests <- request{key, response} //Get 单独获取一个通信频道
	res := <-response
	return res.value, res.err
}

func (memo *Memo) Close() { close(memo.requests) }

//!-get

//!+monitor

/*
首先利用了channel close的广播作用
其次把竞争数据放在同一个monitor goroutine中,实在是棒啊,这个例子很优秀
*/

func (memo *Memo) server(f Func) {
	cache := make(map[string]*entry)
	for req := range memo.requests { //每次获取request 通信
		e := cache[req.key]
		if e == nil {
			// This is the first request for this key.
			e = &entry{ready: make(chan struct{})}
			cache[req.key] = e
			go e.call(f, req.key) // call f(key)
		}
		go e.deliver(req.response) //妙啊,在monitor goroutine中创建一个goroutine 初始化,在创建一个进程返回值,把之前的思路叠加起来
	}
}

func (e *entry) call(f Func, key string) {
	// Evaluate the function.
	e.res.value, e.res.err = f(key)
	// Broadcast the ready condition.
	close(e.ready) //慢函数处理完之后关闭channel作为广播 数据以结束的标志
}

func (e *entry) deliver(response chan<- result) {
	// Wait for the ready condition.
	<-e.ready //非双参数返回格式,一直阻塞到 Close() 表示 数据初始化成功
	// Send the result to the client.
	response <- e.res //不过显然这样处理忽视了一种情况,那就是 那个天命之子获取request结果失败了,虽然他最终释放了ready广播信号,可是等待的各位兄弟们还是空,不过好像害的也就是那几个并发等待的
	//倒霉蛋而已
}

//!-monitor
