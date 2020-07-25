package main

import (
	"log"
	"os/exec"
)

func main() {
	gofmtLocation, err := exec.Command("which", "gofmt").Output()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(gofmtLocation))
}
