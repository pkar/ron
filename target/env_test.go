package target

import (
	"bytes"
	"strings"
	"testing"
)

type parseOSEnvsTest struct {
	name string
	in   []string
	out  MSS
}

var testNewEnvConfig = `
- UNAME: plan9
- RON: was here
- CMD: +echo blah
- ENVS: >-
    -e CMD=$CMD
    -e TEST=$UNAME
- NOOP:
`

func TestEnv_NewEnv(t *testing.T) {
	writer := &bytes.Buffer{}
	envs, _, err := BuiltinDefault()
	ok(t, err)
	NewEnv(&RawConfig{Envs: envs}, MSS{}, writer)
}

func TestEnv_NewEnvStdout(t *testing.T) {
	envs, _, err := BuiltinDefault()
	ok(t, err)
	NewEnv(&RawConfig{Envs: envs}, MSS{}, nil)
}

func TestEnv_ParseOSEnvs(t *testing.T) {
	var parseOSEnvTests = []parseOSEnvsTest{
		{"", []string{"a="}, MSS{"a": ""}},
		{"", []string{"a=b"}, MSS{"a": "b"}},
		{"", []string{"a=b = 1"}, MSS{"a": "b = 1"}},
		{"", []string{"b=> &>1"}, MSS{"b": "> &>1"}},
	}
	for _, tt := range parseOSEnvTests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseOSEnvs(tt.in)
			equals(t, tt.out, got)
		})
	}
}

func TestEnv_Process(t *testing.T) {
	writer := &bytes.Buffer{}
	e, err := NewEnv(&RawConfig{Envs: testNewEnvConfig}, ParseOSEnvs([]string{"APP=ron", "HELLO=hello", "ABC=+pwd"}), writer)
	ok(t, err)
	err = e.Process()
	ok(t, err)

	got := e.Config["HELLO"]
	want := `hello`
	if !strings.Contains(got, want) {
		t.Fatalf("config does not contain %s got \n%s", want, got)
	}

	got = e.Config["UNAME"]
	want = `plan9`
	equals(t, want, got)

	got = e.Config["APP"]
	want = "ron"
	if got != want {
		for _, k := range e.keyOrder {
			t.Log(k, e.Config[k])
		}
		equals(t, want, got)
	}
}

func TestEnv_ProcessEnv(t *testing.T) {
	writer := &bytes.Buffer{}
	e, err := NewEnv(&RawConfig{Envs: testNewEnvConfig}, ParseOSEnvs([]string{}), writer)
	ok(t, err)
	err = e.Process()
	ok(t, err)

	got := e.Config["ENVS"]
	want := `-e CMD=blah -e TEST=plan9`
	if !strings.Contains(got, want) {
		t.Fatalf("config ENVS does not contain %s got \n%s", want, got)
	}
}

func TestEnv_ProcessBadCommand(t *testing.T) {
	writer := &bytes.Buffer{}
	//envs, _, err := BuiltinDefault()
	//ok(t, err)
	e, err := NewEnv(&RawConfig{Envs: testNewEnvConfig + "\nHELLO=+hello"}, ParseOSEnvs([]string{}), writer)
	ok(t, err)
	err = e.Process()
	if err == nil {
		t.Fatal("expected err processing command +hello")
	}
}

func TestEnv_ProcessBadYaml(t *testing.T) {
	e, err := NewEnv(&RawConfig{Envs: `:"`}, MSS{}, nil)
	ok(t, err)
	err = e.Process()
	if err == nil {
		t.Fatal("should have gotten invalid err")
	}
}

func TestEnv_ProcessBadYamlNewEnvs(t *testing.T) {
	e, err := NewEnv(&RawConfig{Envs: `:"`}, MSS{}, nil)
	ok(t, err)
	err = e.Process()
	if err == nil {
		t.Fatal("should have gotten invalid err")
	}
}

func TestEnv_List(t *testing.T) {
	writer := &bytes.Buffer{}
	envs, _, err := BuiltinDefault()
	ok(t, err)
	e, err := NewEnv(&RawConfig{Envs: envs}, MSS{}, writer)
	ok(t, err)
	err = e.Process()
	ok(t, err)
	e.List()
	got := writer.String()
	want := "ron\n"
	if !strings.Contains(got, want) {
		t.Errorf("output does not contain %s got \n%s", want, got)
	}
}

func TestEnv_ListBadWriter(t *testing.T) {
	envs, _, err := BuiltinDefault()
	ok(t, err)
	e, _ := NewEnv(&RawConfig{Envs: envs}, MSS{}, badWriter{})
	e.List()
}

func TestEnv_PrintRaw(t *testing.T) {
	writer := &bytes.Buffer{}
	envs, _, err := BuiltinDefault()
	ok(t, err)
	e, _ := NewEnv(&RawConfig{Envs: envs}, MSS{}, writer)
	err = e.PrintRaw()
	ok(t, err)
	want := envs + "\n"
	got := writer.String()
	equals(t, want, got)
}

func TestEnv_PrintRawBadWriter(t *testing.T) {
	envs, _, err := BuiltinDefault()
	ok(t, err)
	e, err := NewEnv(&RawConfig{Envs: envs}, MSS{}, badWriter{})
	ok(t, err)
	e.PrintRaw()
}
