package models

type User struct {
	ID                     string  `json:"id"`
	TeamID                 string  `json:"team_id"`
	Name                   string  `json:"name"`
	Deleted                bool    `json:"deleted"`
	Color                  string  `json:"color,omitempty"`
	RealName               string  `json:"real_name,omitempty"`
	TZ                     string  `json:"tz,omitempty"`
	TZLabel                string  `json:"tz_label,omitempty"`
	TZOffset               int     `json:"tz_offset,omitempty"`
	Profile                Profile `json:"profile"`
	IsAdmin                bool    `json:"is_admin,omitempty"`
	IsOwner                bool    `json:"is_owner,omitempty"`
	IsPrimaryOwner         bool    `json:"is_primary_owner,omitempty"`
	IsRestricted           bool    `json:"is_restricted,omitempty"`
	IsUltraRestricted      bool    `json:"is_ultra_restricted,omitempty"`
	IsBot                  bool    `json:"is_bot"`
	IsAppUser              bool    `json:"is_app_user"`
	Updated                int64   `json:"updated"`
	IsEmailConfirmed       *bool   `json:"is_email_confirmed,omitempty"`
	WhoCanShareContactCard string  `json:"who_can_share_contact_card,omitempty"`
	IsInvitedUser          bool    `json:"is_invited_user,omitempty"`
}

type Profile struct {
	Title                   string                `json:"title"`
	Phone                   string                `json:"phone"`
	Skype                   string                `json:"skype"`
	RealName                string                `json:"real_name"`
	RealNameNormalized      string                `json:"real_name_normalized"`
	DisplayName             string                `json:"display_name"`
	DisplayNameNormalized   string                `json:"display_name_normalized"`
	Fields                  map[string]FieldValue `json:"fields"`
	StatusText              string                `json:"status_text"`
	StatusEmoji             string                `json:"status_emoji"`
	StatusEmojiDisplayInfo  []interface{}         `json:"status_emoji_display_info"`
	StatusExpiration        int                   `json:"status_expiration"`
	AvatarHash              string                `json:"avatar_hash"`
	ImageOriginal           string                `json:"image_original,omitempty"`
	IsCustomImage           *bool                 `json:"is_custom_image,omitempty"`
	Email                   string                `json:"email,omitempty"`
	HuddleState             string                `json:"huddle_state,omitempty"`
	HuddleStateExpirationTs int                   `json:"huddle_state_expiration_ts,omitempty"`
	FirstName               string                `json:"first_name,omitempty"`
	LastName                string                `json:"last_name,omitempty"`
	Image24                 string                `json:"image_24"`
	Image32                 string                `json:"image_32"`
	Image48                 string                `json:"image_48"`
	Image72                 string                `json:"image_72"`
	Image192                string                `json:"image_192"`
	Image512                string                `json:"image_512"`
	Image1024               string                `json:"image_1024,omitempty"`
	StatusTextCanonical     string                `json:"status_text_canonical"`
	Team                    string                `json:"team"`
	APIAppID                string                `json:"api_app_id,omitempty"`
	BotID                   string                `json:"bot_id,omitempty"`
	AlwaysActive            bool                  `json:"always_active,omitempty"`
}

// FieldValue represents the value and alt text for a custom profile field.
type FieldValue struct {
	Value string `json:"value"`
	Alt   string `json:"alt"`
}

type ChannelTopic struct {
	Value   string `json:"value"`
	Creator string `json:"creator"`
	LastSet int64  `json:"last_set"`
}

type ChannelPurpose struct {
	Value   string `json:"value"`
	Creator string `json:"creator"`
	LastSet int64  `json:"last_set"`
}

type ChannelPin struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Created int64  `json:"created"`
	User    string `json:"user"`
	Owner   string `json:"owner"`
}

type Channel struct {
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	Created    int64          `json:"created"` // Timestamp of channel creation
	IsChannel  bool           `json:"is_channel"`
	Creator    string         `json:"creator"` // User ID of the channel creator
	IsArchived bool           `json:"is_archived"`
	IsGeneral  bool           `json:"is_general"`
	Members    []string       `json:"members"` // List of user IDs
	Topic      ChannelTopic   `json:"topic"`
	Purpose    ChannelPurpose `json:"purpose"`
	Pins       []ChannelPin   `json:"pins,omitempty"`
}
