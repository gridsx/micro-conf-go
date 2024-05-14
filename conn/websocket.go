package conn

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/grdisx/micro-conf-go/utils"
	"github.com/winjeg/go-commons/log"
	"io"
	"math/rand"
	"os"
	"os/signal"
	"strings"
)

var (
	logger    = log.GetLogger(nil)
	socketMgr *WebsocketMgr
)

func InitSocketMgr(servers, appId, token string, port int) *WebsocketMgr {
	return &WebsocketMgr{metaServers: servers, appId: appId, port: port, closed: false, token: token}
}

type WebsocketMgr struct {
	metaServers string
	appId       string
	port        int
	token       string

	closed bool

	recvCh      chan []byte
	sendCh      chan []byte
	interruptCh chan os.Signal
	conn        *websocket.Conn
}

func (m *WebsocketMgr) Start() {
	m.interruptCh = make(chan os.Signal, 1)
	signal.Notify(m.interruptCh, os.Interrupt, os.Kill)
	m.recvCh = make(chan []byte, 1)
	m.sendCh = make(chan []byte, 1)
	m.reconnect()
}

func (m *WebsocketMgr) Send(data []byte) {
	m.sendCh <- data
}

func (m *WebsocketMgr) Recv() chan []byte {
	return m.recvCh
}

func (m *WebsocketMgr) chooseAddr() string {
	return chooseAddr(m.metaServers)
}

func chooseAddr(metaServers string) string {
	addrArr := strings.Split(metaServers, ",")
	if len(addrArr) < 1 {
		logger.Panic("chooseAddr - meta server address incorrect")
	}
	idx := rand.Intn(len(addrArr))
	return addrArr[idx]
}

func (m *WebsocketMgr) composeConnKey() string {
	ip := utils.GetIP()
	return fmt.Sprintf("%s:%s:%d", m.appId, ip, m.port)
}

func (m *WebsocketMgr) connect() {
	for {
		chosenAddr := m.chooseAddr()
		connKey := m.composeConnKey()
		c, resp, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s/api/ws?key=%s",
			chosenAddr, connKey), nil)
		if err != nil {
			logger.Warningln("websocket connect  error: " + err.Error())
		} else {
			_, _ = io.ReadAll(resp.Body)
			logger.Infof("Connect - websocket connected from :  " + connKey)
			m.conn = c
			break
		}
	}
}

func (m *WebsocketMgr) reconnect() {
	if m.conn != nil {
		_ = m.conn.Close()
	}
	m.connect()
	go m.read()
	go m.write()
}

func (m *WebsocketMgr) close() {
	close(m.sendCh)
	close(m.recvCh)
	m.conn.Close()
	m.closed = true
}

func (m *WebsocketMgr) read() {
	for {
		mt, message, err := m.conn.ReadMessage()
		if err != nil {
			logger.Println("read err:", err)
			if !m.closed {
				m.reconnect()
			}
			return
		}
		if mt == websocket.CloseMessage {
			if !m.closed {
				m.reconnect()
			}
			return
		}
		m.recvCh <- message
	}
}

func (m *WebsocketMgr) write() {
	for {
		select {
		case msg := <-m.sendCh:
			err := m.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				logger.Println("write err :", err.Error())
				if !m.closed {
					m.reconnect()
				}
				return
			} else {
				logger.Debugln("send heartbeat data: " + string(msg))
			}
		case <-m.interruptCh:
			logger.Println("interrupt")
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := m.conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			m.close()
			if err != nil {
				logger.Println("write close:", err)
				return
			}
			return
		}
	}
}
