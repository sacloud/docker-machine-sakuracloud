package cli

import (
	"fmt"
	"os"
)

const (
	nilValueString     = "(empty)"
	defaultValueString = "[default]"
	fromNone           = "(Default)"
	fromEnv            = "Environment"
	fromFile           = "File"
)

// Config value of CliOptions
type Config struct {
	KeyName      string
	KeyFullName  string
	CurrentValue string
	DefaultValue string
	EnvName      string
}

// IsDefault return CurrentValue is same DefaultValue
func (c *Config) IsDefault() bool {
	return c.CurrentValue == c.DefaultValue
}

// IsFromEnv return valur from env
func (c *Config) IsFromEnv() bool {
	if c.DefaultValue == c.CurrentValue {
		return false
	}
	if os.Getenv(c.EnvName) != "" {
		return true
	}

	return false
}

// FormatedCurrentValue return FormatedCurrentValue
func (c *Config) FormatedCurrentValue() string {
	strCurrent := c.CurrentValue
	if strCurrent == "" {
		strCurrent = nilValueString
	}
	strDefault := ""
	if c.IsDefault() {
		strDefault = defaultValueString
	}
	return fmt.Sprintf("%s %s", strCurrent, strDefault)
}

// ValueFrom return value from
func (c *Config) ValueFrom() string {
	if c.DefaultValue == c.CurrentValue {
		return fromNone
	}

	if os.Getenv(c.EnvName) != "" {
		return fromEnv
	}

	return fromFile
}

// GetPrintInfo return print value info
func (c *Config) GetPrintInfo(detail bool) []string {

	strCurrent := c.CurrentValue
	if strCurrent == "" {
		strCurrent = nilValueString
	}

	fullName := fmt.Sprintf("--%s", c.KeyFullName)
	envName := fmt.Sprintf("$%s", c.EnvName)

	if detail {
		return []string{
			c.KeyName,
			fullName,
			envName,
			c.ValueFrom(),
			strCurrent,
		}
	}

	return []string{
		c.KeyName,
		c.ValueFrom(),
		strCurrent,
	}

}

// GetPrintHeader return print header info
func GetPrintHeader(detail bool) []string {
	if detail {
		return []string{"Name", "FULL_NAME", "ENVIRONMENT_NAME", "VALUE_FROM", "CURRENT_SETTING"}
	}
	return []string{"NAME", "VALUE_FROM", "CURRENT_SETTING"}
}
