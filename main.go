package main

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type githubAPI struct {
	owner      string
	repository string
}

func (g githubAPI) getRepositoryChanges() []byte {
	var changes []map[string]interface{}
	var checksum string
	var checksumBytes []byte
	hash := sha256.New()
	res, err := http.Get("https://api.github.com/repos/" + g.owner + "/" + g.repository + "/events")
	if res.StatusCode == 403 {
		panic("API rate limit exceeded")
	}
	checkErr(err)
	jsonRes, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(jsonRes))
	checkErr(err)
	err = json.Unmarshal(jsonRes, &changes)
	checkErr(err)
	for _, e := range changes {
		if (e["type"]) == "PushEvent" {
			checksum += e["id"].(string)
		}
	}
	checksumBytes = []byte(checksum)
	hash.Write(checksumBytes)
	return hash.Sum(nil)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var owner string
	var repository string
	var timeout int

	flag.StringVar(&owner, "owner", "axotion", "Provide owner")
	flag.StringVar(&repository, "repository", "github-repo-notify", "Provide repository")
	flag.IntVar(&timeout, "timeout", 60, "Provide timeout")
	flag.Parse()
	githubAPI := &githubAPI{owner: owner, repository: repository}
	tmpChecksum := fmt.Sprintf("%x \n", githubAPI.getRepositoryChanges())
	for {
		if tmpChecksum != fmt.Sprintf("%x \n", githubAPI.getRepositoryChanges()) {
			tmpChecksum = fmt.Sprintf("%x \n", githubAPI.getRepositoryChanges())
			fmt.Println("Change")
		}
		time.Sleep(time.Second * time.Duration(timeout))
		fmt.Println("No changes")
	}
}
