package service

import (
	"X_UGC/biz/model"
	"X_UGC/pkg/common/bloomfilter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"strconv"
	"testing"
)

type MockWeightInfoGetter struct {
	mock.Mock
}

func (m *MockWeightInfoGetter) GetWeightInfo(userid int) ([]*model.Weight, error) {
	args := m.Called(userid)
	return args.Get(0).([]*model.Weight), args.Error(1)
}

type MockRedisClient struct {
	mock.Mock
}

func TestPushArticleFeed(t *testing.T) {
	testCases := []struct {
		name     string
		userID   int
		expected int
	}{
		{"Test Case 1", 1, 8},
		// 添加更多的测试用例
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRedisClient := new(MockRedisClient)
			mockWeightInfoGetter := new(MockWeightInfoGetter)

			// 模拟 GetWeightInfo 的返回值
			mockWeightInfoGetter.On("GetWeightInfo", tc.userID).Return([]*model.Weight{{Weight: 10}}, nil)

			// 模拟 Redis 的行为
			mockRedisClient.On("SRandMember", mock.Anything).Return("1", nil)
			mockRedisClient.On("Exist", mock.Anything).Return(0, nil)
			mockRedisClient.On("Add", mock.Anything).Return(nil)
			mockRedisClient.On("Expire", mock.Anything, mock.Anything).Return(nil)

			bloomFilter := bloomfilter.New(model.BLOOM_FILTER+strconv.Itoa(tc.userID),
				ExpectedNumOfItems, FalsePositiveRate, mockRedisClient)

			recommendSrvs := &RecommendSrvs{
				Weight: mockWeightInfoGetter,
				Rc:     mockRedisClient,
			}

			articleIDs, err := recommendSrvs.PushArticleFeed(tc.userID, bloomFilter)

			assert.Nil(t, err)
			assert.Equal(t, tc.expected, len(articleIDs))
		})
	}
}
func TestPushArticleByTypeID(t *testing.T) {
	testCases := []struct {
		name          string
		userID        int
		articleTypeID int
		expected      int
	}{
		{"Test Case 1", 1, 1, 8},
		// 添加更多的测试用例
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRedisClient := new(MockRedisClient)
			bloomFilter := bloomfilter.New(model.BLOOM_FILTER+strconv.Itoa(tc.userID),
				ExpectedNumOfItems, FalsePositiveRate, mockRedisClient)

			articleIDs, err := PushArticleByTypeID(tc.userID, tc.articleTypeID, bloomFilter)

			assert.Nil(t, err)
			assert.Equal(t, tc.expected, len(articleIDs))
		})
	}
}
func (m *MockRedisClient) SRandMember(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *MockRedisClient) Exist(key string) (int, error) {
	args := m.Called(key)
	return args.Int(0), args.Error(1)
}

func (m *MockRedisClient) Add(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *MockRedisClient) Expire(key string, time int) error {
	args := m.Called(key, time)
	return args.Error(0)
}
func (m *MockRedisClient) GetBit(key string, offset int64) (int64, error) {
	args := m.Called(key, offset)
	return int64(args.Int(0)), args.Error(1)
}
func (m *MockRedisClient) SetBit(key string, offset int64) error {
	args := m.Called(key, offset)
	return args.Error(0)
}
