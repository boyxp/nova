package model

import "github.com/boyxp/nova/database"

//简单模型，直接操作表
var User = database.Model{"user"}

/*
//复杂模型，继承附加方法
type ModelUser struct {
	database.Model
}

func (M *ModelUser) MaxUserId() string {
	return M.Field("max(user_id) as max_id").Value("max_id")
}

var User = ModelUser{database.Model{"user"}}
*/
