package model

import "github.com/boyxp/nova/database"

func User() *database.Orm {
	return database.NewOrm()
}
