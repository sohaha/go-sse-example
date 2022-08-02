// main.go
package main

import (
	"time"

	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/znet/cors"
	"github.com/sohaha/zlsgo/zpprof"
	"github.com/sohaha/zlsgo/ztime"
	"github.com/sohaha/zlsgo/ztype"
)

func main() {
	// 获取一个实例
	r := znet.New()

	// 设置为开发模式
	r.SetMode(znet.DebugMode)

	// 开启 pprof
	zpprof.Register(r, "")

	r.Log.ResetFlags(zlog.BitLevel | zlog.BitTime)

	// 异常处理
	r.Use(znet.Recovery(func(c *znet.Context, err error) {
		e := err.Error()
		c.String(500, e)
	}))

	// 支持跨域
	r.Use(cors.Default())

	// 注册路由
	r.GET("/sse", func(c *znet.Context) {
		noretry := c.DefaultQuery("noretry", "")
		id := 0
		s := znet.NewSSE(c, func(lastID string, opts *znet.SSEOption) {
			opts.RetryTime = 1000
			if lastID != "" {
				id = ztype.ToInt(lastID)
				c.Log.Tips("This is even the request again:", lastID, noretry)
			}
		})

		go func() {
			i := 1
			_ = s.Send(ztype.ToString(id+i), "Hi 😋", "system")

			for {
				select {
				case <-s.Done():
					return
				default:
					i++

					// 定时发送当前时间
					now := ztime.Now()
					_ = s.Send(ztype.ToString(id+i), now)

					// N 次之后取消
					if i >= 5 {
						s.Stop()
					}

					time.Sleep(time.Second * 1)
				}
			}
		}()

		s.Push()
	})

	r.GET("/", func(c *znet.Context) {
		c.HTML(200, `<html><meta name="viewport" content="width=device-width, initial-scale=1.0"><title>SSE</title><body><ul id="app"></ul><script>var source = new EventSource("/sse?key=123"); source.addEventListener("message",function(event) {var li=document.createElement("li");li.innerText=event.data;document.querySelector("#app").appendChild(li);console.log("message",event.data)});source.addEventListener("open",function(event) {console.log("连接已经建立",event)});source.addEventListener("error",function(event,e) {console.log("error",event,e)});</script></body></html>`)
	})

	// 启动
	znet.Run()
}
