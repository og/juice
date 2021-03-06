package vd

import (
	"github.com/hoisie/mustache"
	ge "github.com/og/x/error"
	"regexp"
)

type StringSpec struct {
	Name string
	AllowEmpty bool
	MinRuneLen int
	MinRuneLenMessage string
	MaxRuneLen int
	MaxRuneLenMessage string
	Pattern []string
	BanPattern []string
	PatternMessage string
	Enum []string
	Ext []StringSpec
}
func (s StringSpec) NameIs(name string) StringSpec {
	s.Name = name
	return s
}
type stringSpecRender struct {
	Value interface{}
	StringSpec
}
func (spec StringSpec) render (message string, value interface{}) string {
	context := stringSpecRender{
		Value: value,
		StringSpec: spec,
	}
	return mustache.Render(message, context)
}
func (r *Rule) String(v string, spec StringSpec) {
	if r.Fail { return }
	if v == "" && !spec.AllowEmpty {
		r.Break(r.Format.StringNotAllowEmpty(spec.Name))
		return
	}
	if spec.CheckMinRuneLen(v, r) { return }
	if spec.CheckMaxRuneLen(v, r) { return }
	if spec.CheckPattern   (v, r) { return }
	if spec.CheckBanPattern(v, r) { return }
	if spec.CheckEnum(v, r) {return}
	for _, ext := range spec.Ext {
		ext.Name = spec.Name
		ext.AllowEmpty = spec.AllowEmpty
		if ext.PatternMessage == "" {
			ext.PatternMessage = spec.PatternMessage
		}
		r.String(v, ext)
		if r.Fail {return}
	}
}

func (spec StringSpec) CheckMaxRuneLen(v string, r *Rule) (fail bool) {
	if spec.MaxRuneLen == 0 {
		return false
	}
	length := len([]rune(v))
	pass := length <= spec.MaxRuneLen
	if !pass {
		message := r.CreateMessage(spec.MaxRuneLenMessage, func() string {
			return r.Format.StringMaxRuneLen(spec.Name, v, spec.MaxRuneLen)
		})
		r.Break(spec.render(message, v))
	}
	return r.Fail
}

func (spec StringSpec) CheckMinRuneLen(v string, r *Rule) (fail bool) {
	length := len([]rune(v))
	pass := length >= spec.MinRuneLen
	if !pass {
		message := r.CreateMessage(spec.MinRuneLenMessage, func() string {
			return r.Format.StringMinRuneLen(spec.Name, v, spec.MinRuneLen)
		})
		r.Break(spec.render(message, v))
	}
	return r.Fail
}
type patternData struct {
	Pattern []string
	PatternMessage string
	Name string
}
func checkPattern(data patternData, render func(string, interface{}) string, v string, r *Rule) (fail bool) {
	if len(data.Pattern) == 0 {
		return false
	}
	for _, pattern := range data.Pattern {
		matched, err := regexp.MatchString(pattern, v) ; ge.Check(err)
		pass := matched
		if !pass {
			message := r.CreateMessage(data.PatternMessage, func() string {
				return r.Format.Pattern(data.Name, v, data.Pattern, pattern)
			})
			r.Break(render(message, v))
			break
		}
	}
	return r.Fail
}

func (spec StringSpec) CheckPattern(v string, r *Rule) (fail bool) {
	return checkPattern(patternData{
		Pattern:        spec.Pattern,
		PatternMessage: spec.PatternMessage,
		Name:           spec.Name,
	}, spec.render, v, r)
}
type banPatternData struct {
	BanPattern []string
	PatternMessage string
	Name string
}
func checkBanPattern(data banPatternData, render func(string, interface{}) string, v string, r *Rule) (fail bool) {
	if len(data.BanPattern) == 0 {
		return false
	}
	for _, pattern := range data.BanPattern {
		matched, err := regexp.MatchString(pattern, v) ; ge.Check(err)
		pass := !matched
		if !pass {
			message := r.CreateMessage(data.PatternMessage, func() string {
				return r.Format.BanPattern(data.Name, v, data.BanPattern, pattern)
			})
			r.Break(render(message, v))
			break
		}
	}
	return
}
func (spec StringSpec) CheckBanPattern(v string, r *Rule) (fail bool) {
	return checkBanPattern(banPatternData{
		BanPattern:     spec.BanPattern,
		PatternMessage: spec.PatternMessage,
		Name:           spec.Name,
	}, spec.render, v, r)
}
func (spec StringSpec) CheckEnum(v string, r *Rule) (fail bool) {
	if len(spec.Enum) == 0 {
		return false
	}
	pass := false
	for _, enum := range spec.Enum {
		if enum == v {
			pass = true
		}
	}
	if !pass {
		message := r.Format.StringEnum(spec.Name, v, spec.Enum)
		r.Break(spec.render(message, v))
	}
	return r.Fail
}
