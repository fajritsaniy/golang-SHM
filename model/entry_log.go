package model

import "time"

type RequestLog struct {
	StartTime    time.Time
	EndTime      time.Duration
	StatusCode   int
	ClientIP     string
	Method       string
	RelativePath string
	UserAgent    string
}

type ResponseLog struct {
	StatusCode   int
	ResponseBody string
}
