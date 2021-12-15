package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"isp-engine/lib"
	"isp-engine/utils"
	"net/http"
	"strconv"
)

var (
	upgrade = websocket.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)
const (
	StatusOK = 10200
	StatusFail = 10500
)
type Response struct {
	StatusCode int
	Message    string
}

func ChangePassword(c *gin.Context) {
	c.Writer.WriteString("ok")
}

func OpenPort(c *gin.Context) {

}

func OpenSsh(c *gin.Context) {
	var err error
	msg := c.DefaultQuery("msg", "")
	cols := c.DefaultQuery("cols", "150")
	rows := c.DefaultQuery("rows", "35")
	col, _ := strconv.Atoi(cols)
	row, _ := strconv.Atoi(rows)
	terminal := lib.Terminal{
		Columns: uint32(col),
		Rows:    uint32(row),
	}
	sshClient, err := lib.DecodedMsgToSSHClient(msg)
	if err != nil {
		c.Error(err)
		return
	}
	if sshClient.IpAddress == "" || sshClient.Password == "" {
		c.Error(&utils.ApiError{Message: "IP地址或密码不能为空", Code: 400})
		return
	}

	conn, err := lib.Upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Error(err)
		return
	}
	err = sshClient.GenerateClient()
	if err != nil {
		conn.WriteMessage(1, []byte(err.Error()))
		conn.Close()
		return
	}
	sshClient.RequestTerminal(terminal)
	sshClient.Connect(conn)
}
