// Copyright 2024 chenmingyong0423

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package domain

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"text/template"
)

type Type string

type RecipientType uint

const (
	CommentReceived     Type = "comment"
	UserCommentApproved Type = "user-comment-approval"
	UserCommentRejected Type = "user-comment-disapproval"
	UserCommentReplied  Type = "user-comment-reply"
	FriendApplied       Type = "friend"
	UserFriendApproved  Type = "friend-approval"
	UserFriendRejected  Type = "friend-rejection"
)

const (
	RecipientWebmaster RecipientType = iota
	RecipientUser
)

type Data map[string]any

const (
	VariablePostURL       = "PostURL"
	VariableFriendPageURL = "FriendPageURL"
	VariableReason        = "Reason"
)

type MessageTemplate struct {
	Id            string
	Type          Type
	Name          string
	Title         string
	Content       string
	Active        bool
	IsDefault     bool
	RecipientType RecipientType
	CreatedAt     int64
	UpdatedAt     int64
}

func (mt *MessageTemplate) RenderContent(data Data) error {
	if strings.Contains(mt.Content, "%s") {
		return fmt.Errorf("message template %q contains unsupported legacy placeholders", mt.Type)
	}

	tpl, err := template.New(string(mt.Type)).Option("missingkey=error").Parse(mt.Content)
	if err != nil {
		return fmt.Errorf("parse message template %q: %w", mt.Type, err)
	}

	var rendered bytes.Buffer
	if err = tpl.Execute(&rendered, data); err != nil {
		return fmt.Errorf("render message template %q: %w", mt.Type, err)
	}
	mt.Content = rendered.String()
	return nil
}

func RecipientForType(name Type) (RecipientType, bool) {
	switch name {
	case CommentReceived, FriendApplied:
		return RecipientWebmaster, true
	case UserCommentApproved, UserCommentRejected, UserCommentReplied, UserFriendApproved, UserFriendRejected:
		return RecipientUser, true
	default:
		return 0, false
	}
}

var placeholderPattern = regexp.MustCompile(`^{{\s*\.([A-Za-z][A-Za-z0-9]*)\s*}}$`)
var actionPattern = regexp.MustCompile(`{{[^{}]*}}`)

func RequiredVariables(name Type) []string {
	switch name {
	case UserCommentApproved, UserCommentReplied:
		return []string{VariablePostURL}
	case UserCommentRejected:
		return []string{VariablePostURL, VariableReason}
	case UserFriendApproved:
		return []string{VariableFriendPageURL}
	case UserFriendRejected:
		return []string{VariableFriendPageURL, VariableReason}
	default:
		return []string{}
	}
}

func ValidateContent(name Type, content string) error {
	if strings.TrimSpace(content) == "" {
		return errors.New("message template content cannot be empty")
	}
	if strings.Contains(content, "%s") {
		return errors.New("legacy %s placeholders are not allowed; use named variables")
	}
	if _, err := template.New(string(name)).Parse(content); err != nil {
		return fmt.Errorf("parse message template %q: %w", name, err)
	}

	found := make(map[string]struct{})
	for _, action := range actionPattern.FindAllString(content, -1) {
		match := placeholderPattern.FindStringSubmatch(action)
		if len(match) != 2 {
			return fmt.Errorf("message template %q contains unsupported action %q", name, action)
		}
		found[match[1]] = struct{}{}
	}

	required := RequiredVariables(name)
	for variable := range found {
		if !slices.Contains(required, variable) {
			return fmt.Errorf("variable %q is not allowed for message template %q", variable, name)
		}
	}
	for _, variable := range required {
		if _, ok := found[variable]; !ok {
			return fmt.Errorf("message template %q requires variable %q", name, variable)
		}
	}
	return nil
}
