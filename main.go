package main

import (
	"bufio"
	"fmt"
	"github.com/evq/chunnel/utils"
	"log"
	"os"
	"os/exec"
)

func usage() {
	fmt.Println("usage: chunnel [user@]hostname[:port] [command]")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 3 {
		usage()
	}
	userHost := os.Args[1]
	cmdProg := os.Args[2]
	cmdArgs := os.Args[3:]

	user, host, port, err := utils.ParseUserHost(userHost)
	if err != nil {
		log.Fatal(err)
	}
	tunnel := utils.NewDockerTunnel(user, host, port)
	os.Setenv("DOCKER_HOST", tunnel.LocalAddr)

	cmd := exec.Command(cmdProg, cmdArgs...)
	cmd.Stdin = os.Stdin
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			fmt.Printf("%s\n", scanner.Text())
		}
	}()
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}
