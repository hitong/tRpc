package server

import (
	"github.com/hitong/tRpc"
	"log"
	"net"
	"time"
)

func ListenAndServer(mgr *tRpc.TRpcMgr, serverAddr string) {
	ls, err := net.Listen("tcp", serverAddr)
	if err != nil {
		log.Fatal(err)
	}

	go Server(mgr, ls)
}

func Server(tRpcMgr *tRpc.TRpcMgr, listener net.Listener) <-chan bool {
	stop := make(chan bool)
	go func() {
		for {
			newConn, err := listener.Accept()
			if err != nil {
				log.Println(err)
			}

			log.Println("add server key ", newConn.RemoteAddr().String())
			tRpcMgr.NewTRpc(newConn).Server()
		}
	}()

	go func() {
		tick := time.Tick(5 * time.Second)
		for {
			select {
			case <-tick:
				tRpcMgr.Range(func(key, value interface{}) bool {
					tRpc := value.(*tRpc.TRpc)
					if tRpc.BeClosed {
						log.Println("delete server key ", key)
						tRpcMgr.DeleteKey(key.(string))
					}
					return true
				})
			}
		}
	}()

	return stop
}
