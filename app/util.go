package app
import (
	"encoding/json"
	"net/url"
	"os/exec"
)

func jsonMarshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func urlEncode(s string) string {
	return url.QueryEscape(s)
}

func execCommand(name string, args ...string) *exec.Cmd {
	return exec.Command(name, args...)
}