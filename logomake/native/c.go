package native

import (
	"github.com/mndrix/golog"
	"github.com/mndrix/golog/native"
	"github.com/mndrix/golog/term"
)

func C2(m golog.Machine, args []term.Term) {
	var (
		src string
		dst string
	)
	d := native.NewDecoder(m)
	if _, err := d.Decode(args[0], src); err != nil {
		return golog.ForeignFail()
	}
	if _, err := d.Decode(args[1], dst); err != nil {
		return golog.ForeignFail()
	}

}
