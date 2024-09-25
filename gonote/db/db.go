package db

import (
	"encoding/json"
	"fmt"
	"gonote/struct/note"
	"io/ioutil"
	"os"
)

const dbFile string = "db.json"

func DbFileExists() bool {
	_, err := os.Stat(dbFile)
	if os.IsNotExist(err) {
		panic(err)
	}
	return true
}

func ReadDbFile() []note.Note {
	data, err := ioutil.ReadFile(dbFile)
	if err != nil {
		fmt.Print("Dosya okunamadı")
	}
	var notes []note.Note

	json.Unmarshal(data, &notes)
	return notes
}

func WriteDbFile(notes *[]note.Note) {
	data, _ := json.Marshal(notes)
	err := ioutil.WriteFile(dbFile, data, 0644)
	if err != nil {
		fmt.Println("Veritabanı Güncellenirken bir sorun oluştu")
	} else {
		fmt.Println("Veritabanı Güncellendi")
	}
}
