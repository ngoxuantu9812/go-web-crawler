package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/net/html"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

const keyServerAddr = "serverAddr"

func main() {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/crawler", fetch)
	ctx := context.Background()
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
		BaseContext: func(listener net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, listener.Addr().String())
			return ctx
		},
	}

	err = server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error listening for server: %s\n", err)
	}
}

func fetch(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://dantri.com.vn/")
	if err != nil {
		//return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		//return "", fmt.Errorf("error: failed to fetch URL %s: %s", url, resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		//return "", err
	}
	doc, err := html.Parse(strings.NewReader(string(bodyBytes)))
	if err != nil {
	}

	var links []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					links = append(links, a.Val)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	fmt.Println(links)
}
