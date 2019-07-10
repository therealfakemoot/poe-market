package poe

import (
	"encoding/json"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

// GaugeKey is used to identify an item in the context of a metrics.GaugeSet. Had to move it here to prevent an import cycle.
type HistoKey struct {
	Name      string
	Sockets   int
	Links     int
	FrameType string
}

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
	AdditionalProperties  []Properties        `json:"additionalProperties"`
	ArtFilename           string              `json:"artFilename"`
	Category              map[string][]string `json:"category"`
	Corrupted             bool                `json:"corrupted"`
	CosmeticMods          []string            `json:"cosmeticMods"`
	CraftedMods           []string            `json:"craftedMods"`
	DescrText             string              `json:"descrText"`
	Duplicated            bool                `json:"duplicated"`
	Elder                 bool                `json:"elder"`
	EnchantMods           []string            `json:"enchantMods"`
	ExplicitMods          []string            `json:"explicitMods"`
	FlavourText           []string            `json:"flavourText"`
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
	MaxStackSize          int                 `json:"maxStackSize"`
	Name                  string              `json:"name"`
	NextLevelRequirements []Properties        `json:"nextLevelRequirements"`
	Note                  string              `json:"note"`
	Properties            []Properties        `json:"properties"`
	ProphecyDiffText      string              `json:"prophecyDiffText"`
	ProphecyText          string              `json:"prophecyText"`
	Requirements          []Properties        `json:"requirements"`
	SecDescriptionText    string              `json:"secDescrText"`
	Shaper                bool                `json:"shaper"`
	SocketedItems         []Item              `json:"socketedItems"`
	Sockets               []SocketGroup       `json:"sockets"`
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

func (i Item) SocketLinks() int {
	var (
		max int
	)

	groups := make(map[int]int)

	for _, s := range i.Sockets {
		groups[s.Group]++
	}

	for _, c := range groups {
		if max < c {
			max = c
		}
	}

	return max
}

func (i Item) Key() HistoKey {
	var name = i.Name

	if i.Name == "" {
		name = i.TypeLine
	}

	return HistoKey{
		Name:      name,
		Sockets:   len(i.Sockets),
		Links:     i.SocketLinks(),
		FrameType: i.FrameType.String(),
	}
}

func (i Item) Labels() prometheus.Labels {
	var name = i.Name

	if i.Name == "" {
		name = i.TypeLine
	}

	return prometheus.Labels{
		"name":      name,
		"sockets":   fmt.Sprintf("%d", len(i.Sockets)),
		"links":     fmt.Sprintf("%d", i.SocketLinks()),
		"frametype": i.FrameType.String(),
	}
}

type SocketAttr struct {
	Type  string
	Abyss bool
}

func (sa *SocketAttr) UnmarshalJSON(data []byte) error {
	var val interface{}

	err := json.Unmarshal(data, &val)
	if err != nil {
		return err
	}

	switch val.(type) {
	case string:
		sa.Type = val.(string)
	case bool:
		sa.Abyss = val.(bool)
	}
	return nil
}

type SocketColour struct {
	Colour string
	Abyss  bool
}

func (sc *SocketColour) UnmarshalJSON(data []byte) error {
	var val interface{}

	err := json.Unmarshal(data, &val)
	if err != nil {
		return err
	}

	switch val.(type) {
	case string:
		sc.Colour = val.(string)
	case bool:
		sc.Abyss = val.(bool)
	}
	return nil
}

type SocketGroup struct {
	Group   int
	Attr    SocketAttr
	SColour string
}

type Properties struct {
	Name        string               `json:"name"`
	Values      []RequirementsValues `json:"values"`
	DisplayMode int                  `json:"displayMode"`
	Type        int                  `json:"type,omitempty"`
	Progress    float64              `json:"progress,omitempty"`
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

func (ft FrameType) String() string {
	switch ft {
	case 0:
		return "normal"
	case 1:
		return "magic"
	case 2:
		return "rare"
	case 3:
		return "unique"
	case 4:
		return "gem"
	case 5:
		return "currency"
	case 6:
		return "divination"
	case 7:
		return "quest"
	case 8:
		return "prophecy"
	case 9:
		return "relic"
	}
	return "unknown"
}

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
