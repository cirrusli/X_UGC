package main

import (
	"fmt"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/json"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"net/http"
	"sync"
	"trpc.group/trpc-go/trpc-go"
	"trpc.group/trpc-go/trpc-go/log"
	"trpc.group/trpc-go/trpc-go/metrics"
)

func main() {
	_ = NewClsHandler()
	// 启动tRPC服务
	s := trpc.NewServer()
	s.Service("trpc.tianshu.dyeing_log.clslog")
	if err := s.Serve(); err != nil {
		log.Fatal(err)
	}
}

type ClsHandler struct {
	ClsService *ClsService
}

func NewClsHandler() *ClsHandler {
	return &ClsHandler{
		ClsService: NewClsService(),
	}
}
func (h *ClsHandler) SearchClsLog(w http.ResponseWriter, r *http.Request) {
	wuji := &WujiReq{}
	err := json.NewDecoder(r.Body).Decode(&wuji)
	if err != nil {
		log.Error("ParseRequest error: ", err)
		return
	}
	req := &SearchReq{
		Uin:     wuji.Uin,
		AdPosID: wuji.AdPosID,
		AdID:    wuji.AdID,
	}
	req.ConvertToTimestamp(wuji.BeginTime, wuji.EndTime)
	resp, err := h.ClsService.SearchClsLog(req)
	if err != nil {
		log.Error("searchClsLog error: ", err)
		return
	}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error("json.NewEncoder error: ", err)
		return
	}
}

type ClsService struct {
	ClsRepo *LogRepo
}

func NewClsService() *ClsService {
	return &ClsService{
		ClsRepo: NewLogRepo(),
	}
}

// SearchClsLog 是搜索CLS日志的方法
func (s *ClsService) SearchClsLog(req *SearchReq) (*DyeingLogResp, error) {
	clsReq := cls.NewSearchLogRequest()
	clsReq.TopicId = common.StringPtr("fill-in-the-topic-id")
	clsReq.From = common.Int64Ptr(req.BeginTime)
	clsReq.To = common.Int64Ptr(req.EndTime)
	clsReq.Query = common.StringPtr(fmt.Sprintf("uin:%d AND ad_pos_id:%d AND ad_id:%d", req.Uin, req.AdPosID, req.AdID))

	clsResp, err := s.ClsRepo.ClsCli.SearchLog(clsReq)
	if err != nil {
		log.Error("SearchClsLog error: %v", err)
		metrics.IncrCounter("searchClsLog-failed", 1)
		return nil, fmt.Errorf("SearchClsLog error: %w", err)
	}

	clsLogs := make([]ClsLog, len(clsResp.Response.Results))
	dyeingLogs := make([]DyeingLog, len(clsResp.Response.Results))

	var wg sync.WaitGroup
	errorChan := make(chan error)
	defer close(errorChan)

	wg.Add(len(clsResp.Response.Results))

	for i, logInfo := range clsResp.Response.Results {
		go processLogInfo(i, logInfo, clsLogs, dyeingLogs, errorChan)
	}
	wg.Wait()

	err = <-errorChan
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal error: %w", err)
	}
	// todo 记得对日志根据层级排序
	resp := &DyeingLogResp{
		RequestID: *clsResp.Response.RequestId,
		Logs:      dyeingLogs,
	}

	if len(dyeingLogs) == 0 {
		return resp, fmt.Errorf("no logs found")
	}

	return resp, nil
}

// processLogInfo 用于处理单个日志信息
func processLogInfo(i int, logInfo *cls.LogInfo, clsLogs []ClsLog, dyeingLogs []DyeingLog, errorChan chan<- error) {
	err := json.Unmarshal([]byte(*logInfo.LogJson), &clsLogs[i])
	if err != nil {
		select {
		case errorChan <- err: // 发送错误
		default: // 不需要重复传入，说明日志里字段就是有问题的
		}
		return
	}

	dyeingLogs[i] = DyeingLog{
		Time:        clsLogs[i].Ftime,
		Uin:         clsLogs[i].Uin,
		AdPosID:     clsLogs[i].AdPosID,
		AdID:        clsLogs[i].AdID,
		FunnelValue: convert(clsLogs[i].FunnelKey),
		ExtraInfo:   clsLogs[i].ExtraInfo,
	}
}

// convert 漏斗key转译为可读值
func convert(s string) string {
	return s
}

type LogRepo struct {
	ClsCli *cls.Client
}

func NewLogRepo() *LogRepo {
	return &LogRepo{
		ClsCli: NewClsClient(),
	}
}
func NewClsClient() *cls.Client {
	c, err := newClsClient()
	if err != nil {
		log.Fatal(err)
	}
	return c
}
func newClsClient() (*cls.Client, error) {
	// todo 这里的credit部分要放在远程配置
	secretId := "fill-in-the-secret-id"
	secretKey := "fill-in-the-secret-key"
	credit := common.NewCredential(secretId, secretKey)
	clientProfile := profile.NewClientProfile()
	clientProfile.HttpProfile.Endpoint = "cls.tencentcloudapi.com"
	return cls.NewClient(credit, "ap-guangzhou", clientProfile)
}
