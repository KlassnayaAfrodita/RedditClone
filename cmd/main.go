package main

import (
	"fmt"

	"github.com/KlassnayaAfrodita/RedditClone/pkg/session"
)

func main() {
	session := session.NewSessionRepository()
	fmt.Println(session)
}
