package logging

import (
	"log"
	"os"
)

var (
	Dbug = log.New(os.Stdout, "dbug ", log.LstdFlags)
	Info = log.New(os.Stdout, "info ", log.LstdFlags)
	Warn = log.New(os.Stdout, "warn ", log.LstdFlags)
	Eror = log.New(os.Stderr, "eror ", log.LstdFlags)
)
