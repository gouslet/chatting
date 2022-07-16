/*
 * File: /cmd/websocket/client/lient.go                                        *
 * Project: chatting                                                           *
 * Created At: Saturday, 2022/07/16 , 05:12:53                                 *
 * Author: elchn                                                               *
 * -----                                                                       *
 * Last Modified: Saturday, 2022/07/16 , 13:07:18                              *
 * Modified By: elchn                                                          *
 * -----                                                                       *
 * HISTORY:                                                                    *
 * Date      	By	Comments                                                   *
 * ----------	---	---------------------------------------------------------  *
 */
package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func main() {
	var webCient = flag.Bool("web", false, "true for html client")

	flag.Parse()

	if *webCient {
		webClient()
	} else {
		cmdClient()
	}

}

func cmdClient() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c, _, err := websocket.Dial(ctx, "ws://localhost:2022/ws", nil)
	if err != nil {
		panic(err)
	}
	defer c.Close(websocket.StatusInternalError, "Internal errors")

	err = wsjson.Write(ctx, c, "Hello WebSocket Server")
	if err != nil {
		panic(err)
	}

	var v any
	err = wsjson.Read(ctx, c, &v)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Received Response from Server: %v\n", v)
	c.Close(websocket.StatusNormalClosure, "")
}

func webClient() {
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)

	http.ListenAndServe(":8080", nil)
}
