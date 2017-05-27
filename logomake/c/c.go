package c

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"reflect"
)

type COptions struct {
	CC string `default:"$CC"`
}

func (co *COptions) Defaults() *COptions {
	st := reflect.ValueOf(co).Elem()
	result := reflect.New(st.Type()).Elem()
	for i := 0; i < st.NumField(); i++ {
		f := result.Field(i)
		tf := st.Type().Field(i)
		if tf.PkgPath == "" {
			tag := tf.Tag.Get("default")
			var val string
			if tag != "" {
				val = co.loadDefault(tag)
			}
			f.SetString(val)
		}
	}
	tresult := result.Interface().(COptions)
	return &tresult
}

func (co *COptions) Cmd(src []string, prog string) *exec.Cmd {
	args := append(src, "-o", prog)
	return exec.Command(co.CC, args...)
}

func (co *COptions) AreValid() bool {
	return co.CC != ""
}

func (co *COptions) String() string {
	st := reflect.ValueOf(co).Elem()
	buf := bytes.NewBuffer([]byte{})
	for i := 0; i < st.NumField(); i++ {
		f := st.Field(i)
		tf := st.Type().Field(i)
		if tf.PkgPath != "" {
			continue
		}
		val := f.Interface().(string)
		tag := tf.Tag.Get("default")
		if tag == "" {
			continue
		}
		if tag[0] == '$' {
			buf.WriteString(tag[1:])
		} else {
			buf.WriteString(tag)
		}
		buf.WriteByte('=')
		buf.WriteString(val)
	}
	return buf.String()
}

func (co *COptions) loadDefault(tag string) string {
	if tag[0] == '$' {
		return os.Getenv(tag[1:])
	}
	return tag
}

func CompileCProgram(src []string, prog string, opts *COptions) error {
	if opts == nil {
		opts = (&COptions{}).Defaults()
	}
	if !opts.AreValid() {
		return fmt.Errorf("%s -> %s invalid options: %s", src, prog, opts)
	}
	cmd := opts.Cmd(src, prog)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"%s %s -> %s failed: %s: %s",
			opts, src, prog, err, out,
		)
	}
	return nil
}
