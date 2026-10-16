package domain

import "time"

type Error struct {
	Status  int
	Message string
}

func (e *Error) Error() string { return e.Message }

var (
	ErrNotFound = &Error{404, "对话不存在或无权访问"}
	ErrBusy     = &Error{409, "当前已有回答正在生成，请等待或先停止"}
	ErrDisabled = &Error{503, "AI 助手暂未开启"}
)

type Reference struct {
	Title string `bson:"title" json:"title"`
	Href  string `bson:"href" json:"href"`
}

type Message struct {
	ID         string      `bson:"id" json:"id"`
	RequestID  string      `bson:"request_id,omitempty" json:"request_id,omitempty"`
	Role       string      `bson:"role" json:"role"`
	Content    string      `bson:"content" json:"content"`
	Status     string      `bson:"status" json:"status"`
	References []Reference `bson:"references,omitempty" json:"references,omitempty"`
	Error      string      `bson:"error,omitempty" json:"error,omitempty"`
	CreatedAt  time.Time   `bson:"created_at" json:"created_at"`
}

type Generation struct {
	RequestID string `bson:"request_id" json:"request_id"`
	MessageID string `bson:"message_id" json:"message_id"`
	Status    string `bson:"status" json:"status"`
}

type Conversation struct {
	ID         string      `bson:"_id" json:"id"`
	VisitorID  string      `bson:"visitor_id" json:"-"`
	Title      string      `bson:"title" json:"title"`
	Messages   []Message   `bson:"messages" json:"messages"`
	Generation *Generation `bson:"generation,omitempty" json:"generation,omitempty"`
	CreatedAt  time.Time   `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time   `bson:"updated_at" json:"updated_at"`
}

type Visitor struct {
	ID        string    `bson:"_id"`
	FirstIP   string    `bson:"first_ip"`
	LastIP    string    `bson:"last_ip"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

// Events are sent on the POST response; there is no subscription endpoint.
type Event struct {
	Type           string      `json:"type"`
	ConversationID string      `json:"conversation_id,omitempty"`
	RequestID      string      `json:"request_id,omitempty"`
	MessageID      string      `json:"message_id,omitempty"`
	Delta          string      `json:"delta,omitempty"`
	Status         string      `json:"status,omitempty"`
	Message        string      `json:"message,omitempty"`
	Tool           string      `json:"tool,omitempty"`
	References     []Reference `json:"references,omitempty"`
}
