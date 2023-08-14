package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

const (
	invalidArgument = "invalid argument"
	initialMaxNotes = 10
)

type Note struct {
	Text   string `json:"text"`
	Status bool   `json:"status"`
}

type Notepad struct {
	Notes []*Note
}

func main() {
	notepad := newNotepad(initialMaxNotes)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		color.Blue("Enter a command and data: ")
		fmt.Print("> ")
		command, text := getInput(scanner)
		switch command {
		case "save":
			err := notepad.Save(text)
			if err != nil {
				printError(err.Error())
			} else {
				printOK("Notepad saved to file")
			}
		case "load":
			err := notepad.Load(text)
			if err != nil {
				printError(err.Error())
			} else {
				printOK("Notepad loaded from file")
			}
		case "create":
			err := notepad.Create(text)
			if err != nil {
				printError(err.Error())
			} else {
				printOK("Note was successfully created")
			}
		case "done":
			err := notepad.SetStatus(text, true)
			if err != nil {
				printError(err.Error())
			} else {
				printOK("Note marked as done")
			}
		case "undone":
			err := notepad.SetStatus(text, false)
			if err != nil {
				printError(err.Error())
			} else {
				printOK("Note marked as not done")
			}
		case "update":
			err := notepad.Update(text)
			if err != nil {
				printError(err.Error())
			} else {
				printOK("Note updated")
			}
		case "delete":
			err := notepad.Delete(text)
			if err != nil {
				printError(err.Error())
			} else {
				printOK("note deleted")
			}
		case "list":
			notepad.List(text)
		case "clear":
			notepad.Clear(text)
		case "exit":
			notepad.Exit(text)
		default:
			printError("unknown command")
		}
	}
}

func newNotepad(maxNotes int) *Notepad {
	if maxNotes < 0 {
		maxNotes = 0
	}
	return &Notepad{
		Notes: make([]*Note, 0, maxNotes),
	}
}

func getInput(scanner *bufio.Scanner) (command string, text []string) {
	scanner.Scan()
	input := scanner.Text()
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return "", nil
	}
	command = parts[0]
	text = parts[1:]
	return command, text
}

func (n *Notepad) Save(text []string) error {
	if len(text) != 1 {
		return errors.New(invalidArgument)
	}
	filename := text[0]
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			printError("cannot close file")
		}
	}(file)

	notesJSON, err := json.Marshal(n.Notes)
	if err != nil {
		return err
	}

	_, err = file.Write(notesJSON)
	if err != nil {
		return err
	}

	return nil
}

func (n *Notepad) Load(text []string) error {
	if len(text) != 1 {
		return errors.New(invalidArgument)
	}
	filename := text[0]
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			printError("cannot close file")
		}
	}(file)

	var notes []*Note
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&notes); err != nil {
		return err
	}

	n.Notes = notes

	return nil
}

func (n *Notepad) SetStatus(text []string, status bool) error {
	if len(text) != 1 {
		return errors.New(invalidArgument)
	}
	index, err := convertToNumber(text[0])
	if err != nil {
		return err
	}
	index -= 1
	if index < 0 || index >= len(n.Notes) {
		return errors.New("[Error] invalid position")
	}
	n.Notes[index].Status = status
	return nil
}

func (n *Notepad) Create(text []string) error {
	note := strings.Join(text, " ")
	if note == "" {
		return errors.New("missing note argument")
	}
	newNote := &Note{Text: note, Status: false}
	n.Notes = append(n.Notes, newNote)
	return nil
}

func (n *Notepad) Update(text []string) error {
	if len(text) != 2 {
		return errors.New(invalidArgument)
	}

	index, err := convertToNumber(text[0])
	if err != nil {
		return err
	}
	index -= 1
	if index < 0 || index >= len(n.Notes) {
		return errors.New("incorrect position")
	}
	note := strings.Join(text[1:], " ")
	n.Notes[index].Text = note
	return nil
}

func (n *Notepad) Delete(text []string) error {
	if len(text) != 1 {
		return errors.New(invalidArgument)
	}
	index, err := convertToNumber(text[0])
	if err != nil {
		return err
	}
	index -= 1
	n.Notes = append(n.Notes[:index], n.Notes[index+1:]...)
	return nil
}

func (n *Notepad) List(text []string) {
	if !isTextEmpty(text) {
		return
	}
	if len(n.Notes) == 0 {
		color.Cyan("[Info] Notepad is empty")
		return
	}
	for i, note := range n.Notes {
		fmt.Printf("%d: %s ", i+1, note.Text)
		if note.Status == true {
			fmt.Print(" / Status: Done")
		}
		fmt.Println()
	}
}

func (n *Notepad) Clear(text []string) {
	if !isTextEmpty(text) {
		return
	}
	n.Notes = make([]*Note, 0, len(n.Notes))
}

func (n *Notepad) Exit(text []string) {
	if !isTextEmpty(text) {
		return
	}
	color.Cyan("[Info] Goodbye!")
	os.Exit(0)
}

func isTextEmpty(text []string) bool {
	if len(text) == 0 {
		return true
	} else {
		printError(invalidArgument)
		return false
	}
}

func convertToNumber(text string) (int, error) {
	number, err := strconv.Atoi(text)
	if err != nil {
		return -1, err
	}
	return number, nil
}

func printError(message string) {
	color.Red("[Error] " + message)
}

func printOK(message string) {
	color.Green("[OK] " + message)
}
