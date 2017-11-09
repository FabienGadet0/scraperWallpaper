package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type MyJsonName struct {
	Threads []struct {
		Posts []struct {
			Com string `json:"com"`
			Ext string `json:"ext"`
			No int `json:"no"`
			Sub string `json:"sub"`
			Tim int    `json:"tim"`
		} `json:"posts"`
	} `json:"threads"`
}

var searchedWords = []string{"minimalistic", "Minimalistic", "depressing", "Depressing", "minimal", "Minimal", "Depressing", "Paint", "paint",
	"Aesthetic", "aesthetic", "winter", "Winter"}

var board = "wg"
var url = "https://boards.4chan.org/" + board
var downloadURL = "https://i.4cdn.org/wg/"
var total = 0
var fileNumber = 0 // yes its different from total

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if strings.Contains(a, b) {
			return true
		}
	}
	return false
}

func request(url string) []byte {
	response, err := http.Get(url)
	if err != nil {
		println("REQUEST FAILLED , check your connexion or if you can access to ", url)
	} else {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		return body
	}
	return nil
}

func saveFile(name string, url string, wg *sync.WaitGroup) bool {
	response, err := http.Get("https://" + url)
	if err != nil {
		defer wg.Done()
		fileNumber--
		return false
	}
	os.Mkdir("./imgs/", 0777)
	file, errr := os.Create("./imgs/" + name)
	if errr != nil {
		defer wg.Done()
		defer file.Close()
		fileNumber--
		return false
	}
	_, errrr := io.Copy(file, response.Body)
	file.Close()
	if errrr != nil {
		defer wg.Done()
		os.Remove(name)
		fmt.Println(name , " deleted because it's incomplete")
		return false
	}
	defer wg.Done()
	return true
}

func parsewgfiltered(data MyJsonName, hardcore bool) {
	for _, thrd := range data.Threads {
		if stringInSlice(thrd.Posts[0].Sub, searchedWords) {
				HARDCOREMODE(thrd.Posts[0].No)
			}
		}
	}

func HARDCOREMODE(resto int) {
	var wg sync.WaitGroup
	response := request(url + "/thread/" + strconv.Itoa(resto))
	r := regexp.MustCompile("(i\\.4cdn\\.org\\/wg\\/)(\\d+\\.(?:jpg|png))")
	rs := r.FindAllString(string(response), -1)
	for _, link := range rs {
		if (total % 2 == 0){ // request sent for every pic a copy of it just after
		wg.Add(1)
		go saveFile(strconv.Itoa(fileNumber), link, &wg)
		fileNumber++
		}
		total++
	}
	wg.Wait()
}

func worker(url string, wg *sync.WaitGroup, hardcore bool) {
	var a MyJsonName
	data := request(url)
	json.Unmarshal([]byte(data), &a)
	parsewgfiltered(a, hardcore)
	defer wg.Done()
}

func main() {
	start := time.Now()
	var wg sync.WaitGroup
	fmt.Println("waiting for workers ....")
	for i := 1; i < 5; i++ {
		wg.Add(1)
		worker("https://a.4cdn.org/"+board+"/"+strconv.Itoa(i)+".json", &wg, true)
	}
	wg.Wait()
	fmt.Println("Finished ...")
	fmt.Println(time.Since(start))
	fmt.Println(fileNumber, "file downloaded")
}