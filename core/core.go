package core

import (
	"github.com/sydnash/lotou/log"
	"sync"
)

type Service interface {
	Send(m *Message)
	SetId(id uint)
}

type manager struct {
	id      uint
	mutex   sync.Mutex
	dictory map[uint]Service
	nameDic map[string]uint
}

var c *manager

func init() {
	c = new(manager)
	c.dictory = make(map[uint]Service)
	c.nameDic = make(map[string]uint)
}

func GetService(id uint) Service {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	ser, ok := c.dictory[id]
	if !ok {
		log.Warn("GetService: service %d is not exist.\n", id)
	}
	return ser
}
func RegisterService(s Service) uint {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.id++
	c.dictory[c.id] = s
	s.SetId(c.id)
	return c.id
}

func Send(dest uint, src uint, data ...interface{}) bool {
	return send(dest, src, MSG_TYPE_NORMAL, "go", data...)
}
func SendSocket(dest, src uint, data ...interface{}) bool {
	return send(dest, src, MSG_TYPE_NORMAL, "socket", data...)
}

func getIdByName(name string) (uint, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	id, ok := c.nameDic[name]
	if !ok {
		log.Warn("getIdByName: service: %s is not exist.", name)
		return 0, false
	}
	return id, true
}
func SendName(name string, src uint, data ...interface{}) bool {
	id, ok := getIdByName(name)
	if !ok {
		return false
	}
	return send(id, src, MSG_TYPE_NORMAL, "go", data...)
}

func Name(id uint, name string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if _, ok := c.nameDic[name]; ok {
		log.Warn("Name: service %d is not exist.\n", id)
		return false
	}
	c.nameDic[name] = id
	return true
}

func Close(dest uint, src uint) bool {
	ret := send(dest, src, MSG_TYPE_CLOSE, "go")
	remove(dest)
	return ret
}

func remove(id uint) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.dictory, id)
}

func send(dest, src uint, msgType int, msgEncodeType string, data ...interface{}) bool {
	ser := GetService(dest)
	if ser == nil {
		return false
	}
	m := &Message{dest, src, msgType, msgEncodeType, data}
	ser.Send(m)
	return true
}

const (
	MSG_TYPE_NORMAL = iota
	MSG_TYPE_CLOSE
)

type Message struct {
	Dest          uint
	Src           uint
	Type          int
	msgEncodeType string
	Data          []interface{}
}

func SafeGo(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Error("%s", err)
			}
		}()
		f()
	}()
}
func SafeCall(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Error("%s", err)
			}
		}()
		f()
	}()
}
