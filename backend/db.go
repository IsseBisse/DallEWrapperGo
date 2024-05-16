package main

import (
	"bytes"
	"image"
	"image/png"
	_ "image/png"
	"io"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/nfnt/resize"
)

func insertImageFromUrl(url string, prompt string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	imgData, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return "", err
	}

	imgLowRes := resize.Resize(256, 256, img, resize.Lanczos3)
	imgLowResBuf := new(bytes.Buffer)
	if err := png.Encode(imgLowResBuf, imgLowRes); err != nil {
		return "", err
	}

	var id string
	if err := db.QueryRow(
		"INSERT INTO images(prompt, url, image, imageLowRes) VALUES($1, $2, $3, $4) RETURNING id",
		prompt, url, imgData, imgLowResBuf.Bytes(),
	).Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}

func selectImageById(id string, isHighResolution bool) (*Image, error) {
	image := &Image{}
	var query string
	if isHighResolution {
		query = "SELECT prompt, url, image FROM images WHERE id=$1"
	} else {
		query = "SELECT prompt, url, imageLowRes FROM images WHERE id=$1"
	}

	if err := db.QueryRow(query, id).Scan(&image.Prompt, &image.URL, &image.Data); err != nil {
		return nil, err
	}
	image.ID = id
	return image, nil
}

func selectImageIds() ([]string, error) {
	query := "SELECT id FROM images ORDER BY timestamp DESC"

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}
