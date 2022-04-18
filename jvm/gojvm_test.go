package jvm

import "testing"

func TestStaticInvoke(t *testing.T) {
	jvm := NewJVM(".", "128m", "512m", "384m", "512k")
	env := jvm.Attach()
	c1 := env.FindClass("Hello")
	s0, _ := c1.InvokeString("hello", "rarnu", 8)
	t.Logf("%v\n", s0)
	env.Detach()
	jvm.Free()
}

func TestStaticGetClass(t *testing.T) {
	jvm := NewJVM(".", "128m", "512m", "384m", "512k")
	env := jvm.Attach()

	c1 := env.FindClass("Hello")
	jo := c1.GetObject("h", "H1")
	t.Logf("%v\n", jo)
	jo.Free()
	c1.Free()

	env.Detach()
	jvm.Free()
}

func TestNewClass(t *testing.T) {
	jvm := NewJVM(".", "128m", "512m", "384m", "512k")
	env := jvm.Attach()

	c1 := env.NewObject("Hello")
	t.Logf("%v\n", c1)
	c1.Free()

	env.Detach()
	jvm.Free()
}

func TestClassFields(t *testing.T) {

	jvm := NewJVM(".", "128m", "512m", "384m", "512k")
	env := jvm.Attach()

	c1 := env.NewObject("H1")
	t.Logf("%v\n", c1)

	s1 := c1.GetString("v")
	t.Logf("%v\n", s1)

	c1.SetString("v", "rarnu")
	s2 := c1.GetString("v")
	t.Logf("%v\n", s2)

	c1.Free()

	env.Detach()
	jvm.Free()

}

func TestJar(t *testing.T) {
	jvm := NewJVM("./jarsample.jar:.", "128m", "512m", "384m", "512k")
	env := jvm.Attach()

	c1 := env.FindClass("SampleClass")
	s0, _ := c1.InvokeString("helloStatic", "rarnu")
	t.Logf("%v\n", s0)

	o1 := env.NewObject("SampleClass")
	v1 := o1.GetInt("v1")
	t.Logf("%v\n", v1)

	o1.Free()

	c1.Free()
	env.Detach()
	jvm.Free()
}
