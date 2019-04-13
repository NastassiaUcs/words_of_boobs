package web

import (
	"net/http"
	"fmt"
	"./../generator"
	"log"
	"os"
	"unicode/utf8"
)

const (
	TEXT_LENGTH_LIMIT = 15
)

func writeLog(value string) {
	f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	_, err = f.WriteString(value + "\n")
	if err != nil {
		log.Println(err)
	}
	f.Close()
}


func Start(port int) error {
	http.Handle("/", http.FileServer(http.Dir("./html")))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./results"))))
	http.HandleFunc("/api/generate", func(w http.ResponseWriter, r *http.Request) {
		var (
			values []string
			filename, text string
			ok bool
			err error
		)
		params := r.URL.Query()
		if values, ok = params["text"]; !ok {
			w.Write([]byte("error: no text parameter"))
			return
		} else {
			text = string(values[0])
		}

		if utf8.RuneCountInString(text) > TEXT_LENGTH_LIMIT {
			w.Write([]byte(fmt.Sprintf("error: maximum text length = %d symbols", TEXT_LENGTH_LIMIT)))
			return
		}

		filename, err = generator.GenerateImageForText(text,  "boobs")
		if err != nil {
			log.Println(err)
			w.Write([]byte("error: something wrong"))
			return
		}

		writeLog(text)
		log.Printf("generated %s for text '%s'\n", filename, text	)
		w.Write([]byte(filename))
	})

	log.Printf("started on %d port\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
