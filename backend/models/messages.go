package models

type MessageSubtype string

// Defines the possible values for MessageSubtype based on Slack's API.
const (
	SubtypeDefault         MessageSubtype = ""
	SubtypeChannelJoin     MessageSubtype = "channel_join"
	SubtypeChannelLeave    MessageSubtype = "channel_leave"
	SubtypeChannelName     MessageSubtype = "channel_name"
	SubtypeMessageChanged  MessageSubtype = "message_changed"
	SubtypeMessageDeleted  MessageSubtype = "message_deleted"
	SubtypeThreadBroadcast MessageSubtype = "thread_broadcast"
	SubtypeFileShare       MessageSubtype = "file_share"
	SubtypeMeMessage       MessageSubtype = "me_message"
	SubtypeBotMessage      MessageSubtype = "bot_message"
)

type EventType string

// Defines the possible values for EventType based on Slack's API.
const (
	TypeMessage       EventType = "message"
	TypeAppMention    EventType = "app_mention"
	TypeReactionAdded EventType = "reaction_added"
	TypeTeamJoin      EventType = "team_join"
	TypeFileShared    EventType = "file_shared"
	// Add other event types as needed...
)

// Message represents a single object in the JSON array.
type Message struct {
	Subtype        MessageSubtype `json:"subtype,omitempty"`
	UserId         string         `json:"user,omitempty"`
	Text           string         `json:"text,omitempty"`
	Inviter        string         `json:"inviter,omitempty"`
	Type           EventType      `json:"type,omitempty"`
	Ts             string         `json:"ts,omitempty"`
	ClientMsgID    string         `json:"client_msg_id,omitempty"`
	Team           string         `json:"team,omitempty"`
	UserProfile    *UserProfile   `json:"user_profile,omitempty"` // Pointer for optional object
	Blocks         []Block        `json:"blocks,omitempty"`
	Reactions      []Reaction     `json:"reactions,omitempty"`
	ThreadTs       string         `json:"thread_ts,omitempty"`
	ParentUserID   string         `json:"parent_user_id,omitempty"`
	ReplyCount     int            `json:"reply_count,omitempty"`
	Replies        []Reply        `json:"replies,omitempty"`
	Edited         *Edited        `json:"edited,omitempty"` // Pointer for optional object
	Files          []File         `json:"files,omitempty"`
	Attachments    []Attachment   `json:"attachments,omitempty"`
	UserName       string
	UserProfileUrl string
}

// UserProfile contains information about the message author.
type UserProfile struct {
	FirstName string `json:"first_name"`
	RealName  string `json:"real_name"`
	Image72   string `json:"image_72"`
	Team      string `json:"team"`
	Name      string `json:"name"`
}

// Block represents a content block in a message.
type Block struct {
	Type     string    `json:"type"`
	BlockID  string    `json:"block_id"`
	Elements []Element `json:"elements"`
}

// Element is a recursive struct to handle nested elements within blocks.
type Element struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	UserID   string    `json:"user_id,omitempty"`
	URL      string    `json:"url,omitempty"`
	Name     string    `json:"name,omitempty"`     // For emoji
	Elements []Element `json:"elements,omitempty"` // For nested structures
}

// Reaction represents an emoji reaction to a message.
type Reaction struct {
	Name  string   `json:"name"`
	Users []string `json:"users"`
	Count int      `json:"count"`
}

// Reply represents a reply in a thread.
type Reply struct {
	User string `json:"user"`
	Ts   string `json:"ts"`
}

// Edited contains information about an edited message.
type Edited struct {
	User string `json:"user"`
	Ts   string `json:"ts"`
}

// File represents a file attached to a message.
type File struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Title    string `json:"title"`
	Mimetype string `json:"mimetype"`
	FileURL  string `json:"url_private"`
	User     string `json:"user"`
}

// Attachment represents a URL attachment with a preview.
type Attachment struct {
	Title     string `json:"title"`
	TitleLink string `json:"title_link"`
	Text      string `json:"text"`
	ImageURL  string `json:"image_url"`
}
