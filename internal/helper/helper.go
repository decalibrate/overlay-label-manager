package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

func SaveJSONFile(filePath string, data interface{}) error {

	if filePath != "" {
		if ba, err := json.Marshal(data); err == nil {
			return os.WriteFile(filePath, ba, 0755)
		} else {
			return err
		}
	}

	return errors.New("file path not specified")
}

func ReadJSONFile(filePath string, destinationStruct interface{}) error {
	// check file exists
	if _, err := os.Stat(filePath); err != nil {
		return err
	}

	var fc []byte
	var err error
	// try read file content
	if fc, err = os.ReadFile(filePath); err != nil {
		return err
	}
	// try unmarshal to struct
	if err := json.Unmarshal(fc, destinationStruct); err != nil {
		return err
	}
	// all good
	return nil
}

func ParseFloat(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}

func FormatFloat(fl float64, prec int) string {
	return strconv.FormatFloat(fl, 'f', prec, 64)
}

func TimeDurationFormat(t time.Duration, f string) string {

	var tms = float64(t.Milliseconds())
	var ms, cs, ds, s, m, h, d float64
	var ud, uh, um, us, ums bool
	p := make([]*float64, 0)

	if f == "" {
		ud, uh, um, us, ums = true, true, true, true, true
		f = "%02.[1]fd %02.[2]fh %02.[3]fm %02.[4]fs"
		p = append(p, &d, &h, &m, &s, &ms)
	} else if f == "%F" {
		ud, uh, um, us, ums = true, true, true, true, true
		p = append(p, &d, &h, &m, &s, &ms)
		f = "%02.[4]fs"
		if tms >= 1000*60 {
			f = "%02.[3]fm " + f
		}
		if tms >= 1000*60*60 {
			f = "%02.[2]fh " + f
		}
		if tms >= 1000*60*60*24 {
			f = "%02.[1]fd " + f
		}
	} else {
		// escape special characters
		f = strings.ReplaceAll(f, "\\\\", "[~backslash~]")
		f = strings.ReplaceAll(f, "\\%", "[~percent~]")

		pieces := strings.Split(f, "%")

		if len(pieces) > 1 {
			// escape percent sign if it is not followed by a known replacement value
			// e.g %p will appear as %p in final duration string
			for k, s := range pieces {
				if len(s) >= 1 && !strings.Contains("dhmsS", string(s[0])) {
					pieces[k] = "%" + s
				}
			}

			f = strings.Join(pieces, "%")

			rplcmnts := []struct {
				o, d string
				p    *float64
				u    *bool
			}{
				{"%dd", "%02.[x]f", &d, &ud},
				{"%d", "%.[x]f", &d, &ud},
				{"%hh", "%02.[x]f", &h, &uh},
				{"%h", "%.[x]f", &h, &uh},
				{"%mm", "%02.[x]f", &m, &um},
				{"%m", "%.[x]f", &m, &um},
				{"%ss", "%02.[x]f", &s, &us},
				{"%s", "%.[x]f", &s, &us},
				{"%ss", "%02.[x]f", &s, &us},
				{"%s", "%.[x]f", &s, &us},
				{"%SSS", "%03.[x]f", &ms, &ums}, // millisecond
				{"%SS", "%02.[x]f", &cs, &ums},  // centisecond
				{"%S", "%01.[x]f", &ds, &ums},   // decisecond
			}

			for _, r := range rplcmnts {
				for strings.Contains(f, r.o) {
					*r.u = true
					p = append(p, r.p)
					n := len(p)
					f = strings.Replace(f, r.o, strings.Replace(r.d, "x", strconv.Itoa(n), 1), 1)
				}
			}
			f = strings.ReplaceAll(f, "\\d", "d")
			f = strings.ReplaceAll(f, "\\h", "h")
			f = strings.ReplaceAll(f, "\\m", "m")
			f = strings.ReplaceAll(f, "\\s", "s")
			f = strings.ReplaceAll(f, "\\S", "S")
			f = strings.ReplaceAll(f, "\\S", "S")
		}

		// put in real values in placeholders
		f = strings.ReplaceAll(f, "[~percent~]", "%%")
		f = strings.ReplaceAll(f, "[~backslash~]", "\\")
	}

	// days
	if ud {
		d = math.Floor(tms / (24 * 60 * 60 * 1000))
		tms = math.Mod(tms, (24 * 60 * 60 * 1000))
	}
	if uh {
		h = math.Floor(tms / (60 * 60 * 1000))
		tms = math.Mod(tms, (60 * 60 * 1000))
	}
	if um {
		m = math.Floor(tms / (60 * 1000))
		tms = math.Mod(tms, (60 * 1000))
	}
	if us {
		s = math.Floor(tms / 1000)
		tms = math.Mod(tms, 1000)
	}
	if ums {
		ms = math.Floor(tms)
		cs = math.Floor(tms / 10)
		ds = math.Floor(tms / 100)
	}

	pi := make([]interface{}, 0)
	for _, n := range p {
		pi = append(pi, *n)
	}

	return fmt.Sprintf(f, pi...)
}
