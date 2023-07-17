package bloomfilter

import "X_UGC/dal/redis"

type BloomFilter struct {
	Key            string
	BloomHashFuncs []BloomHashFunc //保存hash函数
}

// New 创建一个布隆过滤器
func New(key string) *BloomFilter {
	return &BloomFilter{Key: key, BloomHashFuncs: NewFunc()}
}

// Add bloom中添加一个值
func (b *BloomFilter) Add(str string) error {
	for _, f := range b.BloomHashFuncs {
		offset := f(str)
		err := redis.SetBit(b.Key, offset)
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
		bitValue, err := redis.GetBit(b.Key, offset)
		if err != nil {
			return -1, err
		}
		if bitValue != a {
			return 0, nil
		}
	}
	return 1, nil
}
