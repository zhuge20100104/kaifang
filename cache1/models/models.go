package models

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
