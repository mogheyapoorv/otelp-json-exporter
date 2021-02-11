package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

type Config struct {
	configmodels.ExporterSettings  `mapstructure:",squash"`
	exporterhelper.TimeoutSettings `mapstructure:",squash"`
	exporterhelper.QueueSettings   `mapstructure:"sending_queue"`
	exporterhelper.RetrySettings   `mapstructure:"retry_on_failure"`

	Brokers         []string `mapstructure:"brokers"`
	ProtocolVersion string   `mapstructure:"protocol_version"`
	Topic           string   `mapstructure:"topic"`
	Encoding        string   `mapstructure:"encoding"`
	Metadata        Metadata `mapstructure:"metadata"`
	AllAsFields     bool     `mapstructure:"all"`
	// Dot replacement for tag keys when stored as object fields
	DotReplacement string `mapstructure:"dot_replacement"`
	File           string `mapstructure:"config_file"`
	Include        string `mapstructure:"include"`
}

type Metadata struct {
	Full  bool          `mapstructure:"full"`
	Retry MetadataRetry `mapstructure:"retry"`
}

type MetadataRetry struct {
	Max     int           `mapstructure:"max"`
	Backoff time.Duration `mapstructure:"backoff"`
}

// TagKeysAsFields returns tags from the file and command line merged
func (c *Config) TagKeysAsFields() ([]string, error) {
	var tags []string

	// from file
	if c.File != "" {
		file, err := os.Open(filepath.Clean(c.File))
		if err != nil {
			return nil, err
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if tag := strings.TrimSpace(line); tag != "" {
				tags = append(tags, tag)
			}
		}
		if err := file.Close(); err != nil {
			return nil, err
		}
	}

	// from params
	if c.Include != "" {
		tags = append(tags, strings.Split(c.Include, ",")...)
	}

	return tags, nil
}
