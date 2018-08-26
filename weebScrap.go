package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

type requestStruct struct {
	Threads []struct {
		Posts []struct {
			Com string `json:"com"`
			Ext string `json:"ext"`
			No  int    `json:"no"`
			Sub string `json:"sub"`
			Tim int    `json:"tim"`
		} `json:"posts"`
	} `json:"threads"`
}

var searchedWords = []string{"minimalistic", "Minimalistic", "Minimalism", "minimalism", "depressing", "Depressing", "minimal", "Minimal", "Depressing",
	"Aesthetic", "aesthetic", "winter", "Winter", "Philosophical", "Art", "art", "sad", "Sad", "Depressed", "depressed", "vaporwave", "Vaporwave", "current wallpaper", "Current wallpaper",
	"dark papes", "Dark papes", "90's"}

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

func saveFile(name string, url string, wg *sync.WaitGroup, folderName string) bool {
	response, err := http.Get("https://" + url)
	ext := filepath.Ext(url)
	if err != nil {
		defer wg.Done()
		fileNumber--
		return false
	}
	os.Mkdir("./imgs/", 0777)
	os.Mkdir("./imgs/"+folderName, 0777)
	file, errr := os.Create("./imgs/" + folderName + "/" + name + ext)
	if errr != nil {
		defer file.Close()
		fileNumber--
		defer wg.Done()
		return false
	}
	defer response.Body.Close()
	_, errrr := io.Copy(file, response.Body)
	if errrr != nil {
		defer wg.Done()
		os.Remove(name)
		color.Red(name, " deleted because it's incomplete")
		return false
	}
	defer file.Close()
	defer wg.Done()
	return true
}

func parsewgfiltered(data requestStruct, hardcore bool) {
	for _, thrd := range data.Threads {
		if stringInSlice(thrd.Posts[0].Sub, searchedWords) {
			color.White("------------------------------------")
			color.White("Thread found : " + thrd.Posts[0].Sub)
			rawMode(thrd.Posts[0].No, thrd.Posts[0].Sub)
		}
	}
}

func rawMode(resto int, folderName string) {
	var wg sync.WaitGroup
	response := request(url + "/thread/" + strconv.Itoa(resto))
	r := regexp.MustCompile("(i.4cdn.org\\/wg\\/)(\\w+.(?:jpg|png))")
	rs := r.FindAllString(string(response), -1)
	color.White("Starting download")
	color.White("------------------------------------")
	for _, link := range rs {
		wg.Add(1)
		link = strings.Replace(link, "s", "", -1)
		go saveFile(strconv.Itoa(fileNumber), link, &wg, folderName)
		fileNumber++
	}
	total++
	wg.Wait()
	color.Green(folderName + " finished")
}

func worker(url string, wg *sync.WaitGroup, hardcore bool) {
	var a requestStruct
	data := request(url)
	json.Unmarshal([]byte(data), &a)
	parsewgfiltered(a, hardcore)
	defer wg.Done()
}

func main() {
	start := time.Now()
	var wg sync.WaitGroup
	for i := 1; i < 10; i++ {
		wg.Add(1)
		worker("https://a.4cdn.org/"+board+"/"+strconv.Itoa(i)+".json", &wg, true)
	}
	wg.Wait()
	if total == 0 {
		color.Red("no keywords match")
		fmt.Println("current keywords are : " + strings.Join(searchedWords, " "))
	} else {
		color.Green("Finished ...")
		color.Green("Total Time : ")
		color.Green("%s", time.Since(start))
		color.Green("%s file downloaded", fileNumber)
		color.Green("all the downloaded files are in the imgs folder")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}

// Fabien Gadet
