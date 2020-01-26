package repo

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func setHeaderCacheForever(ctx *gin.Context) {
	now := time.Now().Unix()
	expires := now + 31536000
	ctx.Writer.Header().Set("Date", fmt.Sprintf("%d", now))
	ctx.Writer.Header().Set("Expires", fmt.Sprintf("%d", expires))
	ctx.Writer.Header().Set("Cache-Control", "public, max-age=31536000")
}

func serviceRpc(ctx *gin.Context, service string)  {
	defer ctx.Request.Body.Close()

	if ctx.Request.Header.Get("Content-Type") != fmt.Sprintf("application/x-git-%s-request", service) {
		ctx.Status(401)
		return
	}
	ctx.Writer.Header().Set("Content-Type", fmt.Sprintf("applicaion/x-git-%s-result", service))

	var (
		reqBody = ctx.Request.Body
		err error
	)

	if ctx.Request.Header.Get("Content-Encoding") == "gzip" {
		reqBody, err = gzip.NewReader(reqBody)
		if err != nil {
			ctx.Status(500)
			log.Fatalf("HTTP.Get fail tocreate gzip reader: %v", err)
			return
		}
	}

	var stderr bytes.Buffer
	cmd := exec.Command("git", service, "--stateless-rpc", getRepoDir(ctx))

	envs := []string{
		"SSH_ORIGINAL_COMMAND=1",
	}

	if service == "receive-pack" {
		cmd.Env = append(os.Environ(), envs...)
	}
	cmd.Dir = ""
	cmd.Stdout = ctx.Writer
	cmd.Stderr = &stderr
	cmd.Stdin = reqBody
	if err = cmd.Run(); err != nil {
		ctx.Status(500)
		log.Fatalf("HTTP.serviceRPC: fail to serve RPC '%s': %v - %s", service, err, stderr.String())
		return
	}
}
func ServiceUploadPack(ctx *gin.Context) {
	serviceRpc(ctx, "upload-pack")
}

func ServiceReceivePack(ctx *gin.Context) {
	serviceRpc(ctx, "receive-pack")
}

func sendFile(ctx *gin.Context, contentType string) {
	dir := getRepoDir(ctx)
	username := ctx.Param("username")
	repo := ctx.Param("repo")
	file := strings.TrimPrefix(ctx.Request.RequestURI, fmt.Sprintf("/v1/git/%s/%s", username, repo))
	reqFile := path.Join(dir, file)
	fi, err := os.Stat(reqFile)
	if os.IsNotExist(err) {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		return
	}

	ctx.Writer.Header().Set("Content-Type", contentType)
	ctx.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", fi.Size()))
	ctx.Writer.Header().Set("Last-Modified", fi.ModTime().Format(http.TimeFormat))
	http.ServeFile(ctx.Writer, ctx.Request, reqFile)
}

func getServiceType(ctx *gin.Context) string {
	serviceType := ctx.Request.FormValue("service")
	if !strings.HasPrefix(serviceType, "git-") {
		return ""
	}
	return strings.TrimPrefix(serviceType, "git-")
}

func gitCommand(dir string, args ...string) []byte {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf(fmt.Sprintf("Git: %v - %s", err, out))
	}
	return out
}

func updateServerInfo(dir string) []byte {
	return gitCommand(dir, "update-server-info")
}

func packetWrite(str string) []byte {
	s := strconv.FormatInt(int64(len(str)+4), 16)
	if len(s)%4 != 0 {
		s = strings.Repeat("0", 4-len(s)%4) + s
	}
	return []byte(s + str)
}
func setHeaderNoCache(ctx *gin.Context) {
	header := ctx.Writer.Header()
	header.Set("Expires", "Fri, 01 Jan 1980 00:00:00 GMT")
	header.Set("Pragma", "no-cache")
	header.Set("Cache-Control", "no-cache, max-age=0, must-revalidate")
}

func getRepoDir(ctx *gin.Context) string {
	username:= ctx.Param("username")
	repo := ctx.Param("repo")
	if !strings.HasSuffix(repo, "git") {
		repo += ".git"
	}
	return fmt.Sprintf("/var/srv/git/%s/%s/", username, repo)
}

func GetInfoRefs(ctx *gin.Context) {
	setHeaderNoCache(ctx)
	service := getServiceType(ctx)
	dir := getRepoDir(ctx)
	if service != "upload-pack" && service != "receive-pack" {
		updateServerInfo(dir)
		sendFile(ctx,"text/plain; charset=utf-8")
		return
	}

	refs := gitCommand(dir, service, "--stateless-rpc", "--advertise-refs", ".")
	ctx.Writer.Header().Set("Content-Type", fmt.Sprintf("application/x-git-%s-advertisement", service))
	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Write(packetWrite("# service=git-" + service + "\n"))
	ctx.Writer.Write([]byte("0000"))
	ctx.Writer.Write(refs)
}

func GetTextFile(ctx *gin.Context) {
	setHeaderNoCache(ctx)
	sendFile(ctx,"text/plain")
}

func GetInfoPacks(ctx *gin.Context) {
	setHeaderCacheForever(ctx)
	sendFile(ctx, "text/plain; charset=utf-8")
}

func GetLooseObject(ctx *gin.Context) {
	setHeaderCacheForever(ctx)
	sendFile(ctx, "application/x-git-loose-object")
}

func GetPackFile(ctx *gin.Context) {
	setHeaderCacheForever(ctx)
	sendFile(ctx, "application/x-git-packed-objects")
}

func GetIdxFile(ctx *gin.Context) {
	setHeaderCacheForever(ctx)
	sendFile(ctx, "application/x-git-packed-objects-toc")
}

var routes = [] struct{
	reg *regexp.Regexp
	handler func(*gin.Context)
}{
	{regexp.MustCompile("/objects/info/alternates$"), GetTextFile},
	{regexp.MustCompile("/objects/info/http-alternates$"), GetTextFile},
	{regexp.MustCompile("/objects/info/packs$"), GetInfoPacks},
	{regexp.MustCompile("/objects/info/[^/]*$"), GetTextFile},
	{regexp.MustCompile("/objects/[0-9a-f]{2}/[0-9a-f]{38}$"), GetLooseObject},
	{regexp.MustCompile("/objects/pack/pack-[0-9a-f]{40}\\.pack$"), GetPackFile},
	{regexp.MustCompile("/objects/pack/pack-[0-9a-f]{40}\\.idx$"), GetIdxFile},
}

func GetObject(ctx *gin.Context) {
	action := ctx.Param("action")
	action = "/objects" + action
	for _, route := range routes {
		m := route.reg.FindStringSubmatch(strings.ToLower(action))
		if m == nil {
			continue
		}
		route.handler(ctx)
		return
	}
	ctx.Status(http.StatusNotFound)
	return
}
