package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v4"
)

type Image struct {
	Title   string `json:"title"`
	AltText string `json:"alt_text"`
	URL     string `json:"url"`
}

func main() {

	dbUrl, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		fmt.Fprint(os.Stderr, "DATABASE_URL needs to be set")
		os.Exit(1)
	}

	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to establish connection: %v", err)
		os.Exit(1)
	}
	fmt.Printf("Db connected ...\n")
	defer conn.Close(context.Background())

	http.HandleFunc("/images.json", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			params := r.URL.Query().Get("indent")
			indent, err := handleIndent(params)
			if err != nil {
				fmt.Fprint(os.Stderr, err)
				os.Exit(1)
			}
			
			images, err := fetchImages(conn)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Cannot fetch images: %v", err)
				os.Exit(1)
			}
			output, err := json.MarshalIndent(images, "", indent)
			if err != nil {
				fmt.Fprint(os.Stderr, err)
				os.Exit(1)
			}
			w.Header().Set("Content-Type", "application/json")

			w.Write(output)

		}

		if r.Method == "POST" {
			var result Image
			reqBody, err := io.ReadAll(r.Body)

			if err != nil {
				fmt.Fprint(os.Stderr, err)
				os.Exit(1)
			}

			err = json.Unmarshal(reqBody, &result)

			if err != nil {
				fmt.Fprint(os.Stderr, err)
				os.Exit(1)
			}

			err = saveImages(conn, result)

			if err != nil {
				fmt.Fprint(os.Stderr, err)
				os.Exit(1)
			}
			resp, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				fmt.Fprint(os.Stderr, err)
				os.Exit(1)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(resp)
		}
	})

	http.ListenAndServe(":8080", nil)
}

func fetchImages(conn *pgx.Conn) ([]Image, error) {
	rows, err := conn.Query(context.Background(), "SELECT title, url, alt_text FROM public.images")

	if err != nil {
		return nil, err
	}
	var title, url, altText string
	var images []Image

	for rows.Next() {
		rows.Scan(&title, &url, &altText)
	}
	images = append(images, Image{Title: title, URL: url, AltText: altText})

	return images, nil
}

func saveImages(conn *pgx.Conn, data Image) error {

	queryString := fmt.Sprintf("INSERT INTO public.images(title, url, alt_text) VALUES ('%v', '%v', '%v')", data.Title, data.URL, data.AltText)

	_, err := conn.Query(context.Background(), queryString)

	return err
}

func handleIndent(v string) (string, error) {
	space := ""

	val, err := strconv.Atoi(v)
	if val > 10 || err != nil {
		return "", errors.New("indent number must be from 0 - 10")
	}

	for i := 0; i < val; i++ {
		space += " "
	}
	return space, nil
}
