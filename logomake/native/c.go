package native

import (
	"github.com/mndrix/golog"
	"github.com/mndrix/golog/native"
	"github.com/mndrix/golog/term"

	"github.com/wvxvw/logomake/c"
	"github.com/wvxvw/logomake/logging"
)

func C2(m golog.Machine, args []term.Term) golog.ForeignReturn {
	var (
		src []string
		dst string
	)
	d := native.NewDecoder(m)
	if _, err := d.Decode(args[0], &src); err != nil {
		logging.Dbug.Printf("Couldn't decode %s, %s", args[0], err)
		return golog.ForeignFail()
	}
	if _, err := d.Decode(args[1], &dst); err != nil {
		logging.Dbug.Printf("Couldn't decode %s, %s", args[1], err)
		return golog.ForeignFail()
	}
	if err := c.CompileCProgram(src, dst, nil); err != nil {
		logging.Dbug.Printf(
			"Couldn't compile %s, into %s: %s",
			src, dst, err,
		)
		return golog.ForeignFail()
	}
	return golog.ForeignTrue()
}
