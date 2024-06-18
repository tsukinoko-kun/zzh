package config

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"google.golang.org/protobuf/proto"
)

func getConfigFolder() (string, error) {
	path := ""
	if xdgConfig, ok := os.LookupEnv("XDG_CONFIG_HOME"); ok {
		path = filepath.Join(xdgConfig, "zzh")
	} else if appdata, ok := os.LookupEnv("APPDATA"); runtime.GOOS == "windows" && ok {
		path = filepath.Join(appdata, "zzh")
	} else {
		if home, err := os.UserHomeDir(); err != nil {
			return "", errors.Join(errors.New("failed to get usere home dir"), err)
		} else {
			path = filepath.Join(home, ".zzh")
		}
	}

	return path, os.MkdirAll(path, 0777)
}

func GetConfigs() ([]*Config, error) {
	path, err := getConfigFolder()
	if err != nil {
		return nil, errors.Join(errors.New("failed to get zzh config dir"), err)
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to read content of config dir %s", path), err)
	}

	configs := make([]*Config, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if filepath.Ext(file.Name()) != ".zzh" {
			continue
		}
		filePath := filepath.Join(path, file.Name())
		in, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}
		config := &Config{}
		if err := proto.Unmarshal(in, config); err != nil {
			continue
		}
		configs = append(configs, config)
	}

	return configs, nil
}

func New() (*Config, error) {
	path, err := getConfigFolder()
	if err != nil {
		return nil, errors.Join(errors.New("failed to get zzh config dir"), err)
	}

	config := &Config{}

	config.Port = "22"

	f := huh.NewForm(huh.NewGroup(
		huh.NewInput().Title("user").Value(&config.User).Validate(
			func(s string) error {
				if len(s) == 0 {
					return errors.New("user must be none empty")
				}
				return nil
			}),
		huh.NewInput().Title("host").Value(&config.Host).Validate(
			func(s string) error {
				if len(s) == 0 {
					return errors.New("user must be none empty")
				}
				return nil
			}),
		huh.NewInput().Title("passwort").EchoMode(huh.EchoModePassword).Value(&config.Password),
		huh.NewInput().Title("port").Validate(func(s string) error {
			if _, err := strconv.Atoi(s); err != nil {
				return errors.Join(fmt.Errorf("invalid port %q", s), err)
			} else {
				return nil
			}
		}).Value(&config.Port),
	))

	if err := f.Run(); err != nil {
		return nil, err
	}

	if out, err := proto.Marshal(config); err != nil {
		return nil, errors.Join(errors.New("failed to marshal config"), err)
	} else {
		fileName := base64.StdEncoding.EncodeToString([]byte(config.Display())) + ".zzh"
		filePath := filepath.Join(path, fileName)
		if err := os.WriteFile(filePath, out, 0666); err != nil {
			return nil, errors.Join(fmt.Errorf("failed to write config to file %s", filePath), err)
		}
	}

	return config, nil
}

func (config *Config) Display() string {
	return fmt.Sprintf("%s@%s:%s", config.User, config.Host, config.Port)
}

func (a *Config) Compare(b *Config) int {
	return strings.Compare(a.Display(), b.Display())
}

func Select(configs []*Config) (*Config, error) {
	options := make([]huh.Option[*Config], 0, len(configs))
	var config *Config

	for _, config := range configs {
		options = append(options, huh.Option[*Config]{Key: config.Display(), Value: config})
	}

	f := huh.NewForm(huh.NewGroup(
		huh.NewSelect[*Config]().Value(&config).Options(options...),
	))

	return config, f.Run()
}
