package trackerapi

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	u "os/user"

	"github.com/rhymefororange/clirescue/cmdutil"
	"github.com/rhymefororange/clirescue/user"
)

var (
	url          = "https://www.pivotaltracker.com/services/v5/me"
	fileLocation = homeDir() + "/.tracker"
	currentUser  = user.New()
	stdout       = os.Stdout
)

// Me performs a request to Pivotal Tracker API,
// gets response in JSON format with information about the account and
// writes API token to file, determined by fileLocation.
func Me() {
	setCredentials()
	parse(makeRequest())
	ioutil.WriteFile(fileLocation, []byte(currentUser.APIToken), 0644)
}

func makeRequest() []byte {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if len(currentUser.APIToken) > 0 {
		fmt.Println("Using token")
		req.Header.Set("X-TrackerToken", currentUser.APIToken)
	} else {
		fmt.Println("\nUsing user/pass pair")
		req.SetBasicAuth(currentUser.Username, currentUser.Password)
	}
	resp, err := client.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("\n****\nAPI response: \n%s\n", string(body))
	return body
}

func parse(body []byte) {
	var meResp = new(meResponse)
	err := json.Unmarshal(body, &meResp)
	if err != nil {
		fmt.Println("error:", err)
	}

	currentUser.APIToken = meResp.APIToken
}

func setCredentials() {
	var username, password, token string
	if _, err := os.Stat(fileLocation); err == nil {
		file, err := os.Open(fileLocation)
		if err != nil {
			log.Println(err)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		scanner.Scan()
		text := scanner.Text()
		if len(text) > 0 {
			token = text
		}
	}
	if len(token) < 1 {
		fmt.Fprint(stdout, "Username: ")
		username = cmdutil.ReadLine()
		cmdutil.Silence()
		fmt.Fprint(stdout, "Password: ")
		password = cmdutil.ReadLine()
		cmdutil.Unsilence()
	}
	currentUser.Login(username, password, token)
}

func homeDir() string {
	usr, _ := u.Current()
	return usr.HomeDir
}

type meResponse struct {
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
