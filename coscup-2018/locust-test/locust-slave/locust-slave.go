package main

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/myzhan/boomer"
)

func foo() {

	start := boomer.Now()
	time.Sleep(100 * time.Millisecond)
	elapsed := boomer.Now() - start

	// Report your test result as a success, if you write it in python, it will looks like this
	// events.request_success.fire(request_type="http", name="foo", response_time=100, response_length=10)
	boomer.Events.Publish("request_success", "http", "foo", elapsed, int64(10))
}

func bar() {

	start := boomer.Now()
	time.Sleep(100 * time.Millisecond)
	elapsed := boomer.Now() - start

	// Report your test result as a failure, if you write it in python, it will looks like this
	// events.request_failure.fire(request_type="udp", name="bar", response_time=100, exception=Exception("udp error"))
	boomer.Events.Publish("request_failure", "udp", "bar", elapsed, "udp error")
}

func exampleGet() {
	start := boomer.Now()
	resp, err := http.Get("https://dl.google.com/go/go1.10.3.src.tar.gz")
	//resp, err := http.Get("https://cache.ruby-lang.org/pub/ruby/2.5/ruby-2.5.1.tar.gz")
	if err == nil {
		func() {
			defer resp.Body.Close()
			_, err = ioutil.ReadAll(resp.Body)
		}()
	}
	elapsed := boomer.Now() - start

	if err == nil {
		boomer.Events.Publish("request_success", "http", "example_get", elapsed, resp.ContentLength)
	} else {
		boomer.Events.Publish("request_failure", "http", "example_get", elapsed, err.Error())
	}
}

func main() {

	task1 := &boomer.Task{
		Name:   "foo",
		Weight: 10,
		Fn:     foo,
	}

	task2 := &boomer.Task{
		Name:   "bar",
		Weight: 20,
		Fn:     bar,
	}

	task3 := &boomer.Task{
		Name:   "example_get",
		Weight: 30,
		Fn:     exampleGet,
	}

	boomer.Run(task1, task2, task3)

}
