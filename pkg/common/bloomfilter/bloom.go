package bloomfilter

import (
	"math"
)

type RedisClient interface {
	SetBit(key string, offset int64) error
	GetBit(key string, offset int64) (int64, error)
}
type BloomFilter struct {
	Key            string
	BloomHashFuncs []BloomHashFunc //保存hash函数
	BitArraySize   int64           //位数组的长度
	Redis          RedisClient
}

// New 创建一个布隆过滤器
// 误判率：p = (1 - e^(-kn/m))^k
// k 是哈希函数的数量，n 是插入到布隆过滤器中的元素数量，m 是布隆过滤器的大小
func New(key string, expectedNumOfItems int64, falsePositiveRate float64, redisClient RedisClient) *BloomFilter {
	// 根据预计要存储的元素数量和可接受的误判率来计算位数组的长度
	m := -1 * (float64(expectedNumOfItems) * math.Log(falsePositiveRate)) / (math.Log(2) * math.Log(2))
	bitArraySize := int64(m)

	return &BloomFilter{Key: key, BloomHashFuncs: NewFunc(), BitArraySize: bitArraySize, Redis: redisClient}
}

// Add bloom中添加一个值
func (b *BloomFilter) Add(str string) error {
	for _, f := range b.BloomHashFuncs {
		offset := f(str)
		err := b.Redis.SetBit(b.Key, offset)
		if err != nil {
			return err
		}
	}
	return nil
}

// Exist  查看该值是否在bloom过滤器中   1是在，0是不在,-1错误
func (b *BloomFilter) Exist(str string) (int, error) {
	var a int64 = 1
	for _, f := range b.BloomHashFuncs {
		offset := f(str)
		bitValue, err := b.Redis.GetBit(b.Key, offset)
		if err != nil {
			return -1, err
		}
		if bitValue != a {
			return 0, nil
		}
	}
	return 1, nil
}
