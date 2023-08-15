package pgk

import (
	"fmt"
)

// Rsqdata 定义json结构体
type Rsqdata struct {
	Type  string `json:"Type"`
	Cause string `json:"Cause"`
	User  string `json:"User"`
}

func Reswsdata(users, cause, types string) string {
	rsqdata := Rsqdata{
		Type:  types,
		Cause: cause,
		User:  users,
	}
	// 发送消息
	fmt.Printf("向 Discord bot 发送数据: %v\n", rsqdata)
	stda := SendMessage(rsqdata)
	return stda
}
