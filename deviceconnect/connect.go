// Copyright 2020 Northern.tech AS
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
package deviceconnect

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 4 * time.Second
	// Maximum message size allowed from peer.
	maxMessageSize = 8192
	// Time allowed to read the next pong message from the peer.
	defaultPingWait = 10 * time.Second
	//
	httpsProtocol = "https"
	httpProtocol  = "http"
	wssProtocol   = "wss"
	wsProtocol    = "ws"
)

//Websocket connection routine. setup the ping-pong and connection settings
func Connect(serverUrl string, connectUrl string, token string) (ws *websocket.Conn, err error) {
	dialer := *websocket.DefaultDialer
	parsedUrl, err := url.Parse(serverUrl)
	if err != nil {
		return nil, err
	}

	scheme := getWebSocketScheme(parsedUrl.Scheme)
	u := url.URL{Scheme: scheme, Host: parsedUrl.Host, Path: connectUrl}
	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+token)
	ws, _, err = dialer.Dial(u.String(), headers)
	if err != nil {
		return nil, err
	}

	// ping-pong
	ws.SetReadLimit(maxMessageSize)
	ws.SetReadDeadline(time.Now().Add(defaultPingWait))
	ws.SetPingHandler(func(message string) error {
		pongWait, _ := strconv.Atoi(message)
		ws.SetReadDeadline(time.Now().Add(time.Duration(pongWait) * time.Second))
		return ws.WriteControl(websocket.PongMessage, []byte{}, time.Now().Add(writeWait))
	})
	return ws, nil
}

func getWebSocketScheme(scheme string) string {
	if scheme == httpsProtocol {
		scheme = wssProtocol
	} else if scheme == httpProtocol {
		scheme = wsProtocol
	}
	return scheme
}
