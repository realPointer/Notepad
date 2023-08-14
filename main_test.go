package main

import (
	"bufio"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestNewNotepad(t *testing.T) {
	testCases := []struct {
		name    string
		input   int
		wantLen int
		wantCap int
	}{
		{
			name:    "Zero value",
			input:   0,
			wantLen: 0,
			wantCap: 0,
		},
		{
			name:    "Positive value",
			input:   5,
			wantLen: 0,
			wantCap: 5,
		},
		{
			name:    "Negative value",
			input:   -1,
			wantLen: 0,
			wantCap: 0,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d, %v", i+1, tc.name), func(t *testing.T) {
			notepad := newNotepad(tc.input)
			if len(notepad.Notes) != tc.wantLen {
				t.Errorf("got %v, want %v len", len(notepad.Notes), tc.wantLen)
			}
			if cap(notepad.Notes) != tc.wantCap {
				t.Errorf("got %v, want %v cap", cap(notepad.Notes), tc.wantCap)
			}
		})
	}
}

func TestGetInput(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		wantCmd string
		wantTxt []string
	}{
		{
			name:    "Empty input",
			input:   "   ",
			wantCmd: "",
			wantTxt: nil,
		},
		{
			name:    "Enter a command",
			input:   "create",
			wantCmd: "create",
			wantTxt: []string{}},
		{
			name:    "Enter a command with spaces after",
			input:   "create        ",
			wantCmd: "create",
			wantTxt: []string{}},
		{
			name:    "Enter a command and a single word",
			input:   "create note",
			wantCmd: "create",
			wantTxt: []string{"note"}},
		{
			name:    "Enter a command and a few words",
			input:   "create note about tasks",
			wantCmd: "create",
			wantTxt: []string{"note", "about", "tasks"}},
		{
			name:    "Enter a command and a single word with spaces before, after and between them",
			input:   "   create   note   ",
			wantCmd: "create",
			wantTxt: []string{"note"}},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i+1), func(t *testing.T) {
			scanner := bufio.NewScanner(strings.NewReader(tc.input))
			cmd, txt := getInput(scanner)
			if !reflect.DeepEqual(cmd, tc.wantCmd) {
				t.Errorf("got command %#v, want %#v", cmd, tc.wantCmd)
			}
			if !reflect.DeepEqual(txt, tc.wantTxt) {
				t.Errorf("got text %#v, want %#v", txt, tc.wantTxt)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	tests := []struct {
		name     string
		notepad  *Notepad
		text     []string
		wantErr  bool
		expected []*Note
	}{
		{
			name:     "Add a note to an empty notepad",
			notepad:  &Notepad{},
			text:     []string{"this is a new note"},
			wantErr:  false,
			expected: []*Note{{Text: "this is a new note", Status: false}},
		},
		{
			name: "Add a note to a non-empty notepad",
			notepad: &Notepad{
				Notes: []*Note{{Text: "first note", Status: false}},
			},
			text:    []string{"second note"},
			wantErr: false,
			expected: []*Note{
				{Text: "first note", Status: false},
				{Text: "second note", Status: false},
			},
		},
		{
			name:     "Error: empty note",
			notepad:  &Notepad{},
			text:     []string{""},
			wantErr:  true,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.notepad.Create(tt.text)
			if tt.wantErr && err == nil {
				t.Errorf("expected error, but got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(tt.notepad.Notes, tt.expected) {
				t.Errorf("unexpected notepad state: got %#v, want %#v", tt.notepad.Notes, tt.expected)
			}
		})
	}
}
