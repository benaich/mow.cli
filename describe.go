package cli

import (
	"github.com/jawher/mow.cli/internal/container"
	"reflect"
	"strings"
)

// Parameter describes a declared option or positional argument, before parsing.
type Parameter struct {
	Names       []string    `json:"names"`
	Description string      `json:"description"`
	Required    bool        `json:"required"`
	Enum        []string    `json:"enum,omitempty"`
	Type        string      `json:"type"`
	Default     interface{} `json:"default,omitempty"`
	Repeatable  bool        `json:"repeatable,omitempty"`
}

// Description contains declaration metadata, never actions or environment values.
type Description struct {
	Path        string        `json:"path"`
	Aliases     []string      `json:"aliases,omitempty"`
	Description string        `json:"description"`
	Usage       string        `json:"usage"`
	Options     []Parameter   `json:"options"`
	Arguments   []Parameter   `json:"arguments"`
	Commands    []Description `json:"commands,omitempty"`
}

// Describe initializes an isolated declaration tree, without parsing or actions.
// Initializers must only declare commands, options and actions.
func (c *Cmd) Describe() (Description, error) {
	clone := *c
	clone.options = append([]*container.Container(nil), c.options...)
	clone.args = append([]*container.Container(nil), c.args...)
	clone.optionsIdx = map[string]*container.Container{}
	for key, value := range c.optionsIdx {
		clone.optionsIdx[key] = value
	}
	clone.argsIdx = map[string]*container.Container{}
	for key, value := range c.argsIdx {
		clone.argsIdx[key] = value
	}
	clone.commands = make([]*Cmd, len(c.commands))
	for i, child := range c.commands {
		copy := *child
		clone.commands[i] = &copy
	}
	if err := clone.doInit(); err != nil {
		return Description{}, err
	}
	d := Description{Path: strings.Join(append(clone.parents, clone.name), " "), Aliases: clone.aliases, Description: clone.desc, Usage: clone.Spec, Options: []Parameter{}, Arguments: []Parameter{}}
	for _, p := range clone.options {
		param := describeParameter(p)
		param.Required = requiredParameter(clone.Spec, param.Names)
		d.Options = append(d.Options, param)
	}
	for _, p := range clone.args {
		param := describeParameter(p)
		param.Required = requiredParameter(clone.Spec, param.Names)
		d.Arguments = append(d.Arguments, param)
	}
	for _, child := range clone.commands {
		if child.Hidden {
			continue
		}
		sub, err := child.Describe()
		if err != nil {
			return Description{}, err
		}
		d.Commands = append(d.Commands, sub)
	}
	return d, nil
}

func describeParameter(p *container.Container) Parameter {
	d := Parameter{Names: p.Names, Description: p.Desc, Type: "string"}
	if len(d.Names) == 0 {
		d.Names = []string{p.Name}
	}
	v := reflect.ValueOf(p.Value)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Bool:
		d.Type = "boolean"
	case reflect.Int, reflect.Int64:
		d.Type = "integer"
	case reflect.Float64:
		d.Type = "number"
	case reflect.Slice:
		d.Type = "array"
		d.Repeatable = true
	}
	if !p.HideValue && v.IsValid() && v.CanInterface() {
		d.Default = v.Interface()
	}
	return d
}

// Required means unconditional in the usage grammar. Alternative groups and
// conditional requirements remain expressed by Usage rather than flattened.
func requiredParameter(spec string, names []string) bool {
	depth := 0
	var plain strings.Builder
	for _, ch := range spec {
		switch ch {
		case '[', '(':
			depth++
			plain.WriteByte(' ')
		case ']', ')':
			depth--
			plain.WriteByte(' ')
		default:
			if depth == 0 {
				plain.WriteRune(ch)
			}
		}
	}
	for _, specWord := range strings.Fields(plain.String()) {
		specWord = strings.TrimSuffix(specWord, "...")
		for _, name := range names {
			if name == specWord {
				return true
			}
		}
	}
	return false
}
