package main

type Envelope struct {
	NextChangeID string  `json:"next_change_id"`
	Stashes      []Stash `json:"stashes"`
}

type Stash struct {
	ID                string      `json:"id"`
	Public            bool        `json:"public"`
	AccountName       interface{} `json:"accountName"`
	LastCharacterName interface{} `json:"lastCharacterName"`
	Stash             interface{} `json:"stash"`
	StashType         string      `json:"stashType"`
	League            string      `json:"league"`
	Items             []Item      `json:"items"`
}

type Item struct {
	AbyssJewel            bool
	AdditionalProperties  map[string]string
	ArtFilename           string
	Category              map[string][]string
	Corrupted             bool
	CosmeticMods          []string
	CraftedMods           []string
	DescrText             string
	Duplicated            string
	Elder                 bool
	EnchantMods           []string
	ExplicitMods          []string
	FlavourText           string
	FrameType             FrameType
	H                     int
	Icon                  string
	ID                    string
	Identified            bool
	Ilvl                  int
	ImplicitMods          []string
	InventoryID           string
	IsRelic               bool
	League                string
	LockedToCharacter     bool
	MaxStackSize          bool
	Name                  string
	NextLevelRequirements Requirements
}

type Requirements Properties

type Properties struct {
	Name        string
	Values      []string
	DisplayMode int
	Type        int
	Progress    int
}

type FrameType int

const (
	FrameTypeNormal = iota
	FrameTypeMagic
	FrameTypeRare
	FrameTypeUnique
	FrameTypeGem
	FrameTypeCurrency
	FrameTypeDiv
	FrameTypeQuest
	FrameTypeProphecy
	FrameTypeRelic
)
