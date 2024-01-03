package util

import (
	"bytes"
	"encoding/json"
	"github.com/bykof/gostradamus"
	"github.com/lithammer/dedent"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"sigs.k8s.io/yaml"
	"strings"
	"time"
)

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

/*
Example
example first lines should start with a tab to make a nice help layout.
*/
func Example(s string) string {
	return "    " + LongDocs(s)
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
