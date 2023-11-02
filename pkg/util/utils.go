package util

import (
	"fmt"
	"github.com/bykof/gostradamus"
	"github.com/lithammer/dedent"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"os"
	"runtime"
	"strings"
	"time"
)

func GetStringAndErr(f *pflag.FlagSet, k string) string {
	v, err := f.GetString(k)
	CliExit(err)
	return v
}
func GetBoolAndErr(f *pflag.FlagSet, k string) bool {
	v, err := f.GetBool(k)
	CliExit(err)
	return v
}

/*
LongDocs
dedents, trims surrounding whitespace, changes !pace for the command Name and changes ° for `
The trick with the ° command is to avoid painful string concatenations that collide with
the backticks used in Go raw strings.
*/
func LongDocs(s string) string {
	return strings.TrimSpace(dedent.Dedent(strings.Replace(
		strings.Replace(s, "!pace", RootCommandName, -1), "°", "`", -1)))
}

func IsoFormat(tz gostradamus.Timezone, t *timestamppb.Timestamp) string {
	return gostradamus.DateTimeFromTime(time.Unix(t.Seconds, int64(t.Nanos))).InTimezone(tz).IsoFormatTZ()
}

/*
	RootCommandName

util.RootCommandName is modified in the Makefile to create the `dpace` completion.
If you move or rename this variable, also fix the Makefile (targetVar). If you get this wrong
completion won't work for dpace.
*/
var RootCommandName = "pace"

func CliExit(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		logrus.WithFields(logrus.Fields{"file": file, "line": line}).Error(err)

		st, ok := status.FromError(err)

		if ok {
			_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf(`Error code = %s
Details = %s`, (*st).Code(), (*st).Message()))
		} else {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}

		os.Exit(1)
	}
}
