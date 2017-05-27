package c

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

var (
	programSrc = []string{
		"hello_world.c",
		`
#include<stdio.h>

int main() {
    printf("Hello World");
    return 0;
}
`,
	}
	programDst = "hello_world"
)

func setupSources() {
	err := ioutil.WriteFile(
		programSrc[0],
		[]byte(programSrc[1]),
		os.ModePerm,
	)
	if err != nil {
		panic(fmt.Sprintf("Test fixture failure: %s", err))
	}
}

func teardownSources() {
	_ = os.Remove(programSrc[0])
	_ = os.Remove(programDst)
}

func TestMain(m *testing.M) {
	setupSources()
	ret := m.Run()
	teardownSources()
	os.Exit(ret)
}

func TestCOptionsString(t *testing.T) {
	defer func() {
		if p := recover(); p != nil {
			t.Fatalf("Paniced: %s", p)
		}
	}()
	opts := &COptions{}
	res := opts.String()
	if res == "" {
		opts.CC = "gcc"
	}
	res = opts.String()
	if res == "" {
		t.Fatalf("CC not set")
	}
}

func TestCompile(t *testing.T) {
	if cc := os.Getenv("CC"); cc == "" {
		os.Setenv("CC", "gcc")
	}
	err := CompileCProgram([]string{programSrc[0]}, programDst, nil)
	if err != nil {
		t.Fatalf("Compiling C program failed: %s", err)
	}
	expected := []byte("Hello World")
	out, err := exec.Command("./" + programDst).CombinedOutput()
	if err != nil {
		t.Fatalf("Compiled C program shouldn't fail: %s", err)
	}

	if bytes.Compare(out, expected) != 0 {
		t.Fatalf("Expected to produce '%s', but got '%s'", expected, out)
	}
}
