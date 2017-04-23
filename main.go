package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/rakyll/portmidi"
	"os"
	"os/signal"
	"time"
)

type MIDINote struct {
	Note     int64 `json:"note"`
	Velocity int64 `json:"velocity"`
	Duration int64 `json:"duration"`
}

type MIDIEvent struct {
	Status int64
	Note   MIDINote
}

func playNote(note MIDINote, noteChannel chan<- MIDIEvent) {
	noteChannel <- MIDIEvent{0x90, note}
	time.Sleep(time.Duration(note.Duration) * time.Millisecond)
	noteChannel <- MIDIEvent{0x80, note}
}

func playSequence(notes []MIDINote, timeBetweenNotes time.Duration, noteChannel chan<- MIDIEvent) {
	for {
		for _, note := range notes {
			go playNote(note, noteChannel)
			time.Sleep(timeBetweenNotes)
		}
	}
}

func main() {
	notesJson := flag.String("notes", "[]", `JSON array of objects with keys "note" (int), "velocity" (int), and "duration" (int). Duration value represents milliseconds.`)
	interval := flag.Int("interval", 100, "Number of milliseconds between notes.")
	flag.Parse()

	fmt.Println("Starting sequence. Press Ctrl+C to quit...")

	portmidi.Initialize()
	defer portmidi.Terminate()

	out, err := portmidi.NewOutputStream(portmidi.DefaultOutputDeviceID(), 1024, 0)
	if err != nil {
		panic(err.Error())
	}

	noteChannel := make(chan MIDIEvent)
	notes := []MIDINote{}
	err = json.Unmarshal([]byte(*notesJson), &notes)
	if err != nil {
		panic(err.Error())
	}

	go playSequence(notes, time.Duration(*interval)*time.Millisecond, noteChannel)

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

LOOP:
	for {
		select {
		case event := <-noteChannel:
			out.WriteShort(event.Status, event.Note.Note, event.Note.Velocity)
		case <-sigint:
			break LOOP
		}
	}

	// send off for every note in the sequence. There's an "all notes off" MIDI message
	// but it wasn't working on the only MIDI synth that I have.
	for _, note := range notes {
		out.WriteShort(0x80, note.Note, note.Velocity)
	}
	out.Close()
}
