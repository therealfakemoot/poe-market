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
	AbyssJewel            bool                `json:"abyssJewel"`
	AdditionalProperties  Properties          `json:"additionalProperties"`
	ArtFilename           string              `json:"artFilename"`
	Category              map[string][]string `json:"category"`
	Corrupted             bool                `json:"corrupted"`
	CosmeticMods          []string            `json:"cosmeticMods"`
	CraftedMods           []string            `json:"craftedMods"`
	DescrText             string              `json:"descrText"`
	Duplicated            string              `json:"duplicated"`
	Elder                 bool                `json:"elder"`
	EnchantMods           []string            `json:"enchantMods"`
	ExplicitMods          []string            `json:"explicitMods"`
	FlavourText           string              `json:"flavourText"`
	FrameType             FrameType           `json:"frameType"`
	Height                int                 `json:"h"`
	Icon                  string              `json:"icon"`
	ID                    string              `json:"id"`
	Identified            bool                `json:"identified"`
	Ilvl                  int                 `json:"ilvl"`
	ImplicitMods          []string            `json:"implicitMods"`
	InventoryID           string              `json:"inventoryId"`
	IsRelic               bool                `json:"isRelic"`
	League                string              `json:"league"`
	LockedToCharacter     bool                `json:"lockedToCharacter"`
	MaxStackSize          bool                `json:"maxStackSize"`
	Name                  string              `json:"name"`
	NextLevelRequirements Properties          `json:"nextLevelRequirements"`
	Note                  string              `json:"note"`
	Properties            Properties          `json:"properties"`
	ProphecyDiffText      string              `json:"prophecyDiffText"`
	ProphecyText          string              `json:"prophecyText"`
	Requirements          []Properties        `json:"requirements"`
	SecDescriptionText    string              `json:"secDescrText"`
	Shaper                bool                `json:"shaper"`
	SocketedItems         string              `json:"socketedItems"`
	Sockets               Sockets             `json:"sockets"`
	StackSize             int                 `json:"stackSize"`
	Support               bool                `json:"support"`
	TalismanTier          int                 `json:"talismanTier"`
	TypeLine              string              `json:"typeLine"`
	UtilityMods           []string            `json:"utilityMods"`
	Verified              bool                `json:"verified"`
	Width                 int                 `json:"w"`
	X                     int                 `json:"x"`
	Y                     int                 `json:"y"`
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

func (rv *RequirementsValues) UnmarshalJSON(data []byte) error {
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
