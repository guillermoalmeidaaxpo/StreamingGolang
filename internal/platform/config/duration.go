package config

import (
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

type duration time.Duration

func (d *duration) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {
	case yaml.ScalarNode:
		var raw string
		if err := value.Decode(&raw); err != nil {
			return err
		}
		if raw == "" {
			*d = 0
			return nil
		}
		parsed, err := time.ParseDuration(raw)
		if err == nil {
			*d = duration(parsed)
			return nil
		}
		seconds, err := strconv.Atoi(raw)
		if err != nil {
			return err
		}
		*d = duration(time.Duration(seconds) * time.Second)
		return nil
	default:
		var seconds int
		if err := value.Decode(&seconds); err != nil {
			return err
		}
		*d = duration(time.Duration(seconds) * time.Second)
		return nil
	}
}

func (d duration) Duration() time.Duration {
	return time.Duration(d)
}
