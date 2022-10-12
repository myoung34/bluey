package emitters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-kit/log/level"
	"github.com/myoung34/bluey/bluey"
	"strconv"
	"text/template"
)

type Stdout struct {
	Enabled       bool
	DeviceUUIDMap string `json:"device_uuids"`
	Template      string
}

func StdoutEmit(payload blue.Payload, emitterConfig interface{}) (string, error) {
	stdout := Stdout{}
	jsonString, _ := json.Marshal(emitterConfig)
	json.Unmarshal(jsonString, &stdout)

	level.Info(blue.Logger).Log("emitters.stdout.enabled", fmt.Sprintf("%v", stdout.Enabled))
	level.Info(blue.Logger).Log("emitters.stdout.template", fmt.Sprintf("%+v", stdout.Template))

	nickname := getNickNameFromUUIDs(stdout.DeviceUUIDMap, payload.ID)

	bodyTemplate := Template{
		Nickname:  nickname,
		Minor:     strconv.Itoa(int(payload.Minor)),
		Mac:       payload.Mac,
		Major:     strconv.Itoa(int(payload.Major)),
		Timestamp: payload.Timestamp,
	}
	level.Info(blue.Logger).Log("emitters.stdout.payload", fmt.Sprintf("%+v", payload))

	tmpl, err := template.New("stdout").Parse(`{"name": "Device {{.Nickname}}", "Minor": {{.Minor}}, "Major": {{.Major}}}`)
	if len(stdout.Template) > 0 {
		tmpl, err = template.New("stdout").Parse(stdout.Template)
	}
	if err != nil {
		level.Error(blue.Logger).Log("emitters.stdout", err)
		return "", err
	}
	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, bodyTemplate); err != nil {
		level.Error(blue.Logger).Log("emitters.stdout", err)
		return "", err
	}

	level.Info(blue.Logger).Log("emitters.stdout.rendered_template", fmt.Sprintf("%+v", tpl.String()))
	return tpl.String(), nil
}
