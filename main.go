package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func main() {
	r := gin.Default()
	r.Use(TraceMiddleware())
	r.Any("/static", static)
	r.Any("/uniform", uniform)
	r.Any("/lognorm", lognorm)
	r.Run()
}

const sampleJSON = "{\n  \"first\": 0.12345,\n  \"second\": 0.12345,\n  \"third\": 0.12345,\n  \"forth\": 0.12345,\n  \"fifth\": 0.12345,\n  \"sixth\": 0.12345\n}"
const defaultStatic = 25
const defaultUniformFrom = 10
const defaultUniformTo = 30
const defaultLognormMean = 25
const defaultLognormStdev = 15

func static(c *gin.Context) {
	params := struct {
		Time  float64 `form:"time"`
		Cores int     `form:"cores"`
	}{
		Time:  defaultStatic,
		Cores: 0,
	}
	c.ShouldBindQuery(&params)
	sleepAndRespond(c, time.Duration(params.Time)*time.Millisecond, params.Cores)
}

func uniform(c *gin.Context) {
	params := struct {
		From  int64 `form:"from"`
		To    int64 `form:"to"`
		Cores int   `form:"cores"`
	}{
		From:  defaultUniformFrom,
		To:    defaultUniformTo,
		Cores: 0,
	}
	c.BindQuery(&params)
	delay := uniformRand(params.From, params.To)
	sleepAndRespond(c, time.Millisecond*time.Duration(delay), params.Cores)
}

func lognorm(c *gin.Context) {
	params := struct {
		Mean  float64 `form:"mean"`
		Stdev float64 `form:"stdev"`
		Cores int     `form:"cores"`
	}{
		Mean:  defaultLognormMean,
		Stdev: defaultLognormStdev,
		Cores: 0,
	}
	c.BindQuery(&params)
	delay := lognormalRand(params.Mean, params.Stdev)
	sleepAndRespond(c, time.Millisecond*time.Duration(delay), params.Cores)
}

func busySleep(done chan int) {
	for {
		select {
		case <-done:
			return
		default:
		}
	}
}

func sleepAndRespond(c *gin.Context, sleepTime time.Duration, cores int) {
	done := make(chan int)
	for i := 0; i < cores; i++ {
		go busySleep(done)
	}
	time.Sleep(sleepTime)
	close(done)
	end := time.Now()
	c.Writer.Header().Set("X-End", fmt.Sprintf("%.7f", float64(end.UnixNano())/float64(1000000000)))
	c.JSON(http.StatusOK, gin.H{
		"first":  1234.1234,
		"second": 1234.1234,
		"third":  1234.1234,
		"fourth": 1234.1234,
		"fifth":  1234.1234,
	})
}
