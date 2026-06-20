package domain

import (
	"strings"
	"testing"
)

func TestMessageTemplateRenderContent(t *testing.T) {
	tests := []struct {
		name    string
		tpl     MessageTemplate
		data    Data
		want    string
		wantErr string
	}{
		{
			name: "named variables",
			tpl: MessageTemplate{
				Name:    UserCommentRejected,
				Content: "post={{.PostURL}}, reason={{.Reason}}",
			},
			data: Data{VariablePostURL: "https://example.com/posts/1", VariableReason: "spam"},
			want: "post=https://example.com/posts/1, reason=spam",
		},
		{
			name: "legacy positional variables",
			tpl: MessageTemplate{
				Name:    UserFriendRejected,
				Content: "page=%s, reason=%s",
			},
			data: Data{VariableFriendPageURL: "https://example.com/friend", VariableReason: "invalid URL"},
			want: "page=https://example.com/friend, reason=invalid URL",
		},
		{
			name: "missing variable",
			tpl: MessageTemplate{
				Name:    UserCommentApproved,
				Content: "post={{.PostURL}}",
			},
			wantErr: "map has no entry for key \"PostURL\"",
		},
		{
			name: "invalid syntax",
			tpl: MessageTemplate{
				Name:    UserCommentApproved,
				Content: "post={{.PostURL",
			},
			wantErr: "unclosed action",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.tpl.RenderContent(tt.data)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("RenderContent() error = %v, want error containing %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("RenderContent() error = %v", err)
			}
			if tt.tpl.Content != tt.want {
				t.Fatalf("RenderContent() content = %q, want %q", tt.tpl.Content, tt.want)
			}
		})
	}
}

func TestMessageTemplateRenderLegacyContent(t *testing.T) {
	tests := []struct {
		name    Name
		content string
		data    Data
		want    string
	}{
		{
			name:    UserCommentApproved,
			content: "post=%s",
			data:    Data{VariablePostURL: "post-url"},
			want:    "post=post-url",
		},
		{
			name:    UserCommentRejected,
			content: "post=%s, reason=%s",
			data:    Data{VariablePostURL: "post-url", VariableReason: "reason"},
			want:    "post=post-url, reason=reason",
		},
		{
			name:    UserCommentReplied,
			content: "post=%s",
			data:    Data{VariablePostURL: "post-url"},
			want:    "post=post-url",
		},
		{
			name:    UserFriendApproved,
			content: "page=%s",
			data:    Data{VariableFriendPageURL: "friend-page-url"},
			want:    "page=friend-page-url",
		},
		{
			name:    UserFriendRejected,
			content: "page=%s, reason=%s",
			data:    Data{VariableFriendPageURL: "friend-page-url", VariableReason: "reason"},
			want:    "page=friend-page-url, reason=reason",
		},
	}

	for _, tt := range tests {
		t.Run(string(tt.name), func(t *testing.T) {
			tpl := MessageTemplate{Name: tt.name, Content: tt.content}
			if err := tpl.RenderContent(tt.data); err != nil {
				t.Fatalf("RenderContent() error = %v", err)
			}
			if tpl.Content != tt.want {
				t.Fatalf("RenderContent() content = %q, want %q", tpl.Content, tt.want)
			}
		})
	}
}

func TestValidateContent(t *testing.T) {
	tests := []struct {
		name      string
		tplName   Name
		content   string
		wantError bool
	}{
		{name: "valid variables", tplName: UserCommentRejected, content: "post={{.PostURL}}, reason={{.Reason}}"},
		{name: "no variables required", tplName: CommentReceived, content: "new comment"},
		{name: "missing required variable", tplName: UserCommentRejected, content: "post={{.PostURL}}", wantError: true},
		{name: "unknown variable", tplName: UserCommentApproved, content: "post={{.Unknown}}", wantError: true},
		{name: "legacy placeholder", tplName: UserCommentApproved, content: "post=%s", wantError: true},
		{name: "unsupported action", tplName: UserCommentApproved, content: "{{if .PostURL}}{{.PostURL}}{{end}}", wantError: true},
		{name: "invalid syntax", tplName: UserCommentApproved, content: "{{.PostURL", wantError: true},
		{name: "empty content", tplName: CommentReceived, content: "  ", wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateContent(tt.tplName, tt.content)
			if (err != nil) != tt.wantError {
				t.Fatalf("ValidateContent() error = %v, wantError = %v", err, tt.wantError)
			}
		})
	}
}
