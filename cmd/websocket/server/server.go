/*
 * File: /cmd/websocket/server/server.go                                       *
 * Project: chatting                                                           *
 * Created At: Saturday, 2022/07/16 , 04:58:56                                 *
 * Author: elchn                                                               *
 * -----                                                                       *
 * Last Modified: Saturday, 2022/07/16 , 13:58:58                              *
 * Modified By: elchn                                                          *
 * -----                                                                       *
 * HISTORY:                                                                    *
 * Date      	By	Comments                                                   *
 * ----------	---	---------------------------------------------------------  *
 */
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(w, "HTTP, Hello")
	})

	wsHandlerFunc := func(w http.ResponseWriter, req *http.Request) {
		conn, err := websocket.Accept(w, req, &websocket.AcceptOptions{
			OriginPatterns: []string{"localhost:8080"},
		})
		if err != nil {
			log.Println(err)
			return
		}

		defer conn.Close(websocket.StatusInternalError, "Internal errors")

		ctx, cancel := context.WithTimeout(req.Context(), time.Second*10)
		defer cancel()

		var v any
		err = wsjson.Read(ctx, conn, &v)
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("Received client: %v", v)
		err = wsjson.Write(ctx, conn, "Hello Websocket Client")
		if err != nil {
			log.Println(err)
			return
		}

		conn.Close(websocket.StatusNormalClosure, "cross origin WebSocket accepted")
	}
	http.HandleFunc("/ws", cors(wsHandlerFunc))

	log.Fatal(http.ListenAndServe(":2022", nil))
}

func cors(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")                                                                    // 允许访问所有域，可以换成具体url，注意仅具体url才能带cookie信息
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token,Upgrade") //header的类型
		w.Header().Add("Access-Control-Allow-Credentials", "true")                                                            //设置为true，允许ajax异步请求带cookie信息
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")                                     //允许请求方法
		w.Header().Set("content-type", "application/json;charset=UTF-8")                                                      //返回数据格式是json
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		f(w, r)
	}
}
