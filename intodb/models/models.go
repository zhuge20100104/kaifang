package models

// KFPerson 开房个人对象
type KFPerson struct {
	ID     int    `db:"id"`
	Name   string `db:"name"`
	IDCard string `db:"idcard"`
}
