package variable

import (
	"time"

	"github.com/decalibrate/overlay-label-manager/internal/helper"
)

type Timer struct {
	VarType           string         `json:"type"`
	Name              string         `json:"name,omitempty"`
	Start             *time.Time     `json:"st,omitempty"`
	End               *time.Time     `json:"et,omitempty"`
	Duration          *time.Duration `json:"d,omitempty"`
	DurationRemaining *time.Duration `json:"dr,omitempty"`
	IsPaused          bool           `json:"p,omitempty"`
	CompletionText    *string        `json:"ct,omitempty"`
}

func (v Timer) Id() string {
	return v.Name
}

func (v Timer) Type() string {
	return v.VarType
}

func (v *Timer) Set(s ...string) error {
	v.VarType = "timer"

	return nil
}

func (v *Timer) Unmarshal(o map[string]interface{}) error {
	v.VarType = "timer"
	if o["name"] != nil {
		v.Name = o["name"].(string)

		if o["st"] != nil {
			st := o["st"].(time.Time)
			v.Start = &st
		}
		if o["et"] != nil {
			en := o["et"].(time.Time)
			v.End = &en
		}
		if o["d"] != nil {
			d := o["d"].(time.Duration)
			v.Duration = &d
		}
		if o["dr"] != nil {
			d := o["dr"].(time.Duration)
			v.DurationRemaining = &d
		}
		if o["p"] != nil {
			pause := o["p"].(bool)
			v.IsPaused = pause
		}
		if o["ct"] != nil {
			ct := o["ct"].(string)
			v.CompletionText = &ct
		}
	}

	return nil
}

func (v Timer) GetTokenValues(s ...string) (vt string, vs string, vn float64, vb bool) {
	baseFormat := "%d:%h:%m:%s.%S"

	if len(s) == 1 && s[0] != "" {
		s := s[0]
		switch s {
		case "start":
			vt = "s"
			vs = v.Start.Format(time.RFC3339)
		case "end":
			vt = "s"
			if v.End != nil {
				vs = v.Start.Format(time.RFC3339)
			}
		case "duration":
			vt = "s"
			vs = helper.TimeDurationFormat(*v.Duration, baseFormat)
		case "remaining":
			vt = "s"
			if v.End != nil {
				vs = helper.TimeDurationFormat(time.Until(*v.End), baseFormat)
			} else if v.DurationRemaining != nil {
				vs = helper.TimeDurationFormat(*v.DurationRemaining, baseFormat)
			} else if v.Start != nil && v.Duration != nil {
				vs = helper.TimeDurationFormat(time.Until(v.Start.Add(*v.Duration)), baseFormat)
			}
		case "completion_text":
			vt = "s"
			vs = *v.CompletionText
		case "is_paused":
			vt = "b"
			vb = v.IsPaused
		}
	} else {
		if v.Start != nil {
			if v.Duration != nil {
				if v.CompletionText != nil && v.Start.Add(*v.Duration).Before(time.Now()) {
					vt = "s"
					vs = *v.CompletionText
				} else {
					vt = "s"
					rm := time.Until(v.Start.Add(*v.Duration))
					vs = helper.TimeDurationFormat(rm, baseFormat)
				}
			} else if v.End != nil {
				if v.CompletionText != nil && time.Now().After(*v.End) {
					vt = "s"
					vs = *v.CompletionText
				} else {
					vt = "s"
					rm := time.Until(*v.End)
					vs = helper.TimeDurationFormat(rm, baseFormat)
				}
			} else {
				vt = "s"
				rm := time.Until(*v.Start)
				vs = helper.TimeDurationFormat(rm, baseFormat)
			}
		} else {
			vt = "s"
			vs = ""
		}
	}

	return vt, vs, vn, vb
}
