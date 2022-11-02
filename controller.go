package main

import (
	"context"
	"fmt"
	"log"
	"my-web/framework"
	"time"
)

func FooControllerHandler(c *framework.Context) error {
	// 这个chan负责通知结束
	finish := make(chan struct{}, 1)
	// 这个chan负责通知panic异常
	panicChan := make(chan interface{}, 1)

	durationCtx, cancel := context.WithTimeout(c.BaseContext(), time.Duration(1*time.Second))
	// 记得当所有事情处理结束后调用cancel,告知durationCtx的后续Context结束
	defer cancel()

	go func() {
		// 异常处理
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()
		// 这里做具体业务
		time.Sleep(10 * time.Second)
		c.Json(200, "ok")
		// 新的goroutine结束时通过finish告知父goroutine
		finish <- struct{}{}
	}()

	// 请求监听的时候增加锁机制
	select {
	// 监听panic
	case p := <-panicChan:
		c.WriterMux().Lock()
		defer c.WriterMux().Unlock()
		c.Json(500, "panic")
		log.Println(p)
	// 监听结束事件
	case <-finish:
		fmt.Println("finish")
	// 监听超时事件
	case <-durationCtx.Done():
		c.WriterMux().Lock()
		defer c.WriterMux().Unlock()
		c.Json(500, "time out")
		c.SetHasTimeout()
	}

	return nil
}
