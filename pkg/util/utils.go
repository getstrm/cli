package util

import (
	"github.com/bykof/gostradamus"
	"github.com/lithammer/dedent"
	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/timestamppb"
	"os"
	"pace/pace/pkg/common"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func GetStringAndErr(f *pflag.FlagSet, k string) string {
	v, err := f.GetString(k)
	common.CliExit(err)
	return v
}
func GetBoolAndErr(f *pflag.FlagSet, k string) bool {
	v, err := f.GetBool(k)
	common.CliExit(err)
	return v
}

func CreateConfigDirAndFileIfNotExists() {
	err := os.MkdirAll(filepath.Dir(common.ConfigPath()), 0700)
	common.CliExit(err)

	configFilepath := path.Join(common.ConfigPath(), common.DefaultConfigFilename+common.DefaultConfigFileSuffix)

	if _, err := os.Stat(configFilepath); os.IsNotExist(err) {
		writeFileError := os.WriteFile(
			configFilepath,
			common.DefaultConfigFileContents,
			0644,
		)

		common.CliExit(writeFileError)
	}
}

// LongDocs dedents, trims surrounding whitespace, changes !pace for the command Name and changes ° for `
func LongDocs(s string) string {
	s2 := DedentTrim(strings.Replace(
		strings.Replace(s, "!pace", common.RootCommandName, -1), "°", "`", -1))
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
