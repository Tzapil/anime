package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	a "github.com/tzapil/anime/entries"

	"golang.org/x/net/html"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"

	"encoding/json"
)

const siteAddress = "http://www.fansubs.ru/"

func Win1251ToUtf8(s []byte) []byte {
	result := []byte{}

	sr := bytes.NewReader(s)
	tr := transform.NewReader(sr, charmap.Windows1251.NewDecoder())
	buf, err := ioutil.ReadAll(tr)
	if err != err {
		// обработка ошибки
	}

	result = buf // строка в UTF-8

	return result
}

func ParseAuPage(id string) *a.Author {
	log.Printf("Try to get page with id: %s\n", id)

	answer, errGet := http.Get(siteAddress + "base.php?au=" + id)
	if errGet != nil {
		log.Printf("Error caught while taking page with id: %s\n%s\n", id, errGet.Error())
		return nil
	}

	defer answer.Body.Close()

	body, errRead := ioutil.ReadAll(answer.Body)
	if errRead != nil {
		log.Printf("Error caught while reading body with id: %s\n%s\n", id, errRead.Error())
		return nil
	}

	log.Printf("Try to parse page with id: %s\n", id)
	z, errParse := html.Parse(bytes.NewBuffer(Win1251ToUtf8(body)))

	if errParse != nil {
		log.Printf("Error caught while parse body with id: %s\n%s\n", id, errParse.Error())
		return nil
	}

	log.Printf("Return information author with id: %s\n", id)
	return a.ParseAuthorPage(z)
}

func ParseAnimePage(id string) *a.Entry {
	log.Printf("Try to get author with id: %s\n", id)

	answer, errGet := http.Get(siteAddress + "base.php?id=" + id)
	if errGet != nil {
		log.Printf("Error caught while taking author with id: %s\n%s\n", id, errGet.Error())
		return nil
	}

	defer answer.Body.Close()

	body, errRead := ioutil.ReadAll(answer.Body)
	if errRead != nil {
		log.Printf("Error caught while reading author with id: %s\n%s\n", id, errRead.Error())
		return nil
	}

	log.Printf("Try to parse page with id: %s\n", id)
	z, errParse := html.Parse(bytes.NewBuffer(Win1251ToUtf8(body)))

	if errParse != nil {
		log.Printf("Error caught while parse author with id: %s\n%s\n", id, errParse.Error())
		return nil
	}

	log.Printf("Return information from author with id: %s\n", id)
	return a.ParseAnime(z)
}

func animeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	animeID := vars["id"]
	log.Printf("Request for anime with id %s\n", animeID)
	if animeID == "" {
		http.Error(w, "Anime id required!", http.StatusBadRequest)
		return
	}

	res := ParseAnimePage(animeID)

	if res == nil {
		http.Error(w, "Anime not found!", http.StatusNotFound)
		return
	}

	j, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	io.WriteString(w, string(j))
}

func auHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	animeID := vars["id"]
	log.Printf("Request for author with id %s\n", animeID)
	if animeID == "" {
		http.Error(w, "Author id required!", http.StatusBadRequest)
		return
	}

	res := ParseAuPage(animeID)

	if res == nil {
		http.Error(w, "Author not found!", http.StatusNotFound)
		return
	}

	j, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	io.WriteString(w, string(j))
}

func main() {
	// build router
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/anime/{id}", animeHandler)
	myRouter.HandleFunc("/au/{id}", auHandler)

	// Start web server
	log.Println("About to listen on 8080. Go to http://127.0.0.1:8080/anime/123")
	err := http.ListenAndServe(":8080", myRouter)
	log.Fatal(err)
}
