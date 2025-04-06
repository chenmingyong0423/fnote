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
	"os"
	"text/template"

	"github.com/chenmingyong0423/gkit/stringx"
)

var (
	domain    = flag.String("domain", "", "the name of domain, e.g. MessageTemplate, must be camel case")
	outputDir = flag.String("output", "", "the output directory, e.g. internal/message_template")
	tableName = flag.String("table_name", "", "the table name, e.g. message_template")

	//go:embed templates/request.tmpl
	reqTpl embed.FS

	//go:embed templates/vo.tmpl
	voTpl embed.FS

	//go:embed templates/wire.tmpl
	wireTpl embed.FS

	//go:embed templates/module.tmpl
	module embed.FS

	//go:embed templates/domain.tmpl
	domainTpl embed.FS

	//go:embed templates/web.tmpl
	web embed.FS

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
	TableName     string
	OutputDir     string
}

func main() {
	flag.Parse()
	if domain == nil || *domain == "" {
		panic("domain name is empty")
	}
	gen := GenDomain{
		DomainName:    *domain,
		UnderlineName: stringx.CamelToSnake(*domain),
		TableName:     *tableName,
		OutputDir:     *outputDir,
	}

	if outputDir == nil || *outputDir == "" {
		outputDir = new(string)
		*outputDir = "internal/" + gen.UnderlineName
		gen.OutputDir = *outputDir
	}

	if tableName == nil || *tableName == "" {
		tableName = new(string)
		*tableName = gen.UnderlineName
		gen.TableName = *tableName
	}

	err := executeTemplate(wireTpl, "templates/wire.tmpl", *outputDir, "/wire.go", gen)
	if err != nil {
		_ = os.RemoveAll(*outputDir)
		panic(err)
	}

	err = executeTemplate(module, "templates/module.tmpl", *outputDir, "/module.go", gen)
	if err != nil {
		_ = os.RemoveAll(*outputDir)
		panic(err)
	}

	err = executeTemplate(domainTpl, "templates/domain.tmpl", *outputDir+"/internal/domain", fmt.Sprintf("/%s.go", gen.UnderlineName), gen)
	if err != nil {
		_ = os.RemoveAll(*outputDir)
		panic(err)
	}

	err = executeTemplate(web, "templates/web.tmpl", *outputDir+"/internal/web", fmt.Sprintf("/%s.go", gen.UnderlineName), gen)
	if err != nil {
		_ = os.RemoveAll(*outputDir)
		panic(err)
	}

	err = executeTemplate(reqTpl, "templates/request.tmpl", *outputDir+"/internal/web", "/request.go", gen)
	if err != nil {
		_ = os.RemoveAll(*outputDir)
		panic(err)
	}

	err = executeTemplate(voTpl, "templates/vo.tmpl", *outputDir+"/internal/web", "/vo.go", gen)
	if err != nil {
		_ = os.RemoveAll(*outputDir)
		panic(err)
	}

	err = executeTemplate(service, "templates/service.tmpl", *outputDir+"/internal/service", fmt.Sprintf("/%s.go", gen.UnderlineName), gen)
	if err != nil {
		_ = os.RemoveAll(*outputDir)
		panic(err)
	}

	err = executeTemplate(repository, "templates/repository.tmpl", *outputDir+"/internal/repository", fmt.Sprintf("/%s.go", gen.UnderlineName), gen)
	if err != nil {
		_ = os.RemoveAll(*outputDir)
		panic(err)
	}

	err = executeTemplate(dao, "templates/dao.tmpl", *outputDir+"/internal/repository/dao", fmt.Sprintf("/%s.go", gen.UnderlineName), gen)
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
