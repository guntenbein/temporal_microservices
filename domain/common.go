package domain

import (
	"context"
	"reflect"
	"runtime"
	"strings"
	"time"

	"go.temporal.io/sdk/activity"
)

func GetActivityName(i interface{}) string {
	fullName := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	parts := strings.Split(fullName, ".")
	lastPart := parts[len(parts)-1]
	parts2 := strings.Split(lastPart, "-")
	firstPart := parts2[0]
	return firstPart
}

type SimpleHeartbit struct {
	cancelc chan struct{}
}

func StartHeartbeat(ctx context.Context, period time.Duration) (h *SimpleHeartbit) {
	h = &SimpleHeartbit{make(chan struct{})}
	go func() {
		for {
			select {
			case <-h.cancelc:
				return
			default:
				{
					activity.RecordHeartbeat(ctx)
					time.Sleep(period)
				}
			}
		}
	}()
	return h
}

func (h *SimpleHeartbit) Stop() {
	close(h.cancelc)
}
