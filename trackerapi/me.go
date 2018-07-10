package trackerapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	osUser "os/user"

	"github.com/volodimyr/clirescue/cmdutil"
	"github.com/volodimyr/clirescue/file"
	"github.com/volodimyr/clirescue/user"
)

const (
	headerTok = "X-TrackerToken"
	URL       = "https://www.pivotaltracker.com/services/v5/me"
)

var (
	FileLocation string     = homeDir() + "/.tracker"
	currentUser  *user.User = user.New()
	Stdout       *os.File   = os.Stdout
	token        string
)

func init() {
	data, err := file.Read(FileLocation)
	if err == nil && len(data) != 0 {
		log.Println("Token has already been submited")
		token = string(data)
	}
}

func Me() {
	req := newRequest("GET")
	if len(token) == 0 || token == "" {
		setCredentials()
		req.SetBasicAuth(currentUser.Username, currentUser.Password)
		parse(makeRequest(req))
		ioutil.WriteFile(FileLocation, []byte(currentUser.APIToken), 0644)
		return
	}

	req.Header.Set(headerTok, token)
	parse(makeRequest(req))
}

func makeRequest(req *http.Request) []byte {
	client := &http.Client{}
	resp, err := client.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err)
	}
	defer resp.Body.Close()
	fmt.Printf("\n****\nAPI response: \n%s\n", string(body))
	return body
}

func newRequest(method string) *http.Request {
	req, err := http.NewRequest(method, URL, nil)
	if err != nil {
		panic(err)
	}
	return req
}

func parse(body []byte) {
	var meResp = new(user.User)
	err := json.Unmarshal(body, &meResp)
	if err != nil {
		fmt.Println("error:", err)
	}

	currentUser.APIToken = meResp.APIToken
}

func setCredentials() {
	fmt.Fprint(Stdout, "Username: ")
	var username = cmdutil.ReadLine()
	cmdutil.Silence()
	fmt.Fprint(Stdout, "Password: ")

	var password = cmdutil.ReadLine()
	currentUser.Login(username, password)
	cmdutil.Unsilence()
}

func homeDir() string {
	usr, err := osUser.Current()
	if err != nil {
		panic(err)
	}
	return usr.HomeDir
}
