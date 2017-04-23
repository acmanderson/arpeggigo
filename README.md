# arpeggigo

`arpeggigo` <sup>cute, right?</sup> is a command-line MIDI sequencer.

## Usage

### Requirements
* [PortMIDI](http://portmedia.sourceforge.net/portmidi/) - available on apt and homebrew.

### Flags
* `interval` (`int`) - Number of milliseconds between notes. (default `100`)
* `notes` (`string`) - JSON array of objects with keys `"note"` (`int`), `"velocity"` (`int`), and `"duration"` (`int`). Duration value represents milliseconds. (default `"[]"`)

### Example

Given a file `notes.json` with the contents:

```json
[
  {"note": 60, "velocity": 100, "duration": 50},
  {"note": 64, "velocity": 75, "duration": 100},
  {"note": 67, "velocity": 50, "duration": 150},
  {"note": 72, "velocity": 25, "duration": 200}
]
```

Execute the following command to loop over the notes defined in `notes.json` with 500 milliseconds between each note:
```bash
./arpeggigo -notes "$(< notes.json)" -interval 500
```