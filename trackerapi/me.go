package trackerapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	u "os/user"
	"strings"

	"github.com/iyuroch/clirescue/cmdutil"
	"github.com/iyuroch/clirescue/user"
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
	//we parse the body of response ->
	//if not valid authentication delete stored credentials
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

func loadUser() error {
	//we have to load them from our storage file,
	//if we encounter any error on the way - return error
	configFile, err := os.Open(FileLocation)
	if err != nil {
		return err
	}
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(currentUser); err != nil {
		return err
	}

	if currentUser.Password == "" && currentUser.Username == "" {
		return errors.New("no username and password present")
	}
	return nil
}

func setCredentials() {
	//try to get them from storage -> if not present, prompt
	err := loadUser()
	if err != nil {
		promptUser()
	}
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
