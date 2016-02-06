package option

import (
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/yamamoto-febc/docker-machine-sakuracloud/lib/api"
	"strings"
)

// List interface of Option holder
type List struct {
	Options []Option
}

// McnFlags return []mcnflag.Flag
func (c *List) McnFlags() []mcnflag.Flag {
	var opts []mcnflag.Flag
	for _, o := range c.Options {
		opts = append(opts, o.RawMcnOption)
	}
	return opts[:]
}

// CliOptions return options(for cli only)
func (c *List) CliOptions() []Option {
	var opts []Option
	for _, o := range c.Options {
		opts = append(opts, o)
	}
	return opts[:]
}

// GetCliOption return option(for cli only)
func (c *List) GetCliOption(key string) *Option {
	for _, o := range c.Options {
		if o.KeyName() == key {
			return &o
		}
	}

	return nil

}

// Option interface of Option
type Option struct {
	// RawMcnOption return mcnflag.Flag
	RawMcnOption mcnflag.Flag
	// CandidateFunc return candidate list for bash complement
	CandidateFunc    func(*api.Client) []string
	UsageStringsFunc func(*api.Client) string
}

// KeyName return cli option keyname
func (o *Option) KeyName() string {
	return strings.Replace(o.RawMcnOption.String(), "sakuracloud-", "", 1)
}

// Description return description
func (o *Option) Description() string {
	switch t := o.RawMcnOption.(type) {
	case mcnflag.StringFlag:
		return t.Usage
	case mcnflag.StringSliceFlag:
		return t.Usage
	case mcnflag.IntFlag:
		return t.Usage
	case mcnflag.BoolFlag:
		return t.Usage
	default:
		return ""
	}

}
