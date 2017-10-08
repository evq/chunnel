package utils

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"io"
	"net"
	"os"
	osu "os/user"
	"strconv"
	"strings"
)

var (
	DEFAULT_SSH_PORT    = 22
	DEFAULT_DOCKER_PORT = 2375
	NEXT_OPEN_PORT      = "127.0.0.1:0"
)

func ParseUserHost(userHost string) (user string, host string, port int, err error) {
	var userIdx, portIdx int

	port = DEFAULT_SSH_PORT
	if userObj, err := osu.Current(); err == nil {
		user = userObj.Username
	}

	if userIdx = strings.Index(userHost, "@"); userIdx != -1 {
		user = userHost[:userIdx]
	}
	if portIdx = strings.Index(userHost, ":"); portIdx != -1 {
		if port, err = strconv.Atoi((userHost[portIdx+1:])); err != nil {
			return
		}
	} else {
		portIdx = len(userHost)
	}
	host = userHost[(userIdx + 1):portIdx]
	return
}

type SSHTunnel struct {
	LocalAddr  string
	RemoteAddr string
	ServerConn *ssh.Client
}

func (tunnel *SSHTunnel) Start() error {
	listener, err := net.Listen("tcp", tunnel.LocalAddr)
	if err != nil {
		return err
	}
	tunnel.LocalAddr = listener.Addr().String()
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go tunnel.forward(conn)
	}
}

func (tunnel *SSHTunnel) forward(localConn net.Conn) {
	remoteConn, err := tunnel.ServerConn.Dial("tcp", tunnel.RemoteAddr)
	if err != nil {
		fmt.Printf("Remote dial error: %s\n", err)
		return
	}

	copyConn := func(writer, reader net.Conn) {
		_, err := io.Copy(writer, reader)
		if err != nil {
			fmt.Printf("io.Copy error: %s", err)
		}
	}

	go copyConn(localConn, remoteConn)
	go copyConn(remoteConn, localConn)
}

func SSHAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}

func NewDockerTunnel(user string, serverHost string, serverPort int) *SSHTunnel {
	sshConfig := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            []ssh.AuthMethod{SSHAgent()},
	}
	serverConn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", serverHost, serverPort), sshConfig)
	if err != nil {
		fmt.Printf("Server dial error: %s\n", err)
		return nil
	}
	tunnel := SSHTunnel{NEXT_OPEN_PORT, fmt.Sprintf("127.0.0.1:%d", DEFAULT_DOCKER_PORT), serverConn}
	go tunnel.Start()

	for tunnel.LocalAddr == NEXT_OPEN_PORT {
	}
	return &tunnel
}
