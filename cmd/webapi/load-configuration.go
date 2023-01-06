package main

import (
	"errors"
	"fmt"
	"github.com/ardanlabs/conf"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"time"
)

// webAPIConfiguration describes the web API configuration. This structure is automatically parsed by
// loadConfiguration and values from flags, environment variable or configuration file will be loaded.
type webAPIConfiguration struct {
	Config struct {
		Path string `conf:"default:/conf/config.yml"`
	}
	Web struct {
		APIHost         string        `conf:"default:0.0.0.0:3000"`
		ReadTimeout     time.Duration `conf:"default:60s"`
		WriteTimeout    time.Duration `conf:"default:5s"`
		ShutdownTimeout time.Duration `conf:"default:5s"`
	}
	Log struct {
		Debug    bool   `conf:"default:true"`
		FileName string `conf:"default:-"`
	}
	DB struct {
		Filename string `conf:"default:wasaphoto.db"`
	}
	// Setup user content storage params
	UserContent struct {
		// The local filesystem path where to store data
		FsDir string `conf:"default:static/user_content"`

		// The virtual API path prefix to prepend to request a static file on this server
		WebPrefix string `conf:"default:/static/user_content"`
	}
}

// loadConfiguration creates a webAPIConfiguration starting from flags, environment variables and configuration file.
// It works by loading environment variables first, then update the config using command line flags, finally loading the
// configuration file (specified in WebAPIConfiguration.Config.Path).
// So, CLI parameters will override the environment, and configuration file will override everything.
// Note that the configuration file can be specified only via CLI or environment variable.
func loadConfiguration() (webAPIConfiguration, error) {
	var cfg webAPIConfiguration

	// Try to load configuration from environment variables and command line switches
	if err := conf.Parse(os.Args[1:], "CFG", &cfg); err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			usage, err := conf.Usage("CFG", &cfg)
			if err != nil {
				return cfg, fmt.Errorf("generating config usage: %w", err)
			}
			fmt.Println(usage) //nolint:forbidigo
			return cfg, conf.ErrHelpWanted
		}
		return cfg, fmt.Errorf("parsing config: %w", err)
	}

	// Override values from YAML if specified and if it exists (useful in k8s/compose)
	fp, err := os.Open(cfg.Config.Path)
	if err != nil && !os.IsNotExist(err) {
		return cfg, fmt.Errorf("can't read the config file, while it exists: %w", err)
	} else if err == nil {
		yamlFile, err := io.ReadAll(fp)
		if err != nil {
			return cfg, fmt.Errorf("can't read config file: %w", err)
		}
		err = yaml.Unmarshal(yamlFile, &cfg)
		if err != nil {
			return cfg, fmt.Errorf("can't unmarshal config file: %w", err)
		}
		_ = fp.Close()
	}

	return cfg, nil
}
