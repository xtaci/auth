package auth_tbl

import (
	"bytes"
	"crypto/md5"
	. "db"
	"encoding/binary"
	"fmt"
	"time"
	. "types"

	log "github.com/gonet2/libs/nsq-logger"
	"gopkg.in/mgo.v2/bson"
)

const (
	COLLECTION = "auth"
)

//----------------------------------------------------------- 新建一个验证(cert,id,gsid, uuid)
func New(uid int32, uuid string, gsid int32, auth_type uint8, unique uint64) *Auth {
	ms, c := C(COLLECTION)
	defer ms.Close()

	auth := &Auth{
		Id:         uid,
		Cert:       _md5(uid),
		Uuid:       uuid,
		BindId:     uuid,
		Gsid:       gsid,
		UniqueId:   unique,
		AuthType:   auth_type,
		CreateTime: time.Now().Unix(),
	}
	c.Insert(auth)
	log.Info("new auth bindid:", auth.BindId, " uuid : ", auth.Uuid, "cert :", auth.CreateTime)
	return auth
}

//----------------------------------------------------------- 通过BindID 查找auth
func FindByBindID(bindid string, auth_type uint8, gsid int32) *Auth {
	ms, c := C(COLLECTION)
	defer ms.Close()

	auth := &Auth{}
	err := c.Find(bson.M{"bindid": bindid, "authtype": auth_type}).One(auth)
	if err != nil {
		log.Info(COLLECTION, "FindByBindID", err, bindid)
		return nil
	}
	return auth
}

func _md5(key interface{}) string {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, key)
	return fmt.Sprintf("%X", md5.Sum(buf.Bytes()))

}
