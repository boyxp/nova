package model

import "github.com/boyxp/nova/database"

func Goods() *database.Orm {
	return database.NewOrm()
}
