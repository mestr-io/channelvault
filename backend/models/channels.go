package models

type ChannelTopic struct {
	Value   string `json:"value"`
	Creator string `json:"creator"` // User ID of the creator of the topic
	LastSet string `json:"last_set"`
}

type ChannelPurpose struct {
	Value   string `json:"value"`
	Creator string `json:"creator"` // User ID of the creator of the purpose
	LastSet string `json:"last_set"`
}

type ChannelPin struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Created string `json:"created"`
	User    string `json:"user"`  // User ID of the user who pinned the item
	Owner   string `json:"owner"` // User ID of the owner of the pinned item
}

type Channel struct {
	ID         string         `json:"id"`
	Name       string         `json:"name"`
	Created    string         `json:"created"` // Timestamp of channel creation
	Creator    string         `json:"creator"` // User ID of the channel creator
	IsArchived bool           `json:"is_archived"`
	IsGeneral  bool           `json:"is_general"`
	Members    []string       `json:"members"` // List of user IDs, won't be used in this version
	Topic      ChannelTopic   `json:"topic"`
	Purpose    ChannelPurpose `json:"purpose"`
	Pins       []ChannelPin   `json:"pins,omitempty"`
}
