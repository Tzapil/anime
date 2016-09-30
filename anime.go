package main

import (
    "log"
    "bytes"
    "io"
    "io/ioutil"
    "net/http"

    "golang.org/x/net/html"

    e "github.com/tzapil/anime/entry"

    "golang.org/x/text/encoding/charmap"
    "golang.org/x/text/transform"

    "encoding/json"
)

const site_adress = "http://www.fansubs.ru/"

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

func ParseSite(id string) *e.Entry {
    log.Printf("Try to get page with id: %s\n", id)

    answer, err_get := http.Get(site_adress + "base.php?id=" + id)
    if err_get != nil {
        log.Printf("Error caught while taking page with id: %s\n%s\n", id, err_get.Error())
        return nil
    }

    defer answer.Body.Close()

    body, err_read := ioutil.ReadAll(answer.Body)
    if err_read != nil {
        log.Printf("Error caught while reading body with id: %s\n%s\n", id, err_read.Error())
        return nil
    }

    log.Printf("Try to parse page with id: %s\n", id)
    z, err_parse := html.Parse(bytes.NewBuffer(Win1251ToUtf8(body)))

    if err_parse != nil {
        log.Printf("Error caught while parse body with id: %s\n%s\n", id, err_parse.Error())
        return nil
    }
  
    log.Printf("Retrun information from page with id: %s\n", id)
    return e.Parse(z)
}

func main() {
    // Create handlers
    handler := func (w http.ResponseWriter, r *http.Request) {
        anime_id := r.FormValue("id")
        log.Printf("Request for anime with id %s\n", anime_id)
        if anime_id == "" {
            http.Error(w, "Anime id parameter required!", http.StatusBadRequest)
            return
        }

        res := ParseSite(anime_id)

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
    http.HandleFunc("/anime", handler)

    // Start web server
    log.Println("About to listen on 8080. Go to http://127.0.0.1:8023/anime?id=123")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
