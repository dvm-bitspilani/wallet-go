package websocket

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return ws, err
	}
	return ws, nil
}

//func Reader(conn *websocket.Conn) {
//	for {
//		messageType, p, err := conn.ReadMessage()
//		if err != nil {
//			log.Println(err)
//			return
//		}
//		fmt.Println(string(p))
//
//		if err := conn.WriteMessage(messageType, p); err != nil {
//			log.Println(err)
//			return
//		}
//	}
//}
//
//func Writer(conn *websocket.Conn) {
//	for {
//		fmt.Println("SENDING")
//		messageType, r, err := conn.NextReader()
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//		w, err := conn.NextWriter(messageType)
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//		if _, err := io.Copy(w, r); err != nil {
//			fmt.Println(err)
//			return
//		}
//		if err := w.Close(); err != nil {
//			fmt.Println(err)
//			return
//		}
//	}
//}
//
//func ServeWs(w http.ResponseWriter, r *http.Request) {
//	ws, err := Upgrade(w, r)
//	if err != nil {
//		fmt.Fprintf(w, "%+V\n", err)
//	}
//	go Writer(ws)
//	Reader(ws)
//}
