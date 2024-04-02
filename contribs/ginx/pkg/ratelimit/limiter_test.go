package ratelimit

import (
	"github.com/dstgo/maxwell/contribs/ginx/pkg/ratelimit/counter"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestCounterLimiter(t *testing.T) {
	engine := gin.Default()
	limit := 10
	window := time.Minute
	url := "http://localhost:8080/"

	limiter := counter.Limiter(
		counter.WithCounter(counter.Cache()),
		counter.WithLimit(limit),
		counter.WithWindow(window),
	)

	engine.Use(
		RateLimit(
			WithLimiter(limiter),
		),
	)

	done := make(chan struct{})

	go func() {
		time.Sleep(time.Second)
		for i := range limit - 1 {
			resp, err := http.Get(url)
			assert.Nil(t, err)
			t.Log(i, resp.StatusCode, err)
		}

		resp, err := http.Get(url)
		assert.Nil(t, err)
		assert.EqualValues(t, resp.StatusCode, http.StatusTooManyRequests)

		done <- struct{}{}
	}()

	go func() {
		engine.Run(":8080")
	}()

	select {
	case <-done:
	}
}
