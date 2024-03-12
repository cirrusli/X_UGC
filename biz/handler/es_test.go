package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFindUser(t *testing.T) {
	// 创建一个新的路由
	router := gin.Default()
	router.GET("/search/searchUser", FindUser)

	// 创建一个模拟的请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET",
		"/search/searchUser?searchString=john&pageSize=10&pageIndex=1", nil)
	router.ServeHTTP(w, req)

	// 检查响应状态码
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	// 解析响应体内容
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	// 检查响应体的 'status' 字段是否为 'success'
	if status, ok := response["status"]; !ok || status != "success" {
		t.Errorf("Expected status 'success', but got '%v'", status)
	}

	// 打印出接口响应的内容
	t.Logf("Response: %v\n", response)
}
