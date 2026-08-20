package policy_engine

import "strings"

type Explanation struct {
	RuleID   string
	Checks   []string
	Selected string
	Fallback bool
}

func Explain(p Policy, q Query) Explanation {
	e := Explanation{}
	for _, r := range p.Rules {
		ok, why := match(r, q)
		e.Checks = append(e.Checks, why)
		if ok {
			e.RuleID = r.ID
			if len(r.Targets) > 0 {
				e.Selected = r.Targets[0]
			}
			return e
		}
	}
	e.Fallback = true
	return e
}
func (e Explanation) String() string {
	parts := append([]string{}, e.Checks...)
	if e.Selected != "" {
		parts = append(parts, "selected="+e.Selected)
	}
	if e.Fallback {
		parts = append(parts, "fallback=true")
	}
	return strings.Join(parts, "; ")
}
