package main

import (
	"encoding/json"
)

// Envelope is the payload envelope containing the pagination ID and
// the actual Stash Tab data.
type Envelope struct {
	NextChangeID string  `json:"next_change_id"`
	Stashes      []Stash `json:"stashes"`
}

// Stash contains account metadata and the actual Item values.
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

// Item describes the attributes of inventory items.
type Item struct {
	AbyssJewel            bool                `json:""`
	AdditionalProperties  map[string]string   `json:""`
	ArtFilename           string              `json:""`
	Category              map[string][]string `json:""`
	Corrupted             bool                `json:""`
	CosmeticMods          []string            `json:""`
	CraftedMods           []string            `json:""`
	DescrText             string              `json:""`
	Duplicated            string              `json:""`
	Elder                 bool                `json:""`
	EnchantMods           []string            `json:""`
	ExplicitMods          []string            `json:""`
	FlavourText           string              `json:""`
	FrameType             FrameType           `json:""`
	Height                int                 `json:"h"`
	Icon                  string              `json:""`
	ID                    string              `json:""`
	Identified            bool                `json:""`
	Ilvl                  int                 `json:""`
	ImplicitMods          []string            `json:""`
	InventoryID           string              `json:""`
	IsRelic               bool                `json:""`
	League                string              `json:""`
	LockedToCharacter     bool                `json:""`
	MaxStackSize          bool                `json:""`
	Name                  string              `json:""`
	NextLevelRequirements Properties          `json:""`
	Note                  string              `json:""`
	Properties            Properties          `json:""`
	ProphecyDiffText      string              `json:""`
	ProphecyText          string              `json:""`
	Requirements          Properties          `json:""`
	SecDescriptionText    string              `json:""`
	Shaper                bool                `json:""`
	// SocketedItems string `json:""`
	Sockets      Sockets  `json:""`
	StackSize    int      `json:""`
	Support      bool     `json:""`
	TalismanTier int      `json:""`
	TypeLine     string   `json:""`
	UtilityMods  []string `json:""`
	Verified     bool     `json:""`
	Width        int      `json:""`
	X            int      `json:""`
	Y            int      `json:""`
}

type Sockets struct{}

type Properties struct {
	Name        string               `json:"name"`
	Values      []RequirementsValues `json:"values"`
	DisplayMode int                  `json:"displayMode"`
	Type        int                  `json:"type,omitempty"`
	Progress    int                  `json:"progress,omitempty"`
}

type ValueType int

const (
	ValueWhitePhys = iota
	ValueBlueModded
	_
	_
	ValueFire
	ValueCold
	ValueLightning
	ValueChaos
)

type RequirementsValues struct {
	Value     string
	ValueType ValueType
}

func (rv *RequirementsValues) Unmarshal(data []byte) error {
	tmp := []interface{}{&rv.Value, &rv.ValueType}
	return json.Unmarshal(data, &tmp)
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
