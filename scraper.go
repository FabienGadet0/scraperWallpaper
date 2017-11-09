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
			// Closed        int    `json:"closed"`
			Com string `json:"com"`
			Ext string `json:"ext"`
			// Filename      string `json:"filename"`
			// Fsize         int    `json:"fsize"`
			// H             int    `json:"h"`
			// Images        int    `json:"images"`
			// Md5           string `json:"md5"`
			// Name          string `json:"name"`
			No int `json:"no"`
			// Now           string `json:"now"`
			// OmittedImages int    `json:"omitted_images"`
			// OmittedPosts  int    `json:"omitted_posts"`
			// Replies       int    `json:"replies"`
			// Resto int `json:"resto"`
			// SemanticURL   string `json:"semantic_url"`
			// Sticky        int    `json:"sticky"`
			Sub string `json:"sub"`
			Tim int    `json:"tim"`
			// Time          int    `json:"time"`
			// TnH           int    `json:"tn_h"`
			// TnW           int    `json:"tn_w"`
			// W             int    `json:"w"`
		} `json:"posts"`
	} `json:"threads"`
}

var searchedWords = []string{"minimalistic", "Minimalistic", "depressing", "Depressing", "minimal", "Minimal", "Depressing", "Paint", "paint",
	"Aesthetic", "aesthetic"}
var board = "wg"
var url = "https://boards.4chan.org/" + board
var downloadURL = "https://i.4cdn.org/wg/"
var total = 0

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
		println(err)
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
		// println(err)
		defer wg.Done()
		return false
	}
	os.Mkdir("./imgs/", 0777)
	file, errr := os.Create("./imgs/" + name)
	if errr != nil {
		// println(err)
		defer wg.Done()
		return false
	}
	_, errrr := io.Copy(file, response.Body)
	file.Close()
	if errrr != nil {
		// println(err)
		defer wg.Done()
		return false
	}
	// fmt.Println("new file number : " + name + " saved")
	total++
	defer wg.Done()
	return true
}

func parsewgfiltered(data MyJsonName, hardcore bool) {
	for _, thrd := range data.Threads {
		if stringInSlice(thrd.Posts[0].Sub, searchedWords) {
			if hardcore {
				HARDCOREMODE(thrd.Posts[0].No)
			} else {
				// for _, pst := range thrd.Posts {
				// 	if pst.Tim != 0 {
				// 		name := strconv.Itoa(pst.Tim) + pst.Ext
				// 		saveFile(name, downloadURL+namem, nil)
				// 	}
				// }
			}
		}
	}
}

func HARDCOREMODE(resto int) {
	var wg sync.WaitGroup
	response := request(url + "/thread/" + strconv.Itoa(resto))
	r := regexp.MustCompile("(i\\.4cdn\\.org\\/wg\\/)(\\d+\\.(?:jpg|png))")
	rs := r.FindAllString(string(response), -1)
	for i, link := range rs {
		wg.Add(1)
		go saveFile(strconv.Itoa(i), link, &wg)
	}
	wg.Wait()
}

func worker(url string, wg *sync.WaitGroup, hardcore bool) {
	defer wg.Done()
	var a MyJsonName
	data := request(url)
	json.Unmarshal([]byte(data), &a)
	parsewgfiltered(a, hardcore)
}

func main() {
	start := time.Now()
	var wg sync.WaitGroup
	for i := 1; i < 10; i++ {
		wg.Add(1)
		go worker("https://a.4cdn.org/"+board+"/"+strconv.Itoa(i)+".json", &wg, true)
	}
	fmt.Println("ON ATTENDS LA FIN DES WORKERS")
	wg.Wait()
	fmt.Println("C EST DEJA FINI OMG")
	fmt.Println(time.Since(start))
	fmt.Println("total de", total, "fichiers")
}