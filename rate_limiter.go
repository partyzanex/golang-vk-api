package vkapi

import "time"

type rateLimiter struct {
	requestsCount   int
	lastRequestTime time.Time
}

func (s *rateLimiter) Wait() {
	if s.requestsCount == 3 {
		secs := time.Since(s.lastRequestTime).Seconds()
		ms := int((1 - secs) * 1000)
		if ms > 0 {
			duration := time.Duration(ms * int(time.Millisecond))
			//fmt.Println("attempted to make more than 3 requests per second, "+
			//"sleeping for", ms, "ms")
			time.Sleep(duration)
		}

		s.requestsCount = 0
	}
}

func (s *rateLimiter) Update() {
	s.requestsCount++
	s.lastRequestTime = time.Now()
}
