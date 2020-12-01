// Copyright 2020 Paul Greenberg greenpau@outlook.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"github.com/greenpau/go-vrop"
	"github.com/greenpau/versioned"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

var (
	app        *versioned.PackageManager
	appVersion string
	gitBranch  string
	gitCommit  string
	buildUser  string
	buildDate  string
)

func init() {
	app = versioned.NewPackageManager("vropcli")
	app.Description = "vRealize API Client"
	app.Documentation = "https://github.com/greenpau/go-vrop/"
	app.SetVersion(appVersion, "1.0.0")
	app.SetGitBranch(gitBranch, "main")
	app.SetGitCommit(gitCommit, "4768e50")
	app.SetBuildUser(buildUser, "")
	app.SetBuildDate(buildDate, "")
}

func main() {
	var logLevel string
	var isShowVersion bool
	var configDir string
	var configFile string
	var host, username, password string
	var getVirtualMachines bool

	flag.StringVar(&configFile, "config", "", "configuration file")
	flag.StringVar(&host, "host", "", "vRealize Operations Manager Hostname")
	flag.StringVar(&username, "username", "", "Username")
	flag.StringVar(&password, "password", "", "Password")

	flag.BoolVar(&getVirtualMachines, "get-virtual-machines", false, "Get virtual machines")

	flag.StringVar(&logLevel, "log-level", "info", "logging severity level")
	flag.BoolVar(&isShowVersion, "version", false, "show version")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n%s - %s\n\n", app.Name, app.Description)
		fmt.Fprintf(os.Stderr, "Usage: %s [arguments]\n\n", app.Name)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nDocumentation: %s\n\n", app.Documentation)
	}
	flag.Parse()

	if isShowVersion {
		fmt.Fprintf(os.Stdout, "%s\n", app.Banner())
		os.Exit(0)
	}

	// Determine configuration file name and extension
	if configFile == "" {
		configDir = "."
		configFile = app.Name + ".yaml"
	} else {
		configDir, configFile = filepath.Split(configFile)
	}
	configFileExt := filepath.Ext(configFile)
	if configFileExt == "" {
		fmt.Fprintf(os.Stderr, "--config specifies a file without an extension, e.g. .yaml or .json\n")
		os.Exit(1)
	}

	configName := strings.TrimSuffix(configFile, configFileExt)
	viper.SetConfigName(configName)
	viper.SetEnvPrefix("vrop")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_", " ", "_"))
	viper.AddConfigPath("$HOME/.config/" + app.Name)
	viper.AddConfigPath(configDir)
	viper.AutomaticEnv()

	// Obtain settings via environment variable

	if host == "" {
		if v := viper.Get("host"); v != nil {
			host = viper.Get("host").(string)
		}
	}

	if username == "" {
		if v := viper.Get("username"); v != nil {
			username = viper.Get("username").(string)
		}
	}

	if password == "" {
		if v := viper.Get("password"); v != nil {
			password = viper.Get("password").(string)
		}
	}

	// Obtain settings via configuration file
	if err := viper.ReadInConfig(); err == nil {
		if host == "" {
			if v := viper.Get("host"); v != nil {
				host = viper.Get("host").(string)
			}
		}
		if username == "" {
			if v := viper.Get("username"); v != nil {
				username = viper.Get("username").(string)
			}
		}

		if password == "" {
			if v := viper.Get("password"); v != nil {
				password = viper.Get("password").(string)
			}
		}
	} else {
		if !strings.Contains(err.Error(), "Not Found in") {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	}

	opts := make(map[string]interface{})
	opts["log_level"] = logLevel
	cli, err := vrop.NewClient(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	defer cli.Close()

	if err := cli.SetHost(host); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	if err := cli.SetUsername(username); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	if err := cli.SetPassword(password); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	cli.Info()

	opts = make(map[string]interface{})
	if getVirtualMachines {
		items, err := cli.GetVirtualMachines(opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		for _, item := range items {
			s, err := item.ToJSONString()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				continue
			}
			fmt.Fprintf(os.Stdout, "%s\n", s)
		}
		os.Exit(0)
	}

	fmt.Fprintf(os.Stderr, "actionable argument is missing\n")
	os.Exit(1)
}
