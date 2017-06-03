package native

import (
	"path/filepath"

	"github.com/mndrix/golog"
	"github.com/mndrix/golog/native"
	"github.com/mndrix/golog/term"
	"github.com/wvxvw/logomake/logging"
)

func Glob2(m golog.Machine, args []term.Term) golog.ForeignReturn {
	var pattern string
	d := native.NewDecoder(m)
	if _, err := d.Decode(args[0], &pattern); err != nil {
		logging.Dbug.Printf("Couldn't decode %s, %s", args[0], err)
		return golog.ForeignFail()
	}
	files, err := filepath.Glob(pattern)
	if err != nil {
		logging.Dbug.Printf("filepath.Glob failed %s, %s", args[0], err)
		return golog.ForeignFail()
	}
	logging.Dbug.Printf("Found files %s using pattern %s", files, pattern)
	e := native.NewEncoder()
	return golog.ForeignUnify(args[1], e.Encode(files))
}
