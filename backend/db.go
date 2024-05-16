package main

import (
	"bytes"
	"database/sql"
	"image"
	"image/png"
	_ "image/png"
	"io"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/nfnt/resize"
)

func insertImageFromUrl(url string, prompt string) string {
	response, e := http.Get(url)
	if e != nil {
		log.Fatal(e)
	}
	defer response.Body.Close()

	imgData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Something went wrong!")
	}

	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		log.Printf("Something went wrong!")
	}

	imgLowRes := resize.Resize(256, 256, img, resize.Lanczos3)
	imgLowResBuf := new(bytes.Buffer)
	err = png.Encode(imgLowResBuf, imgLowRes)
	if err != nil {
		log.Printf("Something went wrong!")
	}

	var id string
	if err := db.QueryRow(
		"INSERT INTO images(prompt, url, image, imageLowRes) VALUES($1, $2, $3, $4) RETURNING id",
		prompt, url, imgData, imgLowResBuf.Bytes(),
	).Scan(&id); err != nil {
		panic(err)
	}

	return id
}

func selectImageById(id string, isHighResolution bool) *Image {
	image := &Image{}
	var query string
	if isHighResolution {
		query = "SELECT prompt, url, image FROM images WHERE id=$1"
	} else {
		query = "SELECT prompt, url, imageLowRes FROM images WHERE id=$1"
	}

	if err := db.QueryRow(query, id).Scan(&image.Prompt, &image.URL, &image.Data); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return nil
	}
	image.ID = id
	return image
}

func selectImageIds() []string {
	query := "SELECT id FROM images ORDER BY timestamp DESC"

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			log.Fatal(err)
			return nil
		}
		ids = append(ids, id)
	}
	if err = rows.Err(); err != nil {
		return nil
	}

	return ids
}
