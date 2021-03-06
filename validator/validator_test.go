package vd_test

import (
	gtest "github.com/og/x/test"
	vd "github.com/og/juice/validator"
	"testing"
)

type RequiredOne struct {
	Name string
}
func (v RequiredOne) VD(r *vd.Rule){
	r.String(v.Name, vd.StringSpec{
		Name: "姓名",
	})
}
func Test_RequiredOne (t *testing.T) {
	c := vd.NewCN()
	as := gtest.NewAS(t)
	as.Equal(c.Check(RequiredOne{}), vd.Report{
		Fail:    true,
		Message: "姓名必填",
	})
	as.Equal(c.Check(RequiredOne{Name:"n"}), vd.Report{
		Fail:    false,
		Message: "",
	})
}
type RequiredTwo struct {
	Name string
	Title string
}
func (v RequiredTwo) VD(r *vd.Rule){
	r.String(v.Name, vd.StringSpec{
		Name: "姓名",
	})
	r.String(v.Title, vd.StringSpec{
		Name: "标题",
	})
}
func Test_RequiredTwo (t *testing.T) {
	c := vd.NewCN()
	as := gtest.NewAS(t)
	as.Equal(c.Check(RequiredTwo{}), vd.Report{
		Fail:    true,
		Message: "姓名必填",
	})
	as.Equal(c.Check(RequiredTwo{Name:"n"}), vd.Report{
		Fail:    true,
		Message: "标题必填",
	})
	as.Equal(c.Check(RequiredTwo{Name:"n",Title:"1"}), vd.Report{
		Fail:    false,
		Message: "",
	})
}
type RequiredThree struct {
	Name string
	Title string
}
func (v RequiredThree) VD(r *vd.Rule){
	r.String(v.Name, vd.StringSpec{
		Name: "姓名",
		AllowEmpty: true,
	})
	r.String(v.Title, vd.StringSpec{
		Name: "标题",
	})
}
func Test_RequiredThree (t *testing.T) {
	c := vd.NewCN()
	as := gtest.NewAS(t)
	as.Equal(c.Check(RequiredThree{}), vd.Report{
		Fail:    true,
		Message: "标题必填",
	})
	as.Equal(c.Check(RequiredThree{Name:"n",Title:"1"}), vd.Report{
		Fail:    false,
		Message: "",
	})
}
type RequiredFour struct {
	Name  string
	Title string
}
func (v RequiredFour) VD(r *vd.Rule){
	r.String(v.Name, vd.StringSpec{
		Name: "姓名",
	})
	r.String(v.Title, vd.StringSpec{
		Name: "标题",
	})
}
func Test_RequiredFour (t *testing.T) {
	c := vd.NewCN()
	as := gtest.NewAS(t)
	as.Equal(c.Check(RequiredFour{}), vd.Report{
		Fail:    true,
		Message: "姓名必填",
	})
	as.Equal(c.Check(RequiredFour{Name:"n",Title:""}), vd.Report{
		Fail:    true,
		Message: "标题必填",
	})
}
