// main.go
package main

import (
	"time"

	"github.com/sohaha/zlsgo/zlog"
	"github.com/sohaha/zlsgo/znet"
	"github.com/sohaha/zlsgo/znet/cors"
	"github.com/sohaha/zlsgo/ztime"
	"github.com/sohaha/zlsgo/ztype"
)

func main() {
	// 获取一个实例
	r := znet.New()

	// 设置为开发模式
	r.SetMode(znet.DebugMode)

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
		s := znet.NewSSE(c, func(lastID string, opts *znet.SSEOption) {
			opts.RetryTime = 5000
			if lastID != "" {
				zlog.Tips("This is even the request again:", lastID, noretry)
				if noretry != "" {
					opts.Verify = false
				}
			}
		})

		id := 0

		s.Send(ztype.ToString(id), "Hi 😋", "system")

		for {
			select {
			case <-s.Done():
				zlog.Warn("disconnects")
				return
			default:
				id++
				// 定时发送当前时间
				now := ztime.Now()
				s.Send(ztype.ToString(id), now)

				time.Sleep(time.Second * 1)

				// 5 次之后取消
				if id >= 5 {
					s.Send("-1", "This is the last one", "system")
					return
				}
			}
		}
	})

	r.GET("/", func(c *znet.Context) {
		c.String(200, "Hello world")
	})

	// 启动
	znet.Run()
}
