package nginx

import (
	"fmt"
	"net/http"
)

type Conf struct {
	pid int64
	path string
}

func runNginx()  {
}

func stopNginx() {}

func getStatus() {
	resp, err := http.Get("http://127.0.0.1:3000/nginx")
	if err != nil {

	}
	fmt.Printf("%v", resp)
}