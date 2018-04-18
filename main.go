package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	data, err := ioutil.ReadFile("gopher.json")
	if err != nil {
		log.Fatal(err)
	}
	story, err := unmarshalStory(data)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range story.Debate.Options {
		fmt.Println(v.Arc + ": " + v.Text + "\n")
	}
}

func unmarshalStory(data []byte) (Story, error) {
	var r Story
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Story) marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Story type from gopher.json
type Story struct {
	Intro     Debate `json:"intro"`
	NewYork   Debate `json:"new-york"`
	Debate    Debate `json:"debate"`
	SeanKelly Debate `json:"sean-kelly"`
	MarkBates Debate `json:"mark-bates"`
	Denver    Debate `json:"denver"`
	Home      Debate `json:"home"`
}

// Debate type from gopher
type Debate struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

// Option type
type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}
