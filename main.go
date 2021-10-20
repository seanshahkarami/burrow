package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
)

func convertGopherToHTML(dst io.Writer, src io.Reader) error {
	scanner := NewScanner(src)

	fmt.Fprintf(dst, "<html><body><pre>\n")

	for scanner.Scan() {
		switch scanner.Code() {
		case "i":
			fmt.Fprintf(dst, "%s\n", scanner.Field(0))
		case "0", "1", "p", "g", "I":
			suffix := ""

			// add / for dir
			if scanner.Code() == "1" {
				suffix = "/"
			}

			icon := ""
			if strings.ContainsAny(scanner.Code(), "pgI") {
				icon = "&#128247; " // camera
			}

			name := scanner.Field(0)
			path := fmt.Sprintf("/%s:%s%s%s", scanner.Field(2), scanner.Field(3), scanner.Field(1), suffix)
			fmt.Fprintf(dst, "<p>%s<a href=\"%s\">%s</a></p>\n", icon, path, name)
		case "h":
			name := scanner.Field(0)
			url := strings.TrimPrefix(scanner.Field(1), "URL:")
			fmt.Fprintf(dst, "<p><a href=\"%s\">%s</a></p>\n", url, name)
		default:
			log.Printf("unknown code=%v fields=%v\n", scanner.Code(), scanner.Fields())
		}
	}

	fmt.Fprintf(dst, "</pre></body></html>\n")

	if scanner.Err() != nil {
		return scanner.Err()
	}

	return nil
}

func openGopher(URL *url.URL) (io.ReadCloser, error) {
	conn, err := net.Dial("tcp", URL.Host)
	if err != nil {
		return nil, err
	}

	if _, err := fmt.Fprintf(conn, "%s\r\n", URL.Path); err != nil {
		conn.Close()
		return nil, err
	}

	return conn, nil
}

func main() {
	addr := flag.String("addr", ":7070", "address to listen on")
	flag.Parse()

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "error: file not found", http.StatusNotFound)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			fmt.Fprintf(w, "<html><body><p>Please start with one of the following sites:</p><ul><li><a href=\"/gopher.floodgap.com:70/\">gopher.floodgap.com</a></li></ul></body></html>\n")
			return
		}

		gopherURL, err := url.Parse("gopher://" + strings.TrimPrefix(r.URL.Path, "/"))
		if err != nil {
			return
		}

		resp, err := openGopher(gopherURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		defer resp.Close()

		// we stream files back as-is (maybe fancier later?)
		if !strings.HasSuffix(gopherURL.Path, "/") {
			log.Printf("serving file %v", r.URL)
			io.Copy(w, resp)
			return
		}

		if err := convertGopherToHTML(w, resp); err != nil {
			log.Printf("convertGopherToHTML error: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	log.Printf("listening on %s", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}
