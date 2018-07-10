package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	u "os/user"
	"path"
)

const url = "https://www.pivotaltracker.com/services/v5/me"

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

func me() {
	username, pass, err := readCredentials()
	if err != nil {
		log.Fatal(err)
	}

	data, err := makeRequest(url, username, pass)
	if err != nil {
		log.Fatal(err)
	}

	me, err := parse(data)
	if err != nil {
		log.Fatal(err)
	}

	err = saveToken(me.APIToken)
	if err != nil {
		log.Fatal(err)
	}
}

func makeRequest(url, username, password string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req.SetBasicAuth(username, password)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Fprintln(os.Stdout, string(body))

	return body, nil
}

func readCredentials() (username, password string, e error) {
	fmt.Fprint(os.Stdout, "Username: ")
	u, err := readLine()
	if err != nil {
		return "", "", err
	}

	silence()
	fmt.Fprint(os.Stdout, "Password: ")

	p, err := readLine()
	if err != nil {
		return "", "", err
	}
	unsilence()
	fmt.Fprintln(os.Stdout, "")

	return u, p, nil
}

func parse(body []byte) (*meResponse, error) {
	r := &meResponse{}
	err := json.Unmarshal(body, r)
	if err != nil {
		return nil, err
	}

	fmt.Fprintln(os.Stdout, r)
	return r, nil
}

func saveToken(t string) error {
	usr, err := u.Current()
	if err != nil {
		return fmt.Errorf("cannot retrieve directory for token: %v", err)
	}
	fileLocation := path.Join(usr.HomeDir, ".tracker")
	ioutil.WriteFile(fileLocation, []byte(t), 0644)
	return nil
}
