package vd

import (
	"github.com/hoisie/mustache"
	gconv "github.com/og/x/conv"
)

type IntSpec struct {
	Name string
	// AllowZero bool // 暂时取消 AllowZero，目的是降低使用者学习成本，观察一段时间后再决定是否完全去掉 (2020年08月07日 by @nimoc)
	Min OptionInt
	MinMessage string
	Max OptionInt
	MaxMessage string
	Pattern []string
	BanPattern []string
	PatternMessage string
}
type intSpecRender struct {
	Value interface{}
	IntSpec
}
func (spec IntSpec) render (message string, value interface{}) string {
	context := intSpecRender{
		Value: value,
		IntSpec: spec,
	}
	return mustache.Render(message, context)
}
func (r *Rule) Uint(v uint, spec IntSpec) {
	r.Int(int(v), spec)
}
func (r *Rule) Uint8(v uint8, spec IntSpec) {
	r.Int(int(v), spec)
}
func (r *Rule) Uint16(v uint16, spec IntSpec) {
	r.Int(int(v), spec)
}
func (r *Rule) Uint32(v uint32, spec IntSpec) {
	r.Int(int(v), spec)
}
func (r *Rule) Uint64(v uint64, spec IntSpec) {
	r.Int(int(v), spec)
}
func (r *Rule) Int8(v int8, spec IntSpec) {
	r.Int(int(v), spec)
}
func (r *Rule) Int16(v int16, spec IntSpec) {
	r.Int(int(v), spec)
}
func (r *Rule) Int32(v int32, spec IntSpec) {
	r.Int(int(v), spec)
}
func (r *Rule) Int64(v int64, spec IntSpec) {
	r.Int(int(v), spec)
}
func (r *Rule) Int(v int, spec IntSpec) {
	if r.Fail {return}
	// if v == 0 && !spec.AllowZero {
	// 	r.Break(r.Format.IntNotAllowEmpty(spec.Name))
	// 	return
	// }
	if spec.CheckMin(v, r) { return }
	if spec.CheckMax(v ,r) { return }
	if spec.CheckPattern(v, r) {return}
	if spec.CheckBanPattern(v, r) {return}
	return
}
func (spec IntSpec) CheckMin(v int, r *Rule) (fail bool) {
	if !spec.Min.Valid() {
		return
	}
	min := spec.Min.Unwrap()
	pass := v >= min
	if !pass {
		message := r.CreateMessage(spec.MinMessage, func() string {
			return r.Format.IntMin(spec.Name, v, min)
		})
		r.Break(spec.render(message, v))
	}
	return
}
func (spec IntSpec) CheckMax(v int, r *Rule) (fail bool) {
	if !spec.Max.Valid() {
		return
	}
	max := spec.Max.Unwrap()
	pass := v <= max
	if !pass {
		message := r.CreateMessage(spec.MaxMessage, func() string {
			return r.Format.IntMax(spec.Name, v, max)
		})
		r.Break(spec.render(message, v))
	}
	return
}
func (spec IntSpec) CheckPattern(v int, r *Rule) (fail bool) {
	return checkPattern(patternData{
		Pattern:        spec.Pattern,
		PatternMessage: spec.PatternMessage,
		Name:           spec.Name,
	}, spec.render, gconv.IntString(v), r)
}

func (spec IntSpec) CheckBanPattern(v int, r *Rule) (fail bool) {
	return checkBanPattern(banPatternData{
		BanPattern:        spec.BanPattern,
		PatternMessage: spec.PatternMessage,
		Name:           spec.Name,
	}, spec.render, gconv.IntString(v), r)
}
