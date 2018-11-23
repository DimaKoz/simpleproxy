package main

import (
	"fmt"
	"github.com/DimaKoz/go-socks5" //forked "github.com/armon/go-socks5"
	"github.com/DimaKoz/goproxy"   //forked "github.com/elazarl/goproxy"
	"github.com/DimaKoz/goproxy/ext/auth"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {

	err := initConfig()
	if err != nil {
		fmt.Print(err)
		return
	}
	var usedPort = configGetHttpPort()
	var usedSocsPort = configGetSocsPort()
	fmt.Println("port:", usedPort)
	proxy := goproxy.NewProxyHttpServer()
	if hasUser() {
		proxy.OnRequest().HandleConnect(auth.BasicConnect("restricted", func(user, passwd string) bool {
			return configIsUserAllowed(user, passwd)
		}))
	}
	proxy.Verbose = true

	go func() {
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(usedPort), proxy))
	}()

	fmt.Println("socs5 port: ", usedSocsPort)

	conf := &socks5.Config{}
	if hasUser() {
		cred := socks5.StaticCredentials{}
		copyCredentials(cred)
		cator := socks5.UserPassAuthenticator{Credentials: cred}
		conf.AuthMethods = []socks5.Authenticator{cator}
	}

	conf.Logger = log.New(os.Stdout, "", log.LstdFlags)
	server, err := socks5.New(conf)
	if err != nil {
		fmt.Print(err)
		return
	}
	//how to send socs5 request: curl --socks5 localhost:32947 --proxy-user user1:pass1 binfalse.de
	//how to see ports: sudo lsof -iTCP -sTCP:LISTEN -n -P
	if err := server.ListenAndServe("tcp", /*"127.0.0.1"+*/ ":"+strconv.Itoa(usedSocsPort)); err != nil {
		fmt.Print(err)
		return
	}
}
