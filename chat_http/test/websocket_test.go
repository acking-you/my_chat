package test

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	chatlog "logger/log"
	"net/http"
	"testing"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func TestSlice(t *testing.T) {
	var slice []int
	slice = append(slice, 332, 2323, 43, 24324, 32)
	for _, item := range slice {
		fmt.Println(item)
	}
}

func TestWebsocket(t *testing.T) {
	http.HandleFunc("/ws", echo)
	err := http.ListenAndServe(":6060", nil)
	if err != nil {
		panic(err)
	}
}

func echo(writer http.ResponseWriter, request *http.Request) {
	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		chatlog.Lg().Errorln(err)
		return
	}
	conn.SetPongHandler(func(string) error {
		err := conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		if err != nil {
			return err
		}
		return nil
	})
	log.Printf("new conn:%s", conn.RemoteAddr())

	for {
		if err != nil {
			log.Println("read:", err)
		}
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv:%s", message)
		err = conn.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
