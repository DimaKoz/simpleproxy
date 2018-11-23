package main

import (
	"fmt"
	"github.com/DimaKoz/goproxy" //forked "github.com/elazarl/goproxy"
	"log"
	"net/http"
	"simpleProxy/ext/auth" //"github.com/elazarl/goproxy/ext/auth"
	"strconv"
)

func main() {

	err := initConfig()
	if err != nil {
		fmt.Print(err)
		return
	}
	var usedPort = configGetPort()
	fmt.Println("port:", usedPort)
	proxy := goproxy.NewProxyHttpServer()
	if hasUser() {
		proxy.OnRequest().HandleConnect(auth.BasicConnect("restricted", func(user, passwd string) bool {
			return configIsUserAllowed(user, passwd)
		}))
	}
	proxy.Verbose = true
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(usedPort), proxy))
}
