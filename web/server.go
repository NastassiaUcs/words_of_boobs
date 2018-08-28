package web

import (
	"net/http"
	"fmt"
	"./../generator"
	"strconv"
	"log"
	"os"
	"unicode/utf8"
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
			width int = 0
			err error
		)
		params := r.URL.Query()
		if values, ok = params["text"]; !ok {
			w.Write([]byte("error: no text parameter"))
			return
		} else {
			text = string(values[0])
		}

		if values, ok = params["width"]; ok {
			if width, err = strconv.Atoi(values[0]); err != nil {
				log.Println(err)
				w.Write([]byte("error: incorrect width value " + values[0]))
				return
			}
		}

		if width > 10000 {
			w.Write([]byte("error: maximum width = 10000"))
			return
		}

		if utf8.RuneCountInString(text) > 15 {
			w.Write([]byte("error: maximum text length = 15 symbols"))
			return
		}

		filename, err = generator.GenerateImageForText(text, "Symbola.ttf", "boobs", 1000, width)
		if err != nil {
			log.Println(err)
			w.Write([]byte("error: something wrong"))
			return
		}

		writeLog(text)
		log.Printf("generated %s for text '%s' with width=%d\n", filename, text, width	)
		w.Write([]byte(filename))
	})

	log.Printf("started on %d port\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
