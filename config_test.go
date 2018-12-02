package main

import (
	"bufio"
	"fmt"
	"github.com/DimaKoz/go-socks5"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestSet(t *testing.T) {

	as := arrayFlags{}
	returnResult := as.Set("test_bad_scan.tmp")
	as.Set("test_bad_scan.tmp")
	if returnResult != nil {
		t.Errorf("return result != nil")
	}
	if "test_bad_scan.tmp, test_bad_scan.tmp" != as.String() {
		t.Errorf("#: String(%s); want %s", as.String(), "test_bad_scan.tmp, test_bad_scan.tmp")
	}
}

func TestStringEmpty(t *testing.T) {
	as := arrayFlags{}
	var testResult = as.String()
	if "[]" != testResult {
		t.Errorf("#: String(%s); want %s", testResult, "[]")
	}
}

func TestStringOneItem(t *testing.T) {
	as := arrayFlags{"abc"}
	var testResult = as.String()
	if "abc" != testResult {
		t.Errorf("#: String(%s); want %s", testResult, "abc")
	}
}

func TestStringMoreItem(t *testing.T) {
	as := arrayFlags{"abc", "abc1", "abc2"}
	var testResult = as.String()
	if "abc, abc1, abc2" != testResult {
		t.Errorf("#: String(%s); want %s", testResult, "abc, abc1, abc2")
	}
}

func TestGetArrayFlagsFromFile(t *testing.T) {
	want := "ffds:dddd, user2:pass2, user3:pass3, baduser, user4:, user5:awesome, :user4:"
	fileName := "test_auth.txt"
	as, _ := getArrayFlagsFromFile(fileName)
	//as := arrayFlags{"abc", "abc1", "abc2"}
	var testResult = as.String()
	if want != testResult {
		t.Errorf("#: String(%s); want %s", testResult, want)
	}

}

func TestFillCredentialsFromFile(t *testing.T) {
	fileName := "test_auth.txt"
	as, _ := getArrayFlagsFromFile(fileName)
	credentialsWait := make(map[string]string)
	credentialsWait["ffds"] = "dddd"
	credentialsWait["user2"] = "pass2"
	credentialsWait["user3"] = "pass3"
	credentialsWait["user5"] = "awesome"
	var credentialsTestResult map[string]string = make(map[string]string)
	fillCredentials(credentialsTestResult, &as)
	eq := reflect.DeepEqual(credentialsWait, credentialsTestResult)
	if !eq {
		t.Errorf("#: map(%s); want %s", credentialsTestResult, credentialsWait)
	}

}

func TestGetArrayFlagsNoFile(t *testing.T) {
	fileName := "no_file"
	testResult, _ := getArrayFlagsFromFile(fileName)
	if testResult != nil {
		t.Errorf("#: String(%s); want %s", testResult, "nil")
	}

}

func TestGetArrayFlagsErrorScan(t *testing.T) {
	if _, err := os.Stat("test_bad_scan.tmp"); os.IsNotExist(err) {
		input := strings.Repeat("x", bufio.MaxScanTokenSize)
		f, err := os.Create("test_bad_scan.tmp")
		if err != nil {
			panic(err)
		}
		w := bufio.NewWriter(f)
		fmt.Fprint(w, input)
		w.Flush()
	}
	as, _ := getArrayFlagsFromFile("test_bad_scan.tmp")

	if as != nil {
		t.Errorf("#: String(%s); want %s", as.String(), "nil")
	}
}

func TestFillCredentialsEmptyFlags(t *testing.T) {
	as := arrayFlags{}
	credentialsWait := make(map[string]string)
	credentialsTestResult := make(map[string]string)
	fillCredentials(credentialsTestResult, &as)
	equalsCredentials(t, credentialsTestResult, credentialsWait)
}

func TestFillCredentialsNotEmptyFlags(t *testing.T) {
	as := arrayFlags{"user1:pass1", "user:pass"}
	credentialsWait := make(map[string]string)
	credentialsWait["user"] = "pass"
	credentialsWait["user1"] = "pass1"
	credentialsTestResult := make(map[string]string)
	fillCredentials(credentialsTestResult, &as)
	equalsCredentials(t, credentialsTestResult, credentialsWait)
}

func TestFillCredentialsBadUsers(t *testing.T) {
	as := arrayFlags{"user1:pass1", "badUser", "badUser:", ":passOfBadUser"}
	credentialsWait := make(map[string]string)
	credentialsWait["user1"] = "pass1"
	credentialsTestResult := make(map[string]string)
	fillCredentials(credentialsTestResult, &as)
	equalsCredentials(t, credentialsTestResult, credentialsWait)
}

func TestFillCredentialsBadPass(t *testing.T) {
	as := arrayFlags{":passOfBadUser"}
	credentialsWait := make(map[string]string)
	credentialsTestResult := make(map[string]string)
	fillCredentials(credentialsTestResult, &as)
	equalsCredentials(t, credentialsTestResult, credentialsWait)
}

func TestCopyCredentials(t *testing.T) {
	cleanAfterInit()
	as := arrayFlags{"user1:pass1"}
	credentialsWait := make(map[string]string)
	credentialsWait["user1"] = "pass1"
	credentialsTestResult := make(map[string]string)
	fillCredentials(credentialsTestResult, &as)
	credentials = credentialsTestResult
	testCredentials := socks5.StaticCredentials{}
	copyCredentials(testCredentials)
	equalsCredentials(t, testCredentials, credentialsWait)
	cleanAfterInit()
}

func TestCopyEmptyCredentials(t *testing.T) {
	cleanAfterInit()
	credentialsWait := make(map[string]string)
	credentialsTestResult := make(map[string]string)
	credentials = credentialsTestResult
	testCredentials := socks5.StaticCredentials{}
	copyCredentials(testCredentials)
	equalsCredentials(t, testCredentials, credentialsWait)
	cleanAfterInit()
}

func equalsCredentials(t *testing.T, result map[string]string, wait map[string]string) {
	eq := reflect.DeepEqual(wait, result)
	if !eq {
		t.Errorf("#: map(%s); want %s", result, wait)
	}
}

func TestConfigGetPort(t *testing.T) {
	port = 2000
	if configGetHttpPort() != 2000 {
		t.Errorf("#: port(%d); want %d", configGetHttpPort(), port)
	}
}

func TestConfigGetSocsPort(t *testing.T) {
	socsPort = 2001
	if configGetSocsPort() != 2001 {
		t.Errorf("#: port(%d); want %d", configGetSocsPort(), port)
	}
}

func TestConfigIsUserAllowed(t *testing.T) {
	cleanAfterInit()
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	expectedUsers := []string{"eww", "222", "user3"}
	expectedPasswords := []string{"www", "3333", "pass3"}

	os.Args = []string{"", "-a=eww:www", "-a=222:3333", "-auth-file=test_auth.txt", "-" + flagNamePort + "=" + strconv.Itoa(defaultFlagPortValue), "-" + flagNameSocsPort + "=" + strconv.Itoa(defaultFlagSocsPortValue)}

	_ = initConfig()
	for i := 0; i < len(expectedUsers) && i < len(expectedPasswords); i++ {

		isAllowedUser := configIsUserAllowed(expectedUsers[i], expectedPasswords[i])
		if !isAllowedUser {
			t.Errorf("Test failed, expected user: '%s', pass:  '%s'", expectedUsers[i], expectedPasswords[i])
		}
	}
	cleanAfterInit()
	os.Args = oldArgs
	_ = initConfig()
}

func TestConfigHasUser(t *testing.T) {

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"", "-a=eww:www", "-a=222:3333", "-auth-file=test_auth.txt", "-" + flagNamePort + "=" + strconv.Itoa(defaultFlagPortValue), "-" + flagNameSocsPort + "=" + strconv.Itoa(defaultFlagSocsPortValue)}

	_ = initConfig()
	if !hasUser() {
		t.Errorf("Test failed, expected hasUser() == true")
	}
	cleanAfterInit()

}

