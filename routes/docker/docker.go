package docker

import (
	"bufio"
	"context"
	"demo-plaform/model/db"
	"demo-plaform/services/docker"
	"demo-plaform/services/user"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)
const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 8192

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Time to wait before force close on connection.
	closeGracePeriod = 10 * time.Second
)

var upgrader = websocket.Upgrader{}

func internalError(ws *websocket.Conn, msg string, err error) {
	log.Println(msg, err)
	_ = ws.WriteMessage(websocket.TextMessage, []byte("Internal server error."))
}

func ServeWs(w http.ResponseWriter, r *http.Request, cloudID string) {

	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	defer ws.Close()

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	cli.NegotiateAPIVersion(ctx)

	c, err := cli.ContainerInspect(ctx, cloudID)

	if err != nil {
		panic(err)
	}

	if !c.State.Running {
		fmt.Printf("You cannot attach to a stopped container, start it first")
	}

	if c.State.Paused {
		fmt.Printf("You cannot attach to a paused container, unpause it first")
	}

	res, err := cli.ContainerAttach(ctx, c.Name, types.ContainerAttachOptions{
		Stream:     true,
		Stdin:      true,
		Stdout:     true,
		Stderr:     true,
		DetachKeys: "",
		Logs:       false,
	})

	if err != nil {
		panic(err)
	}

	defer res.Close()

	go pumpStdin(ws, res)
	go pumpStdout(ws, res.Reader, nil)

}

func pumpStdin(ws *websocket.Conn, res types.HijackedResponse) {
	defer ws.Close()
	ws.SetReadLimit(maxMessageSize)
	_ = ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { _ = ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		message = append(message, '\n')
		if _, err := res.Conn.Write(message); err != nil {
			break
		}
	}
}

func pumpStdout(ws *websocket.Conn, r *bufio.Reader, done chan struct{}) {
	defer func() {
	}()
	s := bufio.NewScanner(r)
	for s.Scan() {
		_ = ws.SetWriteDeadline(time.Now().Add(writeWait))
		if err := ws.WriteMessage(websocket.TextMessage, s.Bytes()); err != nil {
			ws.Close()
			break
		}
	}
	if s.Err() != nil {
		log.Println("scan:", s.Err())
	}
	close(done)

	_ = ws.SetWriteDeadline(time.Now().Add(writeWait))
	_ = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	time.Sleep(closeGracePeriod)
	ws.Close()
}

func ping(ws *websocket.Conn, done chan struct{}) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				log.Println("ping:", err)
			}
		case <-done:
			return
		}
	}
}


func CreateDocker(ctx *gin.Context) {

	var json docker.DockerCreateForm

	sign, _ := ctx.Get("u")

	if err := ctx.ShouldBindJSON(&json); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ID, err := docker.CreateContainer(&json, &db.User{
		Id: sign.(user.SignedData).Ac,
	})

	if err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"containerID": ID,
	})

}

func ListDocker(ctx *gin.Context) {
	sign, _ := ctx.Get("u")
	order := ctx.Param("order")
	pageSize, err:= strconv.Atoi(ctx.Param("pageSize"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	page, err:= strconv.Atoi(ctx.Param("page"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	list, err := docker.ListContainers(&db.User{Id:sign.(user.SignedData).Ac},page,pageSize, order)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, list)
	return
}

func StopDocker(ctx *gin.Context) {
	sign, _ := ctx.Get("u")
	ID := ctx.Param("cloud_id")
	err := docker.StopContainer(ID, &db.User{
		Id: sign.(user.SignedData).Ac,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusOK)
	return
}

func StartDocker(ctx *gin.Context) {
	sign, _ := ctx.Get("u")
	ID := ctx.Param("cloud_id")
	err := docker.StartContainer(ID, &db.User{
		Id: sign.(user.SignedData).Ac,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusOK)
	return
}

func DelDocker(ctx *gin.Context) {
	sign, _ := ctx.Get("u")
	ID := ctx.Param("cloud_id")
	err := docker.RmContainer(ID, &db.User{
		Id: sign.(user.SignedData).Ac,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusOK)
	return
}

func InfoDocker(ctx *gin.Context) {
	sign, _ := ctx.Get("u")
	ID := ctx.Param("cloud_id")
	res, err := docker.StatusContainer(ID, &db.User{
		Id: sign.(user.SignedData).Ac,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, body)
	return
}

func StatDocker(ctx *gin.Context) {
	sign, _ := ctx.Get("u")
	ID := ctx.Param("cloud_id")
	res, err := docker.StatusContainer(ID, &db.User{
		Id: sign.(user.SignedData).Ac,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, body)
	return
}

func AttachDocker(ctx *gin.Context) {
	sign, _ := ctx.Get("u")
	id := ctx.Param("cloud_id")
	c, err := db.Engine.Count(&db.DockerContainer{ContainerID: id, OwnerID:sign.(user.SignedData).Ac})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if c <= 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	}
	go ServeWs(ctx.Writer, ctx.Request, id)
}