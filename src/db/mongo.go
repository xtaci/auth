package db

import (
	"bytes"
	"io"
	"os"
	"runtime"
	"time"

	"gopkg.in/mgo.v2"

	log "github.com/gonet2/libs/nsq-logger"
)

const (
	DEFAULT_MONGODB_URL = "mongodb://127.0.0.1/mydb"
	ENV_MONGODB_URL     = "MONGODB_URL"
	DEFAULT_DBOPS_VALVE = 128
	DEFAULT_MGO_TIMEOUT = 300
)

var (
	_global_ms        *mgo.Session // mongodb session
	_valve_dbops      chan bool
	_valve_dbops_high chan bool

	_high_prio map[string]bool // high priority caller
)

type CQ interface {
	Close()
}

//------------------------------------------------ db high priority
type mgowrap_high struct {
	mgo *mgo.Session
}

func (wrap mgowrap_high) Close() {
	<-_valve_dbops_high
	wrap.mgo.Close()
}

//------------------------------------------------ db normal priority
type mgowrap_normal struct {
	mgo *mgo.Session
}

func (wrap mgowrap_normal) Close() {
	<-_valve_dbops
	wrap.mgo.Close()
}

func init() {
	mongo_host := DEFAULT_MONGODB_URL
	if env := os.Getenv(ENV_MONGODB_URL); env != "" {
		mongo_host = env
	}
	// dial mongodb
	sess, err := mgo.Dial(mongo_host)
	if err != nil {
		log.Critical("mongodb: cannot connect to", mongo_host, err)
		os.Exit(-1)
	}

	// set default session mode to strong for saving player's data
	sess.SetMode(mgo.Strong, true)
	// set a high timout
	sess.SetSocketTimeout(DEFAULT_MGO_TIMEOUT * time.Second)
	// infinite wait cursor
	sess.SetCursorTimeout(0)
	_global_ms = sess

	// value
	valve := DEFAULT_DBOPS_VALVE
	_valve_dbops = make(chan bool, valve)
	_valve_dbops_high = make(chan bool, valve)

	// high priority caller
	_high_prio = make(map[string]bool)
	_high_prio[`agent/AI.LoginProc`] = true
	_high_prio[`agent/net.P_user_login_req`] = true
	_high_prio[`main._flush`] = true
}

//------------------------------------------------ copy connection
// !IMPORTANT!  NEVER FORGET -----> defer ms.Close() <-----
func C(collection string) (CQ, *mgo.Collection) {
	funcName, _, _, _ := runtime.Caller(2)
	caller := runtime.FuncForPC(funcName).Name()

	if _high_prio[caller] {
		_valve_dbops_high <- true
		ms := _global_ms.Copy()
		c := ms.DB("").C(collection)
		return mgowrap_high{ms}, c
	}

	_valve_dbops <- true
	ms := _global_ms.Copy()
	c := ms.DB("").C(collection)
	return mgowrap_normal{ms}, c
}

//------------------------------------------------ copy connection without valve
func _c(collection string) (*mgo.Session, *mgo.Collection) {
	ms := _global_ms.Copy()
	c := ms.DB("").C(collection)
	return ms, c
}

//---------------------------------------------------------- 产生GridFS文件
func SaveFile(filename string, buf []byte) bool {
	ms := _global_ms.Copy()
	defer ms.Close()

	gridfs := ms.DB("").GridFS("fs")

	// 首先删除同名文件
	err := gridfs.Remove(filename)
	if err != nil {
		log.Critical("gridfs", filename, err)
		return false
	}

	// 产生新文件
	file, err := gridfs.Create(filename)
	if err != nil {
		log.Critical("gridfs", filename, err)
		return false
	}

	n, err := file.Write(buf)
	if err != nil {
		log.Critical("gridfs", filename, n, err)
		return false
	}

	err = file.Close()
	if err != nil {
		log.Critical("gridfs", filename, err)
		return false
	}
	log.Info("gridfs", filename, "saved to GridFS!!")
	return true
}

//---------------------------------------------------------- 读取GridFS文件
func LoadFile(filename string) (ok bool, content []byte) {
	ms := _global_ms.Copy()
	defer ms.Close()

	buf := &bytes.Buffer{}
	file, err := ms.DB("").GridFS("fs").Open(filename)
	if err != nil {
		log.Warning("gridfs", filename, err)
		return false, nil
	}

	n, err := io.Copy(buf, file)
	if err != nil {
		log.Error("gridfs", filename, n, err)
		return false, nil
	}

	err = file.Close()
	if err != nil {
		log.Error("gridfs", filename, err)
		return false, nil
	}

	log.Trace("gridfs", filename, "load from GridFS!!")
	return true, buf.Bytes()
}

//---------------------------------------------------------- 删除GridFS文件
func RemoveFile(filename string) bool {
	ms := _global_ms.Copy()
	defer ms.Close()

	gridfs := ms.DB("").GridFS("fs")

	// 删除同名文件
	err := gridfs.Remove(filename)
	if err != nil {
		log.Warning("gridfs", filename, err)
		return false
	}

	log.Trace("gridfs", filename, "removed from GridFS!!")
	return true
}
