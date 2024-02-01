package wsutil

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"wstester/model"
	"wstester/pkg/log"
	"wstester/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type WsUtil struct {
	url             string
	conn            *websocket.Conn
	sendChan        chan []byte
	respChan        chan []byte
	sendMessagePool sync.Map
	respPool        sync.Map
	keepAliveMsg    model.WsMsg
	idGen           utils.IDGenerator
	validate        *validator.Validate
	connected       bool

	sendCtx       context.Context
	receiveCtx    context.Context
	keepAliveCtx  context.Context
	handleRespCtx context.Context
	cancleFunc    context.CancelFunc

	respPoolLock sync.RWMutex
	sendPoolLock sync.RWMutex
}

func NewWsUtil(url string, keepAlive bool) (*WsUtil, error) {
	parent, cancleFunc := context.WithCancel(context.Background())
	sendCtx, _ := context.WithCancel(parent)
	receiveCtx, _ := context.WithCancel(parent)
	keepAliveCtx, _ := context.WithCancel(parent)
	handleRespCtx, _ := context.WithCancel(parent)
	w := &WsUtil{
		url:             url,
		sendChan:        make(chan []byte, 1000),
		respChan:        make(chan []byte, 1000),
		sendMessagePool: sync.Map{},
		respPool:        sync.Map{},
		idGen:           utils.NewIDGenerator(),
		validate:        validator.New(),
		respPoolLock:    sync.RWMutex{},
		sendPoolLock:    sync.RWMutex{},
		sendCtx:         sendCtx,
		receiveCtx:      receiveCtx,
		keepAliveCtx:    keepAliveCtx,
		handleRespCtx:   handleRespCtx,
		cancleFunc:      cancleFunc,
	}
	if err := w.Conn(url, keepAlive, nil); err != nil {
		log.Errorf("err: %s", err.Error())
		return nil, errors.Wrap(err, "conn error")
	}
	return w, nil
}

func (w *WsUtil) SetKeepAliveMsg(keepAliveMsg model.WsMsg) {
	w.keepAliveMsg = keepAliveMsg
}

func (w *WsUtil) Conn(url string, keepAlive bool, requestHeader map[string][]string) error {
	var wsClient = websocket.Dialer{
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		HandshakeTimeout: 30 * time.Second,
	}

	/*
		内网代理：4010812200、3121549363
		外网代理：4252400967
	*/
	url = fmt.Sprintf("%s%s", url, "&proxy1=3121549363&proxy2=4010812200&proxy3=4252400967")

	conn, _, err := wsClient.Dial(url, requestHeader)
	if err != nil {
		return errors.Wrapf(err, "ws dial %s error", url)
	}

	log.Info("success to dial ws")
	w.conn = conn
	w.connected = true
	go w.send()
	go w.receive()
	go w.handleResp()
	if keepAlive {
		go w.keepAlive()
	}
	return nil
}

func (w *WsUtil) SendMsg(msg interface{}) (int64, error) {
	msg.ID = w.idGen.Generate()
	log.Infof("msg id: %d", msg.ID)

	if !w.IsConnected() {
		return 0, errors.New("login first")
	}

	if err := w.validate.Struct(msg); err != nil {
		return 0, errors.Wrapf(err, "send msg error")
	}

	msgContent, err := json.Marshal(msg)
	if err != nil {
		return 0, errors.Wrapf(err, "send msg error")
	}
	w.sendMessagePool.Store(msg.ID, msg)
	w.sendChan <- msgContent

	return msg.ID, nil
}

func (w *WsUtil) send() {
	log.Info("Ready to send ws message")

	for {
		select {
		case msg := <-w.sendChan:
			if err := w.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Errorf("send msg error: %s", err.Error())
				w.CloseConnected()
			}
		case <-w.sendCtx.Done():
			log.Errorf("stop to send message...")
			return
		}
	}
}

