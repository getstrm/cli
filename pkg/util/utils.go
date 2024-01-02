package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bykof/gostradamus"
	"github.com/lithammer/dedent"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"os"
	"runtime"
	"sigs.k8s.io/yaml"
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
	return strings.TrimSpace(dedent.Dedent(replaceContent(s)))
}

func replaceContent(s string) string {
	return strings.Replace(
		strings.Replace(s, "!pace", RootCommandName, -1), "°", "`", -1)
}
func PlainExample(s string) string {
	return "\n" + replaceContent(s)
}

/*
Example
example first lines should start with a tab to make a nice help layout.
*/
func Example(s string) string {
	return "\n    " + LongDocs(s)
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
			var additionalDetails string
			if len(st.Details()) > 0 {
				details := st.Details()[0]
				yamlBytes := ProtoMessageToYaml(details.(proto.Message))
				additionalDetails = string(yamlBytes.Bytes())
			} else {
				additionalDetails = ""
			}
			formattedMessage := fmt.Sprintf(dedentAndTrimMultiline(`
						Error code = %s
						Details = %s

						%s`), (*st).Code(), (*st).Message(), additionalDetails)

			_, _ = fmt.Fprintln(os.Stderr, formattedMessage)
		} else {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}

		os.Exit(1)
	}
}

func dedentAndTrimMultiline(s string) string {
	return strings.TrimLeft(dedent.Dedent(s), "\n")
}

func ProtoMessageToRawJson(proto proto.Message) bytes.Buffer {
	// As protojson.Marshal adds random spaces, we use json.Compact to omit the random spaces in the output.
	// Linked issue in google/protobuf: https://github.com/golang/protobuf/issues/1082
	marshal, _ := protojson.MarshalOptions{
		UseProtoNames: true,
	}.Marshal(proto)
	buffer := bytes.Buffer{}
	_ = json.Compact(&buffer, marshal)
	return buffer
}

func ProtoMessageToPrettyJson(proto proto.Message) bytes.Buffer {
	prettyJson := bytes.Buffer{}
	rawJson := ProtoMessageToRawJson(proto)
	_ = json.Indent(&prettyJson, rawJson.Bytes(), "", "    ")
	return prettyJson
}

func ProtoMessageToYaml(proto proto.Message) bytes.Buffer {
	rawJson := ProtoMessageToRawJson(proto)
	m, _ := yaml.JSONToYAML(rawJson.Bytes())
	return *bytes.NewBuffer(m)
}
