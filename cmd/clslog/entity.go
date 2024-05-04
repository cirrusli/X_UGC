package main

import (
	"time"
	"trpc.group/trpc-go/trpc-go/log"
)

// WujiReq 无极传入参数
type WujiReq struct {
	Uin       uint64
	AdPosID   uint64
	AdID      uint64
	BeginTime string
	EndTime   string
}

// SearchReq 搜索CLS日志传入参数，改为时间戳
type SearchReq struct {
	Uin       uint64
	AdPosID   uint64
	AdID      uint64
	BeginTime int64
	EndTime   int64
}

// ConvertToTimestamp 将时间转换为毫秒时间戳
func (s *SearchReq) ConvertToTimestamp(beginTime, endTime string) {
	bt, err := time.Parse(time.DateTime, beginTime)
	if err != nil {
		log.Info("time.Parse error: %v", err)
	}
	s.BeginTime = bt.UnixNano() / int64(time.Millisecond)

	et, err := time.Parse(time.DateTime, endTime)
	if err != nil {
		log.Info("time.Parse error: %v", err)
	}
	s.EndTime = et.UnixNano() / int64(time.Millisecond)
}

// DyeingLogResp 返回给无极的染色日志
type DyeingLogResp struct {
	// 透传CLS的请求ID
	RequestID string
	Logs      []DyeingLog
}

// DyeingLog 染色日志
type DyeingLog struct {
	// 取上报的Ftime
	Time    string
	Uin     uint64
	AdPosID uint64
	AdID    uint64
	// 漏斗key转译为可读值
	FunnelValue string
	ExtraInfo   string
}

// ClsLog agent接出到cls的上报日志信息
type ClsLog struct {
	Ftime     string
	ExtInfo   string
	Uin       uint64
	AdPosID   uint64
	AdID      uint64
	FunnelKey string
	ExtraInfo string
}
