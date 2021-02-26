package json

import (
	"bytes"
	"crypto/sha1"
	"errors"
)

func Replace(w *bytes.Buffer, b []byte, from, to []Field) error {
	if len(from) != len(to) {
		return errors.New("'from' and 'to' must be of the same length")
	}

	state := expectKey
	ws, we := 0, len(b)

	s := 0
	fi := -1

	fmap := make(map[[20]byte]int, (len(from) * 2))
	tmap := make(map[[20]byte]int, (len(from)))

	for i, f := range from {
		h1 := sha1.Sum(f.Key)
		n, ok := fmap[h1]
		if !ok {
			fmap[h1] = i
			n = i
		}

		h2 := sha1.Sum(f.Value)
		fmap[h2] = n
		tmap[h2] = i
	}

	for i := 0; i < len(b); i++ {
		switch {
		case ws == 0 && b[i] == '{' || b[i] == '[':
			ws = i

		case state == expectKey && b[i] == '"':
			state = expectKeyClose
			s = i + 1
			continue

		case state == expectKeyClose && b[i] == '"':
			state = expectColon

		default:
			continue
		}

		if state != expectColon {
			continue
		}

		h1 := sha1.Sum(b[s:i])
		if n, ok := fmap[h1]; ok {
			we = s
			fi = n
		}

		e := 0
		d := 0
		for ; i < len(b); i++ {
			if state == expectObjClose || state == expectListClose {
				switch b[i] {
				case '{', '[':
					d++
				case '}', ']':
					d--
				}
			}

			switch {
			case state == expectColon && b[i] == ':':
				state = expectValue

			case state == expectValue && b[i] == '"':
				state = expectString
				s = i

			case state == expectString && b[i] == '"':
				e = i

			case state == expectValue && b[i] == '[':
				state = expectListClose
				s = i
				d++

			case state == expectListClose && d == 0 && b[i] == ']':
				e = i

			case state == expectValue && b[i] == '{':
				state = expectObjClose
				s = i
				d++

			case state == expectObjClose && d == 0 && b[i] == '}':
				e = i

			case state == expectValue && (b[i] >= '0' && b[i] <= '9'):
				state = expectNumClose
				s = i

			case state == expectNumClose &&
				((b[i] < '0' || b[i] > '9') &&
					(b[i] != '.' && b[i] != 'e' && b[i] != 'E' && b[i] != '+' && b[i] != '-')):
				i--
				e = i

			case state == expectValue &&
				(b[i] == 'f' || b[i] == 'F' || b[i] == 't' || b[i] == 'T'):
				state = expectBoolClose
				s = i

			case state == expectBoolClose && (b[i] == 'e' || b[i] == 'E'):
				e = i
			}

			if e != 0 {
				e++
				h2 := sha1.Sum(b[s:e])

				if n, ok1 := fmap[h2]; ok1 && n == fi {
					ti, ok2 := tmap[h2]

					if ok2 {
						if _, err := w.Write(b[ws:we]); err != nil {
							return err
						}
						if _, err := w.Write(to[ti].Key); err != nil {
							return err
						}
						if _, err := w.WriteString(`":`); err != nil {
							return err
						}
						if _, err := w.Write(to[ti].Value); err != nil {
							return err
						}
						ws = e
					}
				}

				if ws != e && (b[s] == '[' || b[s] == '{') {
					i = s - 1
				}

				state = expectKey
				we = len(b)
				fi = -1
				break
			}
		}
	}

	w.Write(b[ws:we])
	return nil
}
