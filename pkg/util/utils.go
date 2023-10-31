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

// LongDocs dedents, trims surrounding whitespace, changes !pace for the command Name and changes ° for `
func LongDocs(s string) string {
	s2 := DedentTrim(strings.Replace(
		strings.Replace(s, "!pace", RootCommandName, -1), "°", "`", -1))
	return s2
}

func LongDocsUsage(s string) string {
	return LongDocs(s) + "\n\n### Usage"
}

func DedentTrim(s string) string {
	return strings.TrimSpace(dedent.Dedent(s))

}

func IsoFormat(tz gostradamus.Timezone, t *timestamppb.Timestamp) string {
	tt := time.Unix(t.Seconds, int64(t.Nanos))
	n := gostradamus.DateTimeFromTime(tt)
	return n.InTimezone(tz).IsoFormatTZ()
}

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
