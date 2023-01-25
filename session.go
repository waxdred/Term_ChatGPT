package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var user = os.Getenv("USER")

var Path = fmt.Sprint("/Users/", user, "/.config/chatGpt/")

type setting struct {
	Temperature      float64 `json:"temperature"`
	TopP             float64 `json:"topp"`
	FrequencyPenalty float64 `json:"frequency_penalty"`
	PresencePenalty  float64 `json:"presence_penalty"`
	MaxTokens        int64   `json:"max_token"`
}

type Session struct {
	Id         int64   `json:"id"`
	Title      string  `json:"title"`
	Setting    setting `json:"setting"`
	Content    string  `json:"content"`
	Created_at string  `json:"create_at"`
}

type Sessions []Session

func (s Sessions) Len() int {
	return len(s)
}

func (s Sessions) Less(i, j int) bool {
	return s[i].Created_at < s[j].Created_at
}

func (s Sessions) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func sortSessions(sessions Sessions) {
	sort.Sort(sessions)
}

func (s *Session) save() error {
	name := s.Created_at + ".json"
	name = strings.Replace(name, " ", "", -1)
	pathfile := Path + name
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(pathfile, b, 0644); err != nil {
		return err
	}
	return nil
}

func (s Sessions) deleteFile(idx int) Sessions {
	var tmp []Session
	for i, se := range s {
		if i+1 == idx {
			name := strings.Replace(se.Created_at, " ", "", -1)
			filename := Path + name + ".json"
			os.Remove(filename)
		} else {
			tmp = append(tmp, se)
		}
	}
	return tmp
}

func (s Sessions) getList(idx int64) string {
	var ret string
	var count int64
	for i, se := range s {
		if i+1 >= int(idx) && count < int64(heightSession) {
			tmp := se.Created_at
			tmp = strings.Replace(tmp, " ", "", -1)
			if i+1 == int(idx) {
				tmp = styleSettingSelectTitle.Render(">" + tmp)
			} else {
				tmp = styleSettingTitle.Render(tmp)
			}
			ret = fmt.Sprintf("%s\n%s", ret, tmp)
			count++
		}
	}
	return ret
}

func (s Sessions) init() Sessions {
	if _, err := os.Stat(Path); os.IsNotExist(err) {
		os.MkdirAll(Path, os.ModePerm)
	}
	files, err := ioutil.ReadDir(Path)
	if err != nil {
		return s
	}
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}
		filePath := filepath.Join(Path, file.Name())
		b, err := ioutil.ReadFile(filePath)
		if err != nil {
			return s
		}
		var session Session
		if err := json.Unmarshal(b, &session); err != nil {
			return s
		}
		s = append(s, session)
	}
	return s
}
