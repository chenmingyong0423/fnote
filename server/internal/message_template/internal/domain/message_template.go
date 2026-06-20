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
	"fmt"
	"strings"
	"text/template"
)

type Name string

type RecipientType uint

const (
	CommentReceived     Name = "comment"
	UserCommentApproved Name = "user-comment-approval"
	UserCommentRejected Name = "user-comment-disapproval"
	UserCommentReplied  Name = "user-comment-reply"
	FriendApplied       Name = "friend"
	UserFriendApproved  Name = "friend-approval"
	UserFriendRejected  Name = "friend-rejection"
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
	Name    Name
	Title   string
	Content string
}

func (mt *MessageTemplate) RenderContent(data Data) error {
	content := normalizeLegacyContent(mt.Name, mt.Content)
	if strings.Contains(content, "%s") {
		return fmt.Errorf("message template %q contains unsupported legacy placeholders", mt.Name)
	}

	tpl, err := template.New(string(mt.Name)).Option("missingkey=error").Parse(content)
	if err != nil {
		return fmt.Errorf("parse message template %q: %w", mt.Name, err)
	}

	var rendered bytes.Buffer
	if err = tpl.Execute(&rendered, data); err != nil {
		return fmt.Errorf("render message template %q: %w", mt.Name, err)
	}
	mt.Content = rendered.String()
	return nil
}

func normalizeLegacyContent(name Name, content string) string {
	var variables []string
	switch name {
	case UserCommentApproved, UserCommentReplied:
		variables = []string{VariablePostURL}
	case UserCommentRejected:
		variables = []string{VariablePostURL, VariableReason}
	case UserFriendApproved:
		variables = []string{VariableFriendPageURL}
	case UserFriendRejected:
		variables = []string{VariableFriendPageURL, VariableReason}
	}
	for _, variable := range variables {
		content = strings.Replace(content, "%s", "{{."+variable+"}}", 1)
	}
	return content
}
