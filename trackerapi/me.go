package trackerapi

import (
	"encoding/json"
	"fmt"
	"strings"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	u "os/user"

	"github.com/GoBootcamp/clirescue/cmdutil"
	"github.com/GoBootcamp/clirescue/user"
)

var (
	URL          string     = "https://www.pivotaltracker.com/services/v5/me"
	FileLocation string     = homeDir() + "/.tracker"
	currentUser  *user.User = user.New()
	Stdout       *os.File   = os.Stdout
)

func Me() {
	setCredentials()
	err := parse(makeRequest())
	if err != nil {
		err := os.Remove(FileLocation)
		if err != nil {
			fmt.Println("err: ", err)
		}
	} else {
		data, _ := json.Marshal(currentUser)
		ioutil.WriteFile(FileLocation, data, 0644)
	}
}

func makeRequest() []byte {
	client := &http.Client{}
	req, err := http.NewRequest("GET", URL, nil)
	req.SetBasicAuth(currentUser.Username, currentUser.Password)
	resp, err := client.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("\n****\nAPI response: \n%s\n", string(body))
	return body
}

func parse(body []byte) error {
	var meResp = new(MeResponse)
	err := json.Unmarshal(body, &meResp)
	if err != nil {
		fmt.Println("error:", err)
	}

	if strings.Contains(string(body), "invalid_authentication") {
		//we stored username and password, have to delete them
		return errors.New("not valid credentials")
	}

	currentUser.APIToken = meResp.APIToken
	return nil
}

func promptUser() {
	fmt.Fprint(Stdout, "Username: ")
	var username = cmdutil.ReadLine()
	cmdutil.Silence()
	fmt.Fprint(Stdout, "Password: ")

	var password = cmdutil.ReadLine()
	currentUser.Login(username, password)
	cmdutil.Unsilence()
	return
}

func setCredentials() {
	//try to get them from storage -> if not present, prompt
	configFile, err := os.Open(FileLocation)
    if err != nil {
		promptUser()
		return
    }
    jsonParser := json.NewDecoder(configFile)
    if err = jsonParser.Decode(currentUser); err != nil {
		promptUser()
		return
	}

	if currentUser.Password != "" && currentUser.Username != "" {
		return
	}
	promptUser()
}

func homeDir() string {
	usr, _ := u.Current()
	return usr.HomeDir
}

type MeResponse struct {
	APIToken string `json:"api_token"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Initials string `json:"initials"`
	Timezone struct {
		Kind      string `json:"kind"`
		Offset    string `json:"offset"`
		OlsonName string `json:"olson_name"`
	} `json:"time_zone"`
}
