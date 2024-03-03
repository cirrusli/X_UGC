package service

import (
	"X_UGC/biz/dal/redis"
	"X_UGC/biz/model"
	"X_UGC/pkg/common/bloomfilter"
	"strconv"
	"time"
)

const (
	ExpectedNumOfItems = 10000
	FalsePositiveRate  = 0.01
	ExpireTime         = 180
)

func PushArticleFeed(userid int) (ArticleIDs []int, err error) {
	weightList, err := GetWeightInfo(userid)
	var weights, typeIDs []int
	for _, w := range weightList {
		weights = append(weights, w.Weight)
	}
	//权重选择器
	weightSelector := GetWeightSelector(weights)
	for i := 0; i < 8; i++ {
		time.Sleep(time.Nanosecond)
		index := weightSelector.PickIndex()
		typeIDs = append(typeIDs, index+1)
	}
	//根据用户id创建布隆过滤器
	bloomFilter := bloomfilter.New(model.BLOOM_FILTER+strconv.Itoa(userid),
		ExpectedNumOfItems, FalsePositiveRate, &redis.Bitmap{})
	for _, typeId := range typeIDs {
		member, err := redis.SRandMember(model.ARTICLE_POOL + strconv.Itoa(typeId))
		if err != nil {
			return nil, err
		}
		imember, _ := strconv.Atoi(member)
		exist, err := bloomFilter.Exist(member)
		if err != nil {
			return nil, err
		}
		if exist == 0 && imember > 0 {
			ArticleIDs = append(ArticleIDs, imember)
			err = bloomFilter.Add(member)
			if err != nil {
				return nil, err
			}

		}
	}
	err = redis.Expire(bloomFilter.Key, ExpireTime)
	if err != nil {
		return nil, err
	}
	return
}

// PushArticleByTypeID 根据具体类型进行推荐
func PushArticleByTypeID(userid int, articleTypeId int, bf *bloomfilter.BloomFilter) (ArticleIDs []int, err error) {
	//根据用户id创建布隆过滤器
	bf = bloomfilter.New(model.BLOOM_FILTER+strconv.Itoa(userid),
		ExpectedNumOfItems, FalsePositiveRate, &redis.Bitmap{})
	var count int64 = 8
	members, err := redis.SRandMemberN(model.ARTICLE_POOL+strconv.Itoa(articleTypeId), count)
	if err != nil {
		return nil, err
	}
	for _, member := range members {
		imember, _ := strconv.Atoi(member)
		exist, err := bf.Exist(member)
		if err != nil {
			return nil, err
		}
		if exist == 0 && imember > 0 {
			ArticleIDs = append(ArticleIDs, imember)
			err = bf.Add(member)
			if err != nil {
				return nil, err
			}

		}
	}

	err = redis.Expire(bf.Key, ExpireTime)
	if err != nil {
		return nil, err
	}

	return

}

// todo 单元测试需要依赖注入，会涉及handler层，主要在于需要mock的组件，后面有时间了看看

type RecommendSrvs struct {
	Weight WeightInfoGetter
	Rc     RedisClient
}

type WeightInfoGetter interface {
	GetWeightInfo(userid int) ([]*model.Weight, error)
}

type RedisClient interface {
	SRandMember(key string) (string, error)
	Exist(key string) (int, error)
	Add(key string) error
	Expire(key string, time int) error
}

func (r *RecommendSrvs) PushArticleFeed(uid int, bf *bloomfilter.BloomFilter) (ArticleIDs []int, err error) {
	weightList, err := r.Weight.GetWeightInfo(uid)
	var weights, typeIDs []int
	for _, w := range weightList {
		weights = append(weights, w.Weight)
	}
	//权重选择器
	weightSelector := GetWeightSelector(weights)
	for i := 0; i < 8; i++ {
		time.Sleep(time.Nanosecond)
		index := weightSelector.PickIndex()
		typeIDs = append(typeIDs, index+1)
	}
	//根据用户id创建布隆过滤器
	bf = bloomfilter.New(model.BLOOM_FILTER+strconv.Itoa(uid),
		ExpectedNumOfItems, FalsePositiveRate, &redis.Bitmap{})
	for _, typeId := range typeIDs {
		member, err := r.Rc.SRandMember(model.ARTICLE_POOL + strconv.Itoa(typeId))
		if err != nil {
			return nil, err
		}
		imember, _ := strconv.Atoi(member)
		exist, err := bf.Exist(member)
		if err != nil {
			return nil, err
		}
		if exist == 0 && imember > 0 {
			ArticleIDs = append(ArticleIDs, imember)
			err = bf.Add(member)
			if err != nil {
				return nil, err
			}

		}
	}
	err = r.Rc.Expire(bf.Key, ExpireTime)
	if err != nil {
		return nil, err
	}
	return
}
