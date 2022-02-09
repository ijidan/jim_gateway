package test

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
	"jim_gateway/pkg"
	"net/url"
	"testing"
)

//var connNum = flag.Int("conn_num", 1, "conn num")

func TestWSConn(t *testing.T) {
	var address = fmt.Sprintf("%s:%d", pkg.Conf.Websocket.Host, pkg.Conf.Websocket.Port)
	u := url.URL{Scheme: "ws", Host: address, Path: "/"}
	clientConn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	require.Nil(t, err, err)

	go func() {
		for {
			_, messageContent, err1 := clientConn.ReadMessage()
			require.Nil(t, err1, err1)
			t.Logf("received message:%s", string(messageContent))
		}
	}()
	t.Log("waiting...")
	select {}
}
