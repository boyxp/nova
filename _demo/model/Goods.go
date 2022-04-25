package model

import "github.com/boyxp/nova/database"

func Payment() *database.Orm {
	return database.NewOrm()
}
