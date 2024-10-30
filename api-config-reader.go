package goapiconfigutilis

import (
	"github.com/wlevene/ini"
)

type ApiConfigReader struct {
	configFilePath string
	config         *ini.Ini
}

// Read an ini file
func InitApiConfigReader(configFilePath string) *ApiConfigReader {
	return &ApiConfigReader{
		configFilePath: configFilePath,
		config:         ini.New().LoadFile(configFilePath),
	}
}

// Get a value from the specified section/key
func (c *ApiConfigReader) Get(section, key string) any {
	return c.config.Section(section).Get(key)
}
