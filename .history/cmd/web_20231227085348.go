package cmd

import (
	"embed"
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
	"os"
	"os/signal"
	"time"
	ws2 "wstester/internal/base/ws"

	"wstester/model"
	"wstester/pkg/log"
	"wstester/utils/wsutil"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const UIStaticPath = "public"

var webCommand = &cobra.Command{
	Use:   "web",
	Short: "run web restful api",
	Run:   runWeb,
}

func init() {
	webCommand.Flags().StringP("filename", "f", "config", "set filename")
	webCommand.Flags().StringArrayP("configpath", "p", []string{"./config"}, "set config path")
	viper.BindPFlags(webCommand.Flags())
	rootCommand.AddCommand(webCommand)
}

type WsHandler struct {
	ws *wsutil.WsUtil
}

var (
	upgrade = &websocket.Upgrader{
		HandshakeTimeout: time.Second * 10,

		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func runWeb(cmd *cobra.Command, args []string) {
	log.Info("run web")
	//var ws WsHandler

	//c := config.NewConfig()
	//c.LoadConfigFile()

	cmd.Flags().GetString("")

	engine := gin.New()
	engine.Use(CORSMiddleware())
	engine.Use(gin.Recovery())
	if err := engine.SetTrustedProxies(nil); err != nil {
		panic(err)
	}
	gin.SetMode(gin.DebugMode)

	//engine.POST("/login", loginFunc(&ws))
	//engine.POST("/send-ws-msg", sendMsg(&ws))
	//engine.GET("/ping", func(c *gin.Context) {
	//	c.String(http.StatusOK, "pong")
	//})

	engine.GET("/ws", func(ctx *gin.Context) {
		var (
			err  error
			conn *websocket.Conn
			//ws   *ws2.WsConn
		)
		if conn, err = upgrade.Upgrade(ctx.Writer, ctx.Request, nil); err != nil {
			return
		}
		ws2.NewWsConn(conn).Start()
		// 使得inChan和outChan耦合起来
		//for {
		//	var data []byte
		//	if data, err = ws.InChanRead(); err != nil {
		//		goto ERR
		//	}
		//	if err = ws.OutChanWrite(data); err != nil {
		//		goto ERR
		//	}
		//}
		//ERR:
		//	ws.CloseConn()
	})

	// engine.GET("/test", uiFile)

	go func() {
		if err := engine.Run("0.0.0.0:" + "8080"); err != nil {
			panic(err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)
	for range quit {
		//if ws.ws != nil && ws.ws.IsConnected() {
		//	ws.ws.CloseConnected()
		//}
		log.Infof("stop...")
		time.Sleep(1 * time.Second)
		return
	}
}

type _resource struct {
	fs embed.FS
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}

//func loginFunc(ws *WsHandler) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		req := &model.ReqLoginParam{}
//		if err := c.Bind(req); err != nil {
//			log.Errorf("bind json error: %s", err.Error())
//			c.JSON(http.StatusBadRequest, gin.H{"result": "bind request param error"})
//			return
//		}
//		_ = req.CheckAndFillDefaultValue()
//
//		log.Infof("login request: %+v", req)
//
//		if ws.ws != nil && ws.ws.IsConnected() {
//			ws.ws.CloseConnected()
//		}
//
//		login := service.NewLogin(req.User, req.Password, "/mesh/user/login")
//		wsUrl, err := login.Login(model.NewServer(req.Ip, req.Port, req.Protocol))
//		if err != nil {
//			log.Errorf("login error: %s", err.Error())
//			c.String(http.StatusInternalServerError, fmt.Sprintf("login error: %s", err.Error()))
//			return
//		}
//
//		ws.ws, err = wsutil.NewWsUtil(wsUrl, true)
//		if err != nil {
//			log.Errorf("login error: %s", err.Error())
//			c.JSON(http.StatusOK, gin.H{"result": fmt.Sprintf("login error: %s", err.Error())})
//			return
//		}
//		ws.ws.SetKeepAliveMsg(model.NewWsMsg(code.MSG_KEEP_ALIVE, model.NewKeepAlive(60, time.Now().UTC().String())))
//
//		log.Infof("-------login ws: %p", ws)
//		c.JSON(http.StatusOK, gin.H{"result": "login success"})
//	}
//}

func sendMsg(ws *WsHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if ws == nil || ws.ws == nil || !ws.ws.IsConnected() {
			c.String(http.StatusOK, "login first")
			return
		}

		msg := model.WsMsg{}
		if err := c.Bind(&msg); err != nil {
			log.Errorf("bind msg error: %s", err.Error())
			c.JSON(http.StatusOK, gin.H{"Msg": "json格式错误"})
			return
		}
		msgID, err := ws.ws.SendMsg(msg)
		if err != nil {
			log.Errorf("send msg error: %s", err.Error())
			c.JSON(http.StatusOK, gin.H{"Msg": err.Error()})
			return
		}

		resp := ws.ws.GetResp(msgID)
		m := make(map[string]interface{})
		if err = json.Unmarshal(resp, &m); err != nil {
			log.Errorf("unmarshal resp error: %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"status": "failed"})
			return
		}
		c.JSON(http.StatusOK, m)
	}
}

// func webPage(c *gin.Context) {
// 	filePath := ""
// 	var file []byte
// 	var err error
// 	filePath = "public/index.html"
// 	c.Header("content-type", "text/html;charset=utf-8")
// 	file, err = ui.Build.ReadFile(filePath)
// 	if err != nil {
// 		log.Errorf("err: %s", err.Error())
// 		c.Status(http.StatusNotFound)
// 		return
// 	}
// 	c.String(http.StatusOK, string(file))
// }

// func uiFile(c *gin.Context) {
// 	log.Infof("123")
// 	path := c.FullPath()
// 	if !strings.ContainsRune(path, '.') {
// 		path = filepath.Join("public", path, "/index.html")
// 	}
// 	log.Infof("path: %s", path)
// 	file, err := ui.Build.ReadFile(path)
// 	if err != nil {
// 		log.Errorf("err: %s", err.Error())
// 		c.Status(http.StatusNotFound)
// 		return
// 	}
// 	c.String(http.StatusOK, string(file))
// 	log.Infof("path: %s", path)
// }
