package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type VisionContent struct {
	IsUrlType bool
	Text      string
}

type VisionMessage struct {
	Role    string          `json:"role"`
	Content []VisionContent `json:"content"`
}

type VisionRequestBody struct {
	Model    string          `json:"model"`
	Messages []VisionMessage `json:"messages"`
}

type DallERequestBody struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	N      int    `json:"n"`
	Size   string `json:"size"`
}

type ChatResponseChoice struct {
	Index   int         `json:"index"`
	Message ChatMessage `json:"message"`
}

type ChatResponse struct {
	Id      string               `json:"id"`
	Choices []ChatResponseChoice `json:"choices"`
}

type DallEData struct {
	RevisedPrompt string `json:"revised_prompt"`
	URL           string `json:"url"`
}

type DallEResponse struct {
	Created int         `json:"created"`
	Data    []DallEData `json:"data"`
}

func (content VisionContent) MarshalJSON() ([]byte, error) {
	if content.IsUrlType {
		return json.Marshal(struct {
			Type     string `json:"type"`
			ImageURL string `json:"image_url"`
		}{
			Type:     "image_url",
			ImageURL: content.Text,
		})
	}

	return json.Marshal(struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}{
		Type: "text",
		Text: content.Text,
	})
}

func PromptFromURL(url string, isStyle bool) (string, error) {
	var instruction string
	if isStyle {
		instruction = "Create a prompt for Dall-E 3 to recreate the style of this image. Your answer should only contain information about the style of the image (i.e. tone, color scale, technique), NOT the contents"
	} else {
		instruction = "Create a prompt for Dall-E 3 to recreate the scene depicted in this image. Your answer should only contain information about the objects (i.e. people, animals, items, surroundings etc.) in the image, NOT the style in which the image"
	}

	contents := []VisionContent{{false, instruction}, {true, url}}
	messages := []VisionMessage{{"user", contents}}
	requestBody := VisionRequestBody{"gpt-4-vision-preview", messages}
	requestBodyMarshalled, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(requestBodyMarshalled))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("OPENAI_API_KEY")))

	client := http.Client{Timeout: 100 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var chatResponse ChatResponse
	responseBody, err := io.ReadAll(response.Body)
	json.Unmarshal(responseBody, &chatResponse)
	if err != nil {
		log.Printf("Something went wrong!")
	}

	return chatResponse.Choices[0].Message.Content, nil
}

func GenerateDallEImage(scene string, style string, size string) (string, string, error) {
	// TODO: Re-use parts of this that are identical to PromptFromURL()
	prompt := scene + "\n\n" + style
	requestBody := DallERequestBody{"dall-e-3", prompt, 1, size}

	requestBodyMarshalled, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/images/generations", bytes.NewReader(requestBodyMarshalled))
	req.Header.Set("Content-Type", "application/json")
	apiKey := os.Getenv("OPENAI_API_KEY")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := http.Client{Timeout: 100 * time.Second}
	response, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer response.Body.Close()

	var dalleResponse DallEResponse
	responseBody, err := io.ReadAll(response.Body)
	json.Unmarshal(responseBody, &dalleResponse)
	if err != nil {
		return "", "", err
	}

	return dalleResponse.Data[0].URL, prompt, nil
}
