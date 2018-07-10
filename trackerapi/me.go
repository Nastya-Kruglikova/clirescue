package trackerapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	u "os/user"

	"github.com/GoBootcamp/clirescue/cmdutil"
	"github.com/GoBootcamp/clirescue/user"
)

const (
	url      = "https://www.pivotaltracker.com/services/v5/me"
	fileName = "tracker"
)

var (
	//currentUser  *user.User
	stdout *os.File
)

func init() {
	stdout = os.Stdout
}

// Me reads the username an password and makes a request to get the API token. If the request is successful stores
// the obtained API token to a file
func Me() {

	home, err := homeDir()
	if err != nil {
		fmt.Println(err)
		return
	}
	fileLocation := home + "/" + fileName

	currentUser, err := cmdutil.Credentials()
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := makeRequest(currentUser)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = parse(currentUser, resp); err != nil {
		fmt.Println(err)
		return
	}

	if err = ioutil.WriteFile(fileLocation, []byte(currentUser.APIToken), 0644); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("file %s was created in %s, API token: %s", fileName, fileLocation, currentUser.APIToken)
}

func makeRequest(currentUser *user.User) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(currentUser.Username, currentUser.Password)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("could not get API token, response status " + resp.Status)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Printf("\n****\nAPI response: \n%s\n", string(body))
	return body, nil
}

func parse(currentUser *user.User, body []byte) error {
	resp := make(map[string]interface{})

	err := json.Unmarshal(body, &resp)
	if err != nil {
		return err
	}

	token, ok := resp["api_token"]
	if ok {
		apiToken, ok := token.(string)
		if ok && apiToken != "" {
			currentUser.APIToken = apiToken
			return nil
		}
	}
	return errors.New("could not parse API token")
}

func homeDir() (string, error) {
	usr, err := u.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}
