/*
 * File: /cmd/tcp/user.go                                                      *
 * Project: chatting                                                           *
 * Created At: Thursday, 2022/07/14 , 14:54:43                                 *
 * Author: elchn                                                               *
 * -----                                                                       *
 * Last Modified: Friday, 2022/07/15 , 12:30:39                                *
 * Modified By: elchn                                                          *
 * -----                                                                       *
 * HISTORY:                                                                    *
 * Date      	By	Comments                                                   *
 * ----------	---	---------------------------------------------------------  *
 */
package logic

import (
	"strconv"
	"time"
)

type User struct {
	ID             int         // unique user identification,generated by GenUserID
	Addr           string      // IP and port of users
	EnterAt        time.Time   // the time when an user enter the chatting
	MessageChannel chan string // used for sending messages
}

func (u *User)String() string {
	return u.Addr + ", UID: " + strconv.Itoa(u.ID) + ", Enters At: "  + u.EnterAt.Format("2006-01-02 15:04:05+8000")
}
