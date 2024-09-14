package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type SpellError struct {
	Code int      `json:"code"`
	Pos  int      `json:"pos"`
	Row  int      `json:"row"`
	Col  int      `json:"col"`
	Len  int      `json:"len"`
	Word string   `json:"word"`
	S    []string `json:"s"`
}

func CheckSpelling(text string) (bool, error) {
	apiURL := "https://speller.yandex.net/services/spellservice.json/checkText"
	form := url.Values{}
	form.Add("text", text)

	resp, err := http.Post(apiURL, "application/x-www-form-urlencoded", bytes.NewBufferString(form.Encode()))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	//fmt.Println(text)
	//fmt.Println("Response Body:", string(body))

	var errors []SpellError
	if err := json.Unmarshal(body, &errors); err != nil {
		return false, err
	}

	if len(errors) == 0 {
		return true, nil
	}
	return false, nil
}
