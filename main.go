package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/micmonay/keybd_event"
)

var kb keybd_event.KeyBonding

func getKey(key int) int {
	if key == 0 {
		return keybd_event.VK_0
	}
	return key + 1
}

func main() {
	time.Sleep(5 * time.Second)

	var err error
	kb, err = keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}

	song1 := "35 66 67 88 89 77 655 l6 | 35 66 67 88 89 77 65 l6 | 35 66 67 88 90 Sh4Sh4 090 6 | 67 889 0 6 68 77 86 l7 (repeat) l0 Shl4 000 09 | l9 l8 787 l6 (repeat) 67 8 90 987 890 l9 | 89 l0 98 787 66 75 l6 | 67 8 78 989 098 l6 | 67 890 Sh469 l8 l7 l0 Shl4 000 09 | l9 l8 787 l6 (repeat) "

	song2 := "5 6 7 7 | 6 5 6 7 5 2 | 5 6 7 7 | 6 5 6 7 l5 | 2254579 | 9098765 | 999 999 | l9 | 5 6 7 7 | 6 5 6 7 5 2 | 5 6 7 7 | 6 5 6 7 l5"

	song3 := "0 7 8 9 8 7 | 6 6 8 0 9 8 | 7 7 8 9 0 | 8 6 6 | 9 sh4 sh6 sh5 sh4 | 0 8 0 9 8 | 7 7 8 9 0"

	song4 := "5 5 5 2 3 3 l2 7 7 6 6 l5 | 5 5 5 2 3 3 l2 7 7 6 6 l5 | 2 2 5 5 l5 2 2 5 5 l5 5 5 5 5 5 5 5 5 5 5 l5 l5 | 5 5 5 2 3 3 l2 7 7 6 6 l5"

	song5 := "469 469 469 469 469 469 469 469 | 369 369 369 | 368 368 368 368 368 368 368 368"

	song6 := "24l9 24l9 shl3 sh4sh3sh4sh3 8l6 6 2 45l6 6 2 45l3"

	song7 := "spl8 spl7 63 | 3333 6666 6564 | 4444 6666 7 8 l5 | 5555 888 99 l7 | l8 l7 l6 l3"

	playlist := []string{song1, song2, song3, song4, song5, song6, song7}

	song := strings.Join(playlist, "(pause)")

	song = strings.ToLower(song)

	repeatDone := make([]bool, len(song))
	var sb strings.Builder

	lastPointer := -1

	longMod := false
	shMod := false
	spMod := false

	for i := 0; i < len(song); i++ {

		note := song[i]

		if note == ' ' {
			time.Sleep(100 * time.Millisecond)
			longMod = false
			shMod = false
			spMod = false
			continue
		}

		if note == '|' {
			lastPointer = i
			continue
		}

		if note == 'l' {
			longMod = true
			continue
		}

		noteR := int(note - '0')
		isDigit := noteR >= 0 && noteR <= 9

		if (note == ')' || isDigit) && sb.Len() > 0 {

			if !isDigit {
				sb.WriteByte(note)
			}

			str := sb.String()
			sb.Reset()
			fmt.Printf("mod: %v\n", str)
			if str == "(repeat)" && !repeatDone[i] {
				repeatDone[i] = true
				i = lastPointer
				fmt.Printf("Goto: %v\n", i)
				continue
			}

			if str == "(pause)" {
				time.Sleep(2 * time.Second)
				continue
			}

			if str == "sh" {
				shMod = true
			}

			if str == "sp" {
				spMod = true
			}

			if !isDigit {
				continue
			}

		}

		if isDigit {

			kb.SetKeys(getKey(noteR))

			dur := 100

			if longMod {
				dur += 100
			}

			if spMod {
				kb.AddKey(keybd_event.VK_SPACE)
			}

			kb.HasSHIFT(shMod)

			fmt.Printf("Playing %v for %v ms space: %t shift: %t index: %v\n", noteR, dur, spMod, shMod, i)
			kb.Press()
			time.Sleep(time.Duration(dur) * time.Millisecond)
			kb.Release()
			time.Sleep(30 * time.Millisecond)
			continue
		}

		sb.WriteByte(note)
	}

}