func (w *WsUtil) receive() {
	log.Info("Ready to receive ws message")
	t := time.NewTicker(500 * time.Millisecond)

	for {
		select {
		case <-t.C:
			messageType, resp, err := w.conn.ReadMessage()
			if err != nil {
				log.Errorf("read message error: %s", err.Error())
				log.Errorf("stop the ws conn")
				w.CloseConnected()
				continue
			}

			if messageType != websocket.TextMessage {
				log.Errorf("message is not supported: %d", messageType)
				continue
			}

			w.respChan <- resp
			t.Reset(500 * time.Millisecond)
		case <-w.receiveCtx.Done():
			log.Errorf("stop to receive message...")
			return
		}
	}
}

func (w *WsUtil) handleResp() {
	log.Info("Ready to handle response.")
	for {
		select {
		case resp := <-w.respChan:
			var temp map[string]interface{}
			if err := json.Unmarshal(resp, &temp); err != nil {
				log.Errorf("unmarshal resp error: %s", err.Error())
				continue
			}

			if _, exist := temp["id"]; !exist {
				log.Errorf("resp is not exist id element")
				continue
			}
			id := utils.ParseToInt64(temp["id"])

			if _, exist := w.sendMessagePool.Load(id); !exist {
				log.Errorf("message pool is not exist, id: %d", id)
				continue
			}

			w.respPool.Store(id, resp)
		case <-w.handleRespCtx.Done():
			log.Errorf("stop handle resp")
			return
		}
	}
}

func (w *WsUtil) getRespByID(id int64) []byte {
	w.respPoolLock.TryRLock()
	defer w.respPoolLock.RUnlock()

	r, ok := w.respPool.Load(id)
	if !ok {
		return nil
	}

	if resp, ok := r.([]byte); ok {
		return resp
	}
	return nil
}

func (w *WsUtil) keepAlive() {
	log.Info("Start to keep alive")

	t := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-t.C:
			id, err := w.SendMsg(w.keepAliveMsg)
			if err != nil {
				log.Errorf("keep alive error: %s", err.Error())
			}
			var (
				resp        []byte
				timeoutFlag bool
			)
			log.Infof("send keepalive message")
			tryGetRespInternal := time.NewTicker(1 * time.Second)
			timeout := time.NewTicker(10 * time.Second)
			for {
				select {
				case <-tryGetRespInternal.C:
					resp = w.getRespByID(id)
					if resp != nil {
						timeoutFlag = true
					}
				case <-timeout.C:
					log.Errorf("get keep alive response time out")
					w.cancleFunc()
					timeoutFlag = true
				}
				if timeoutFlag {
					break
				}
			}

			if resp == nil {
				continue
			}

			result := model.RespKeepAlive{}
			if err := json.Unmarshal(resp, &result); err != nil {
				log.Errorf("keep alive error, unmarshal error: %s", err.Error())
			}
			if result.ID != id {
				log.Errorf("msg error, id is not the same as before")
			}
			log.Infof("keepAlive result: %+v", result)
		case <-w.keepAliveCtx.Done():
			log.Errorf("stop to keep alive...")
			return
		}
	}
}

func (w *WsUtil) GetResp(msgID int64) []byte {
	var (
		resp        []byte
		timeoutFlag bool
	)

	tryGetRespInternal := time.NewTicker(1 * time.Second)
	timeout := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-tryGetRespInternal.C:
			resp = w.getRespByID(msgID)
			if resp != nil {
				timeoutFlag = true
			}
		case <-timeout.C:
			log.Errorf("get response time out")
			timeoutFlag = true
		}
		if timeoutFlag {
			break
		}
	}
	go w.delRespPool(msgID)
	return resp

}

func (w *WsUtil) delRespPool(msgID int64) {
	w.respPoolLock.TryLock()
	defer w.respPoolLock.Unlock()

	w.respPool.Delete(msgID)
}

func (w *WsUtil) IsConnected() bool {
	return w.connected
}

func (w *WsUtil) CloseConnected() {
	w.cancleFunc()
	w.connected = false
}
