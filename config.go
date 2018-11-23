package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
)

const (
	flagNameA                = "a"
	flagNamePort             = "port"
	flagNameSocsPort         = "socs-port"
	flagNameAuthFile         = "auth-file"
	defaultFlagPortValue     = 32946
	defaultFlagSocsPortValue = 32947
)

var port int
var socsPort int
var authFile string
var credentials map[string]string
var mutex = &sync.Mutex{}

type arrayFlags []string

func (i *arrayFlags) String() string {
	if *i == nil || len(*i) == 0 {
		return "[]"
	}
	var b bytes.Buffer

	if i != nil {
		for index, each := range *i {
			if index == 0 {
				b.WriteString(each)
			} else {
				b.WriteString(", ")
				b.WriteString(each)
			}

		}

	}
	return b.String()
}

func configGetHttpPort() int {
	return port
}

func configGetSocsPort() int {
	return socsPort
}

func hasUser() bool {
	var result bool
	mutex.Lock()
	result = len(credentials) > 0
	mutex.Unlock()
	return result
}

func configIsUserAllowed(userName string, userPass string) bool {
	var result bool
	mutex.Lock()
	if len(credentials) == 0 { //no user
		result = true
	} else {
		storedPass := credentials[userName]
		if len(storedPass) != 0 {
			if storedPass == userPass {
				result = true
			}
		}
	}
	mutex.Unlock()
	return result
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func initConfig() error {

	var aFlags arrayFlags
	var authFileFlags arrayFlags

	if flag.Lookup(flagNameA) == nil {
		flag.Var(&aFlags, flagNameA, "HTTP basic auth username and password, for example user:password, for multiple users just repeat -a parameters ,such as: -a user1:password1 -a user2:password2")
	}
	if flag.Lookup(flagNamePort) == nil {
		flag.IntVar(&port, flagNamePort, defaultFlagPortValue, "Provide a port to connect to")
	}
	if flag.Lookup(flagNameSocsPort) == nil {
		flag.IntVar(&socsPort, flagNameSocsPort, defaultFlagSocsPortValue, "Provide a port to connect to")
	}
	if flag.Lookup(flagNameAuthFile) == nil {
		flag.StringVar(&authFile, flagNameAuthFile, "", "A path of HTTP basic auth file,\"username:password\" on each line in a file, the file contains, for example:\n user1:pass1 \n user2:pass2 \n userN:passN")
	}

	flag.Parse()

	var errPort error
	errPort = checkPort(port)
	if errPort != nil {
		return errPort
	}
	errPort = checkPort(socsPort)
	if errPort != nil {
		return errPort
	}

	credentials = make(map[string]string)
	fillCredentials(credentials, &aFlags)
	/*	for key, value := range credentials {
			fmt.Println("user:", key, "pass:", value)
		}
	*/
	if len(authFile) > 0 {
		authFileFlags, _ = getArrayFlagsFromFile(authFile)
		if authFileFlags != nil {
			fillCredentials(credentials, &authFileFlags)
		}
	}
	return nil
}

func checkPort(port int) error {
	if port > 65535 || port < 1 {
		errMessage := fmt.Sprintf("TCP port must be in the range 1 - 65535,  the port[%d] is wrong\n", port)
		return errors.New(errMessage)
	}
	return nil
}

func getArrayFlagsFromFile(pPath string) (arrayFlags, error) {
	file, err := os.Open(pPath)
	if err != nil {
		fmt.Println("Faled to open file:", pPath)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	result := make(arrayFlags, 0)
	for scanner.Scan() {
		_ = result.Set(scanner.Text())
	}
	fmt.Println(result)
	if err := scanner.Err(); err != nil {
		fmt.Println("can't scan '"+pPath+"' file , error:", err)
		return nil, err
	}

	return result, nil
}

func fillCredentials(pCredentials map[string]string, aFlags *arrayFlags) {
	for _, each := range *aFlags {
		charIndex := strings.Index(each, ":")
		if charIndex == -1 || charIndex == 0 || charIndex+1 == len(each) { //A wrong format of the 'user:pass' pair
			continue
		}
		key := each[:charIndex]
		var value string
		if charIndex+1 < len(each) {
			value = each[charIndex+1:]
			pCredentials[key] = value
		}

	}
}

func copyCredentials(pCopyCredentials map[string]string) {
	mutex.Lock()
	if len(credentials) == 0 { //no user
		mutex.Unlock()
		return
	}

	for key, value := range credentials {
		newValue := make([]byte, len(value))
		newKey := make([]byte, len(key))
		copy(newKey, key)
		copy(newValue, value)
		pCopyCredentials[string(newKey)] = string(newValue)
	}
	mutex.Unlock()
	return
}
