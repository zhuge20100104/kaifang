package models

// TimedData 缓存数据对象，缓存框架的数据接口规范
type TimedData interface {
	GetCacheTime() int64
}

// KFPerson 开房个人对象
type KFPerson struct {
	ID     int    `db:"id"`
	Name   string `db:"name"`
	IDCard string `db:"idcard"`
}

// QueryResult 查询结果结构
type QueryResult struct {
	Value     []KFPerson // 开房者数据切片
	CacheTime int64      // 缓存时间时间戳
	Count     int        // 被查询次数
}

// GetCacheTime 实现 TimedData接口
func (q *QueryResult) GetCacheTime() int64 {
	return q.CacheTime
}
