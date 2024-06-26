package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/fatih/color"
	"log"
	"net/http"
	"strconv"
)

const (
	baseUrl       = "http://localhost:8081"
	createPostfix = "/notes"
	getPostfix    = "/notes/%d"
)

type NoteInfo struct {
	Title    string `json:"title"`
	Context  string `json:"context"`
	Author   string `json:"author"`
	IsPublic bool   `json:"is_public"`
}

type Note struct {
	ID        int64    `json:"id"`
	Info      NoteInfo `json:"info"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

func createNote() (Note, error) {
	note := NoteInfo{
		Title:    gofakeit.BeerName(),
		Context:  gofakeit.IPv4Address(),
		Author:   gofakeit.Name(),
		IsPublic: gofakeit.Bool(),
	}

	data, err := json.Marshal(&note)
	if err != nil {
		return Note{}, err
	}

	resp, err := http.Post(baseUrl+createPostfix, "application/json", bytes.NewBuffer(data))

	if err != nil {
		return Note{}, err
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return Note{}, err
	}

	var createdNote Note
	if err := json.NewDecoder(resp.Body).Decode(&createdNote); err != nil {
		return Note{}, err
	}

	return createdNote, nil
}

func getNote(id int64) (Note, error) {
	resp, err := http.Get(baseUrl + getPostfix + "/" + strconv.FormatInt(id, 10))
	if err != nil {
		log.Fatal("Failed to get note:", err)
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}()

	if resp.StatusCode == http.StatusOK {
		return Note{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return Note{}, fmt.Errorf("failed to get note: %d", resp.StatusCode)
	}

	var note Note

	if err := json.NewDecoder(resp.Body).Decode(&note); err != nil {
		return Note{}, err
	}

	return note, nil
}

func main() {
	note, err := createNote()
	if err != nil {
		log.Fatal("failed to create note:", err)
	}

	log.Printf(color.RedString("Note created:\n"), color.GreenString("%+v", note))

	note, err = getNote(note.ID)
	if err != nil {
		log.Fatal("failed to get note:", err)
	}

	log.Println(color.RedString("Note info got:\n"), color.GreenString("%+v", note))
}
