package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	notepadFull     = "[Error] Notepad is full"
	missingNoteArg  = "[Error] Missing note argument"
	unknownCommand  = "[Error] Unknown command: %s\n"
	invalidPosition = "[Error] Invalid position: %d\n"
	nothingToUpdate = "[Error] There is nothing to update"
	nothingToDelete = "[Error] There is nothing to delete"
	invalidArgument = "[Error] Invalid argument"
	cannotConvert   = "[Error] Cannot convert index to a number"
	notepadEmpty    = "[Info] Notepad is empty"
	noteDone        = "[Info] Note done"
	noteUndone      = "[Info] Note undone"
	noteAlreadyDone = "[Info] Note already done"
	notepadExtended = "[Info] Notepad is extended to %d notes\n"
	noteCreated     = "[OK] The note was successfully created"
	noteUpdated     = "[OK] The note was successfully updated"
	noteDeleted     = "[OK] The note was successfully deleted"
	noteStatusDone  = "Status: Done"
	farewell        = "[Info] Goodbye!"
)

type Note struct {
	Text string
	Done bool
}

type Notepad struct {
	Notes      []*Note
	MaxNotes   int
	CurrentIdx int
}

func main() {
	notepad := NewNotepad(10)
	scanner := bufio.NewScanner(os.Stdin)

	for {
		command, text := GetInput(scanner)

		switch command {
		case "extend":
			notepad.Extend(text)
		case "create":
			notepad.Create(text)
		case "done":
			notepad.Done(text)
		case "undone":
			notepad.Undone(text)
		case "update":
			notepad.Update(text)
		case "delete":
			notepad.Delete(text)
		case "list":
			notepad.List(text)
		case "clear":
			notepad.Clear(text)
		case "exit":
			notepad.Exit(text)
		default:
			fmt.Printf(unknownCommand, command)
		}
	}
}

func NewNotepad(maxNotes int) *Notepad {
	return &Notepad{
		Notes:      make([]*Note, maxNotes),
		MaxNotes:   maxNotes,
		CurrentIdx: 0,
	}
}

func GetInput(scanner *bufio.Scanner) (command string, text []string) {
	fmt.Println("Enter a command and data: ")
	fmt.Print("> ")
	scanner.Scan()
	input := scanner.Text()
	parts := strings.Fields(input)
	command = parts[0]
	text = parts[1:]
	return command, text
}

func (n *Notepad) Extend(text []string) {
	if isEmpty(text) {
		return
	}

	extend, err := convertToNumber(text[0])
	if err != nil {
		return
	}

	if extend <= 0 {
		fmt.Println("Cannot extend to a number that is negative or equal to 0")
		return
	}

	n.MaxNotes += extend
	fmt.Printf(notepadExtended, n.MaxNotes)
}

func GetIndex(text []string) (index int, err error) {
	index, err = strconv.Atoi(text[0])
	if err != nil {
		return -1, err
	}
	return index, nil
}

func (n *Notepad) Done(text []string) {
	index, err := GetIndex(text)
	if err != nil {
		fmt.Println("Note cannot be converted to done")
		return
	}
	index -= 1

	if n.Notes[index].Done == true {
		fmt.Println(noteAlreadyDone)
		return
	} else {
		n.Notes[index].Done = true
	}

	fmt.Println(noteDone)
}

func (n *Notepad) Undone(text []string) {
	index, err := GetIndex(text)
	if err != nil {
		fmt.Println("Note cannot be converted to undone")
		return
	}
	n.Notes[index-1].Done = false
	fmt.Println(noteUndone)
}

func (n *Notepad) Create(text []string) {
	note := strings.Join(text[:], " ")
	if n.CurrentIdx >= n.MaxNotes {
		fmt.Println(notepadFull)
		return
	}
	if note == "" {
		fmt.Println(missingNoteArg)
		return
	}

	newNote := &Note{Text: note, Done: false}
	n.Notes[n.CurrentIdx] = newNote
	n.CurrentIdx++
	fmt.Println(noteCreated)
}

func (n *Notepad) GetArguments(text []string) (position int, note string) {
	position, err := convertToNumber(text[0])
	if err != nil {
		return -1, ""
	}

	if position < 0 || position >= n.MaxNotes {
		fmt.Printf(invalidPosition, position)
		return -1, ""
	}

	note = strings.Join(text[1:], " ")
	return position, note
}

func (n *Notepad) Update(text []string) {
	if isEmpty(text) {
		return
	}

	index, note := n.GetArguments(text)
	if index == -1 {
		return
	}

	index -= 1

	if n.Notes[index] == nil {
		fmt.Println(nothingToUpdate)
		return
	}

	if note == "" {
		fmt.Println(missingNoteArg)
		return
	}

	n.Notes[index].Text = note
	fmt.Println(noteUpdated)
}

func (n *Notepad) Delete(text []string) {
	if isEmpty(text) {
		return
	}

	if len(text) != 1 {
		fmt.Println(invalidArgument)
		return
	}

	index, err := convertToNumber(text[0])
	if err != nil {
		return
	}
	index -= 1

	if n.Notes[index] == nil {
		fmt.Println(nothingToDelete)
		return
	}

	n.Notes = append(n.Notes[:index], n.Notes[index+1:]...)
	n.CurrentIdx -= 1
	fmt.Println(noteDeleted)
}

func (n *Notepad) List(text []string) {
	if len(text) != 0 {
		fmt.Println(invalidArgument)
		return
	}
	if n.CurrentIdx == 0 {
		fmt.Println(notepadEmpty)
		return
	}
	for i, note := range n.Notes {
		if note != nil {
			fmt.Printf("%d: %s ", i+1, note.Text)
			if note.Done == true {
				fmt.Println(noteStatusDone)
			}
			fmt.Println()
		}
	}
}

func (n *Notepad) Clear(text []string) {
	if len(text) != 0 {
		fmt.Println(invalidArgument)
		return
	}
	n.Notes = make([]*Note, n.MaxNotes)
	n.CurrentIdx = 0
}

func (n *Notepad) Exit(text []string) {
	if len(text) != 0 {
		fmt.Println(invalidArgument)
		return
	}
	fmt.Println(farewell)
	os.Exit(0)
}

func isEmpty(text []string) bool {
	if len(text) == 0 {
		fmt.Println(missingNoteArg)
		return true
	}
	return false
}

func convertToNumber(text string) (int, error) {
	number, err := strconv.Atoi(text)
	if err != nil {
		fmt.Println(cannotConvert)
		return -1, err
	}
	return number, nil
}
