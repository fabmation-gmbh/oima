package log

import "github.com/apsdehal/go-logger"
import "os"

var Log, _ = logger.New("oima", 1, os.Stdout)