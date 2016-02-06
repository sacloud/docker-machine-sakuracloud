package cli

import (
	"fmt"
	"github.com/docker/machine/commands/mcndirs"
	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/yamamoto-febc/docker-machine-sakuracloud/lib/api"
	"github.com/yamamoto-febc/docker-machine-sakuracloud/lib/option"
	"github.com/yamamoto-febc/docker-machine-sakuracloud/lib/persist"
	sakura "github.com/yamamoto-febc/docker-machine-sakuracloud/spec"
	"os"
	"strconv"
	"strings"
)

const (
	// defaultConfigName default setting name
	defaultConfigName = "default"
)

// API Type of Cli Options Handling API interface
type API interface {
	GetDriversDir() string
	GetName() string
	GetConfigValue(string) (*Config, error)
	SetConfigValue(string, string) error
	ListConfigValue() ([]*Config, error)
	ListOptions() []option.Option
	GetOption(string) *option.Option
	ClearConfigValue(string) error
	GetClient() *api.Client
	GetDriverOptions(drivers.DriverOptions) drivers.DriverOptions
	persist.Store
}

// Client is implements CliOptions
type client struct {
	targetSettingName string
	config            *option.Config
	*persist.Filestore
}

// GetName return target setting name
func (c *client) GetName() string {
	return c.targetSettingName
}

// GetDriversDir return drivers dir path
func (c *client) GetDriversDir() string {
	return c.Filestore.GetDriversDir()
}

// NewClient create CliOptions instance
func NewClient() API {
	return &client{
		targetSettingName: defaultConfigName,
		Filestore:         persist.NewFilestore(mcndirs.GetBaseDir()),
	}
}

func (c *client) GetConfigValue(key string) (*Config, error) {
	opt, err := c.getCliConfig(key)
	return opt, err
}

func (c *client) SetConfigValue(key string, value string) error {
	opt, err := c.getCliConfig(key)
	if err != nil {
		return err
	}

	if opt.CurrentValue = value; !opt.IsDefault() {
		c.config.Set(key, value)
		if err := c.Filestore.Save(c.config); err != nil {
			return err
		}
	}

	return nil
}

func (c *client) ListConfigValue() ([]*Config, error) {
	opts := sakura.Options.McnFlags()
	var options []*Config
	for i := range opts {
		o, err := c.getCliConfigWithFlag(opts[i])
		if err != nil {
			return nil, err
		}
		if o.EnvName != "" {
			options = append(options, o)
		}
	}
	return options, nil
}

func (c *client) ListOptions() []option.Option {
	return sakura.Options.CliOptions()
}

func (c *client) GetOption(key string) *option.Option {
	return sakura.Options.GetCliOption(key)
}

func (c *client) ClearConfigValue(key string) error {
	opt, err := c.getCliConfig(key)
	if err != nil {
		return err
	}

	if !opt.IsDefault() {
		c.config.Clear(key)
		if err := c.Filestore.Save(c.config); err != nil {
			return err
		}
	}

	return nil
}

func (c *client) GetClient() *api.Client {
	token, _ := c.GetConfigValue("access-token")
	secret, _ := c.GetConfigValue("access-token-secret")
	region, _ := c.GetConfigValue("region")

	if token.CurrentValue != "" && secret.CurrentValue != "" && region.CurrentValue != "" {
		return api.NewClient(token.CurrentValue, secret.CurrentValue, region.CurrentValue)
	}
	return nil
}

func (c *client) loadOption(key string) error {
	exists, err := c.Filestore.Exists(c.GetName())
	if err != nil {
		return err
	}

	if exists {
		conf, err := c.Filestore.Load(c.GetName())
		if err != nil {
			return err
		}
		c.config = conf
	} else {
		conf := c.CreateNewConfig(defaultConfigName)
		c.Filestore.Save(conf)
		c.config = conf
	}

	return nil
}

func (c *client) getCliConfig(key string) (*Config, error) {
	keyName := c.getOptionFullName(key)
	options := sakura.Options.McnFlags()
	var opt mcnflag.Flag
	for i := range options {
		if options[i].String() == keyName {
			opt = options[i]
			break
		}
	}
	if opt == nil {
		return nil, nil //not found (however has no error)
	}
	return c.getCliConfigWithFlag(opt)
}

func (c *client) getCliConfigWithFlag(f mcnflag.Flag) (*Config, error) {

	key := c.getOptionName(f.String())
	var envName string
	switch t := f.(type) {
	case mcnflag.BoolFlag:
		envName = t.EnvVar
	case mcnflag.IntFlag:
		envName = t.EnvVar
	case mcnflag.StringFlag:
		envName = t.EnvVar
	case mcnflag.StringSliceFlag:
		envName = t.EnvVar
	default:
		envName = ""
	}

	def := f.Default()
	var defValue string
	if def == nil {
		defValue = ""
	} else {
		switch t := def.(type) {
		case bool:
			defValue = strconv.FormatBool(t)
		case string:
			defValue = t
		case []string:
			defValue = strings.Join(t, ",")
		case int:
			defValue = strconv.Itoa(t)
		default:
			return nil, fmt.Errorf("parse error (default value type is unknown) : %v", f)
		}
	}

	c.loadOption(c.GetName())
	conf := c.config
	value, _ := conf.Get(key)

	if envVal := os.Getenv(envName); envVal != "" {
		value = envVal
	}

	// if file and env is not set , use default.
	if value == "" {
		value = defValue
	}

	return &Config{
		KeyName:      key,
		KeyFullName:  f.String(),
		CurrentValue: value,
		DefaultValue: defValue,
		EnvName:      envName,
	}, nil
}

func (c *client) getOptionFullName(name string) string {
	return fmt.Sprintf("sakuracloud-%s", name)
}

func (c *client) getOptionName(fullName string) string {
	return strings.Replace(fullName, "sakuracloud-", "", 1)
}

func (c *client) GetDriverOptions(flags drivers.DriverOptions) drivers.DriverOptions {
	return wrapDriverOptions{
		values: flags,
		client: c,
	}
}

type wrapDriverOptions struct {
	values drivers.DriverOptions
	client *client
}

func (r wrapDriverOptions) String(key string) string {
	optionName := r.client.getOptionName(key)
	config, _ := r.client.GetConfigValue(optionName)
	if config == nil || config.IsDefault() {
		return r.values.String(key)
	}
	return config.CurrentValue
}

func (r wrapDriverOptions) StringSlice(key string) []string {
	return r.values.StringSlice(key)
}

func (r wrapDriverOptions) Int(key string) int {
	optionName := r.client.getOptionName(key)
	config, _ := r.client.GetConfigValue(optionName)
	if config == nil || config.IsDefault() {
		return r.values.Int(key)
	}
	ret, _ := strconv.Atoi(config.CurrentValue)
	return ret
}

func (r wrapDriverOptions) Bool(key string) bool {
	optionName := r.client.getOptionName(key)
	config, _ := r.client.GetConfigValue(optionName)
	if config == nil || config.IsDefault() {
		return r.values.Bool(key)
	}
	ret, _ := strconv.ParseBool(config.CurrentValue)
	return ret
}
