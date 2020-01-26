package repo

import (
	"demo-plaform/services/repo"
	"log"
	"testing"
)

func TestCreateRepo(t *testing.T) {
	err := repo.Create("/var/srv/git/asuis/test2.git", nil)
	if err != nil {
		log.Fatalf("%v", err)
	}
}