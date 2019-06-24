package main

import (
	"encoding/json"
	"testing"
)

/*
func TestItemDecode(t *testing.T) {
	t.Run("simple item", func(t *testing.T) {
		in := `{"verified": false, "w": 1, "h": 1, "ilvl": 45, "icon": "http://web.poecdn.com/image/Art/2DItems/Rings/Ring4.png?scale=1&w=1&h=1&v=6ef89e2afcf83dc3cc86b1522597ed2b", "league": "Standard", "id": "0ef6377fbeb361d2397d260047fbf6e1d4fb9f4c85cb56a3d9ed24bc5a311cc3", "name": "Kraken Grasp", "typeLine": "Gold Ring", "identified": true, "requirements": [{"name": "Level", "values": [["33", 0]], "displayMode": 0}], "implicitMods": ["6% increased Rarity of Items found"], "explicitMods": ["+11 to Strength", "+17 to maximum Energy Shield", "+33 to maximum Life", "46% increased Mana Regeneration Rate", "13% increased Rarity of Items found"], "frameType": 2, "category": {"accessories": ["ring"]}, "x": 6, "y": 3, "inventoryId": "Stash1"}`

		var item Item

		err := json.Unmarshal([]byte(in), &item)
		if err != nil {
			t.Logf("unmarshal error: %s", err)
			t.Fail()
		}
	})

	t.Run("complex item", func(t *testing.T) {
		t.Fail()
	})
}
*/

func TestRequirementsDecode(t *testing.T) {
	t.Run("simple properties", func(t *testing.T) {
		in := `[{"name":"Level","values":[["33",0]],"displayMode":0}]`
		var ps []Properties
		err := json.Unmarshal([]byte(in), &ps)
		if err != nil {
			t.Logf("could not decode properties: %s", err)
			t.Fail()
		}
	})

	/*
		t.Run("complex properties", func(t *testing.T) {
			t.Fail()
		})
	*/
}
