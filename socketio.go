//|------------------------------------------------------------------
//|        __
//|     __/  \
//|  __/  \__/_
//| /  \__/    \
//|/\__/CellGo /_
//|\/_/NetFW__/  \
//|  /\__ _/  \__/
//|  \/_/  \__/_/
//|    /\__/_/
//|    \/_/
//| ------------------------------------------------------------------
//| Cellgo Framework tcpip/socketio file
//| All rights reserved: By cellgo.cn CopyRight
//| You are free to use the source code, but in the use of the process,
//| please keep the author information. Respect for the work of others
//| is respect for their own
//|-------------------------------------------------------------------
// Author:Tommy.Jin Dtime:2016-08-06

package tcpip

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/googollee/go-socket.io"
)

type socketConf struct {
	Conn    string `json:"Conn"`    //Connection function name
	Disconn string `json:"Disconn"` //Disconnection function name
	Error   string `json:"Error"`   //Error function name
	Auth    string `json:"Auth"`    //Auth function name
	Push    string `json:"Push"`    //Push content function name
	Pull    string `json:"Pull"`    //Pull content function name
}

func RunSocketIO() {
	for _, v := range Tcp[SOCKETIO] {
		go func(v *TcpRun) {
			socketConf := &socketConf{}
			err := json.Unmarshal([]byte(v.TcpConf), socketConf)
			if err != nil {
				log.Fatal("socketio [", v.TcpName, "] error:", err)
			}
			checkDefault(socketConf)

			server := v.Handle.(*socketio.Server)
			server.On(socketConf.Conn, func(so socketio.Socket) {
				log.Println("on connection")
				so.On(socketConf.Pull, func(msg string) string {
					return msg
				})
				so.On(socketConf.Disconn, func() {
					log.Println("on disconnect")
				})
			})
			server.On(socketConf.Error, func(so socketio.Socket, err error) {
				log.Println("error:", err)
			})
			http.Handle(v.Route, server)
			log.Println(v.TcpName, "Serving at", v.Addr, "to", v.Route)
			log.Fatal(http.ListenAndServe(v.Addr, nil))
		}(v)
	}
}

func checkDefault(s *socketConf) {
	switch {
	case s.Conn == "":
		s.Conn = "connection"
		fallthrough
	case s.Disconn == "":
		s.Disconn = "disconnection"
		fallthrough
	case s.Error == "error":
		s.Error = "error"
		fallthrough
	case s.Auth == "":
		s.Auth = "auth"
		fallthrough
	case s.Push == "":
		s.Push = "push"
		fallthrough
	case s.Pull == "":
		s.Pull = "pull"
	}
}