func TestConfigCheckWrongPort(t *testing.T) {

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	initialArgs := [][]string{
		{"", "-" + flagNameSocsPort + "=65536", "-" + flagNamePort + "=" + strconv.Itoa(defaultFlagPortValue)},
		{"", "-" + flagNameSocsPort + "=-1", "-" + flagNamePort + "=" + strconv.Itoa(defaultFlagPortValue)},
		{"", "-" + flagNamePort + "=65536", "-" + flagNameSocsPort + "=" + strconv.Itoa(defaultFlagSocsPortValue)},
		{"", "-" + flagNamePort + "=-1", "-" + flagNameSocsPort + "=" + strconv.Itoa(defaultFlagSocsPortValue)},
	}

	expectedErrorMessage := []string{"TCP port must be in the range 1 - 65535,  the port[65536] is wrong\n",
		"TCP port must be in the range 1 - 65535,  the port[-1] is wrong\n",
		"TCP port must be in the range 1 - 65535,  the port[65536] is wrong\n",
		"TCP port must be in the range 1 - 65535,  the port[-1] is wrong\n",
	}

	for i := 0; i < len(initialArgs); i++ {
		os.Args = initialArgs[i]
		err := initConfig()
		if err == nil || err.Error() != expectedErrorMessage[i] {
			t.Errorf("Test failed, expected [%s], got[%s], for args[%s] ", expectedErrorMessage[i], err, initialArgs[i])
		} else {
			fmt.Print(err)
		}
	}

	cleanAfterInit()
	os.Args = []string{"", "-" + flagNamePort + "=65535", "-" + flagNameSocsPort + "=" + strconv.Itoa(defaultFlagSocsPortValue)}

	err := initConfig()
	if err != nil {
		msg := fmt.Sprintf("Test failed, error [%s]", err)
		t.Errorf(msg)
	}
	cleanAfterInit()
	os.Args = []string{"", "-" + flagNamePort + "=1", "-" + flagNameSocsPort + "=" + strconv.Itoa(defaultFlagSocsPortValue)}

	err = initConfig()
	if err != nil {
		msg := fmt.Sprintf("Test failed, error [%s]", err)
		t.Errorf(msg)
	}

	cleanAfterInit()
	os.Args = []string{"", "-" + flagNamePort + "=1", "-" + flagNameSocsPort + "=1"}
	err = initConfig()
	if err == nil {
		t.Errorf("Test failed, expected a non-nil error")
	} else {
		fmt.Print(err)
	}

}

func cleanAfterInit() {
	credentials = nil
	authFile = ""
	socsPort = 0
	port = 0
}
