package author

import (
	"fmt"
	"gonote/db"
	"gonote/struct/note"
)

type Author struct {
	Username string
	Password int
	Notes    []note.Note
}

func (author *Author) ListNotes() {
	fmt.Println("Not List:")
	for index, note := range author.Notes {
		fmt.Printf("Date: %s\nLocation: %s\nText: %s", note.Date.Format("02-01-2006 3:4:5 pm"), note.Location["city"], note.Text)
		if index != len(author.Notes) {
			fmt.Printf("\n-----------------------\n")
		}
	}
}

func (author *Author) AddNote(note note.Note) {
	author.Notes = append(author.Notes, note)
	db.DbFileExists()
	db.WriteDbFile(&author.Notes)
}

func (author *Author) DeleteNote(id int) {
	remove := func(index int) {
		author.Notes = append(author.Notes[:index], author.Notes[index+1:]...)
	}
	_, index, err := author.FindNote(id)
	if err != "" {
		fmt.Println(err)
		return
	}
	remove(index)
	db.WriteDbFile(&author.Notes)
}

func (author *Author) FindNote(id int) (note.Note, int, string) {
	var findedNote note.Note
	var findedNoteIndex int
	var errorMessage string
	for index, note := range author.Notes {
		if note.Id == id {
			findedNote = note
			findedNoteIndex = index
			break
		}
	}
	if findedNote.Id == 0 {
		errorMessage = "Note not found"
	}
	return findedNote, findedNoteIndex, errorMessage
}

func (author *Author) UpdateNote(noteIndex int, text string) {
	_, index, err := author.FindNote(noteIndex)
	if err != "" {
		fmt.Println(err)
		return
	}
	author.Notes[index].Text = text
	fmt.Println("Note Updated")
	db.WriteDbFile(&author.Notes)
}
