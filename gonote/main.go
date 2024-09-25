package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gonote/db"
	"gonote/struct/author"
	"gonote/struct/note"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func commandController() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("command: ")
	command, _ := reader.ReadString('\n')
	command = strings.Trim(command, "\r\n ")
	return command
}

func readInput() string {
	reader := bufio.NewReader(os.Stdin)
	noteid, _ := reader.ReadString('\n')
	str := strings.TrimRight(noteid, "\r\n")
	return str
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	tutku := author.Author{
		Username: "tutkuonight",
		Notes:    []note.Note{},
	}

	data := db.ReadDbFile()
	tutku.Notes = data

	for {
		command := commandController()
		switch command {
		case "list":
			tutku.ListNotes()
		case "add":
			fmt.Print("Write your note: ")
			text, _ := reader.ReadString('\n')
			res, err := http.Get("http://ip-api.com/json/")
			var location map[string]interface{}
			if err != nil {
				panic(err)
			}
			json.NewDecoder(res.Body).Decode(&location)
			fmt.Println(location)
			tutku.AddNote(note.Note{Id: len(tutku.Notes) + 1, Text: text, Location: location, Date: time.Now()})
		case "delete":
			fmt.Print("Note Id: ")
			str := readInput()
			marks, _ := strconv.Atoi(str)
			tutku.DeleteNote(marks)
		case "update":
			fmt.Print("[Note Id, Text]: ")
			str := readInput()
			var splitted []string = strings.Split(str, " ")
			marks, _ := strconv.Atoi(splitted[0])
			tutku.UpdateNote(marks, strings.Join(splitted[1:], " "))
		case "exit":
			os.Exit(0)
		}
	}
}
