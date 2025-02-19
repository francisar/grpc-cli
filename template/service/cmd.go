// Copyright © 2020 The gRPC Kit Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

func (t *templateService) fileDirectoryCmd() {
	t.files = append(t.files, &templateFile{
		name:  "cmd/server/main.go",
		parse: true,
		body: `
// Code generated by "grpc-kit-cli/{{ .Global.ReleaseVersion }}". DO NOT EDIT.

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/grpc-kit/pkg/cfg"
	"github.com/grpc-kit/pkg/signal"
	"github.com/grpc-kit/pkg/version"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"

	"{{ .Global.Repository }}/handler"
)

var (
	flagCfgFile *string
	flagVersion *bool
)

func init() {
	flagCfgFile = flag.StringP("config", "c", "./config/app.yaml", "config file")
	flagVersion = flag.BoolP("version", "v", false, "print version and exit")
	flag.Parse()

	if *flagVersion {
		fmt.Println(version.Get())
		os.Exit(0)
	}
}

func main() {
	viper.SetConfigFile(*flagCfgFile)

	if err := viper.ReadInConfig(); err == nil {
		fmt.Printf("Using config file: %v\n", *flagCfgFile)
	} else {
		fmt.Printf("Load config file: %v, err: %v\n", *flagCfgFile, err)
		os.Exit(1)
	}

	ctx := context.Background()

	m, err := startServer(ctx)
	if err != nil {
		fmt.Printf("Start server err: %v\n", err)
		os.Exit(1)
	}

	signal.WaitQuit()

	if err := m.Shutdown(ctx); err != nil {
		fmt.Printf("Shutdown server err: %v\n", err)
		os.Exit(1)
	}
}

func startServer(ctx context.Context) (*handler.Microservice, error) {
	lc, err := cfg.New(viper.GetViper())
	if err != nil {
		return nil, err
	}

	m, err := handler.NewMicroservice(lc)
	if err != nil {
		return nil, err
	}

	if err := m.Register(ctx); err != nil {
		return nil, err
	}

	return m, nil
}
`,
	})
}
