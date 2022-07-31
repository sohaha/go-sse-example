package main

import (
	"time"

	"github.com/sohaha/zlsgo/zhttp"
	"github.com/sohaha/zlsgo/zlog"
)

func main() {
	zlog.ResetFlags(zlog.BitLevel | zlog.BitTime)

	sse := zhttp.SSE("http://127.0.0.1:3788/sse")

	go func() {
		time.Sleep(time.Second * 15)
		zlog.Warn("Manual close")
		sse.Close()
	}()

sseFor:
	for {
		select {
		case <-sse.Done():
			break sseFor
		case ev := <-sse.Event():
			zlog.Tipsf("id:%s msg:%s [%s]\n", ev.ID, string(ev.Data), ev.Event)
		}
	}

	zlog.Success("Done")
}
