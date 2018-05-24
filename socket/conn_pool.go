package socket

import (
	"sync"
	"net"
	"fmt"
)

// golang内建的map不是并发安全的
// 解决方案一     （csp）
// 解决方案二     （加锁）
// 解决方案三     （sync.Map）

type ConnPool struct {
	sync.RWMutex
	ConnMap map[string]net.Conn
}

func (cp *ConnPool) Get(key string) net.Conn {
	cp.Lock()
	result := cp.ConnMap[key]
	cp.Unlock()
	return result
}

func (cp *ConnPool) Delete(key string) {
	fmt.Println("BEFORE LOCK")
	cp.Lock()
	fmt.Println("d1")
	delete(cp.ConnMap, key)
	fmt.Println("d2")
	cp.Unlock()
	fmt.Println("AFTER LOCK")
}

func (cp *ConnPool) Set(key string, conn net.Conn) {
	cp.Lock()
	cp.ConnMap[key] = conn
	cp.Unlock()
}
