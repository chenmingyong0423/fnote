// Copyright 2023 chenmingyong0423

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"embed"
	"flag"
	"fmt"
	"github.com/chenmingyong0423/gkit/stringx"
	"os"
	"text/template"
)

var (
	domain    = flag.String("domain", "", "the name of domain, e.g. MessageTemplate, must be camel case")
	outputDir = flag.String("output", "", "the output directory, e.g. internal/message_template")

	//go:embed templates/handler.tmpl
	handler embed.FS

	//go:embed templates/service.tmpl
	service embed.FS

	//go:embed templates/repository.tmpl
	repository embed.FS

	//go:embed templates/dao.tmpl
	dao embed.FS
)

type GenDomain struct {
	DomainName    string
	UnderlineName string
}

func main() {
	flag.Parse()
	if domain == nil || *domain == "" {
		panic("domain name is empty")
	}
	gen := GenDomain{
		DomainName:    *domain,
		UnderlineName: stringx.CamelToSnake(*domain),
	}

	if outputDir == nil || *outputDir == "" {
		outputDir = new(string)
		*outputDir = "internal/" + gen.UnderlineName
	}

	err := executeTemplate(handler, "templates/handler.tmpl", *outputDir+"/handler", fmt.Sprintf("/%s_handler.go", gen.UnderlineName), gen)
	if err != nil {
		_ = os.RemoveAll(*outputDir)
		panic(err)
	}

	err = executeTemplate(service, "templates/service.tmpl", *outputDir+"/service", fmt.Sprintf("/%s_service.go", gen.UnderlineName), gen)
	if err != nil {
		_ = os.RemoveAll(*outputDir)
		panic(err)
	}

	err = executeTemplate(repository, "templates/repository.tmpl", *outputDir+"/repository", fmt.Sprintf("/%s_repository.go", gen.UnderlineName), gen)
	if err != nil {
		_ = os.RemoveAll(*outputDir)
		panic(err)
	}

	err = executeTemplate(dao, "templates/dao.tmpl", *outputDir+"/repository/dao", fmt.Sprintf("/%s_dao.go", gen.UnderlineName), gen)
	if err != nil {
		_ = os.RemoveAll(*outputDir)
		panic(err)
	}
	_, _ = fmt.Fprintf(os.Stdout, "generate domain %s success\n", *domain)
}

func executeTemplate(fs embed.FS, templatePath string, outputDir string, output string, gen GenDomain) error {
	tpl, err := template.ParseFS(fs, templatePath)
	if err != nil {
		return err
	}
	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return err
	}
	created, err := os.Create(outputDir + output)
	if err != nil {
		return err
	}
	err = tpl.Execute(created, gen)
	if err != nil {
		return err
	}
	return nil
}
