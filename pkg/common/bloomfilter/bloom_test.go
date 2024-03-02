package bloomfilter

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) SetBit(key string, offset int64) error {
	args := m.Called(key, offset)
	return args.Error(0)
}

func (m *MockRedisClient) GetBit(key string, offset int64) (int64, error) {
	args := m.Called(key, offset)
	return args.Get(0).(int64), args.Error(1)
}

func TestAdd(t *testing.T) {
	mockRedis := new(MockRedisClient)
	mockRedis.On("SetBit", "test", mock.Anything).Return(nil)

	bf := New("test", 1000, 0.01, mockRedis)
	err := bf.Add("teststring")
	assert.Nil(t, err)
}

func TestExist(t *testing.T) {
	mockRedis := new(MockRedisClient)
	mockRedis.On("SetBit", "test", mock.Anything).Return(nil)
	mockRedis.On("GetBit", "test", mock.Anything).Return(int64(1), nil)

	bf := New("test", 1000, 0.01, mockRedis)
	_ = bf.Add("teststring")
	exist, err := bf.Exist("teststring")
	assert.Nil(t, err)
	assert.Equal(t, 1, exist)
}
