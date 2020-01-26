package repo

import (
	"demo-plaform/services/repo"
	"fmt"
	"io"
	"log"
	"testing"
)

func TestSearchRepoDir(t *testing.T) {
	list, err := repo.SearchDir("/var/srv/git/asuis/test.git")
	if err != nil {
		log.Fatalf("%v", err)
	}
	for _, item := range list {
		log.Printf("Name %s", item.Name())
	}
}

func TestGetRawFile(t *testing.T)  {
	list, err := repo.SearchDir("/var/srv/git/asuis/test.git")
	if err != nil {
		log.Fatalf("%v", err)
	}
	for _, item := range list {
		log.Printf("Name %s", item.Name())
		blob ,err:=repo.GetRawFile("/var/srv/git/asuis/test.git", item.ID.String())
		if err != nil {
			log.Fatalf("%v", err)
		}
		reader,err := blob.Data()
		if err != nil {
			log.Fatalf("%v", err)
		}
		buf := make([]byte, 1024)
		for {
			_, err := reader.Read(buf) //读到一个换行就结束
			if err == io.EOF {                  //io.EOF 表示文件的末尾
				break
			}
			fmt.Print(buf)
		}
	}
}