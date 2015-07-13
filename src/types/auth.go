package types

type Auth struct {
	Id         int32  //User Id
	Uuid       string //UUID
	BindId     string //bind id
	UniqueId   uint64 //user, alliance, room and others unique id. use for chat
	Cert       string //证书
	Gsid       int32  //所在服务器
	AuthType   uint8  //登录方式
	CreateTime int64  //创建时间
	Email      string //邮箱
}
