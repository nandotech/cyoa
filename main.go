package cyoa

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
	story, err := UnmarshalChapter(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(story.Title)

}

func UnmarshalChapter(data []byte) (Chapter, error) {
	var r Chapter
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Chapter) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
