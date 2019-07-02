package poe

import (
	"encoding/json"
	"os"
	"testing"
)

func TestEnvelopeDecode(t *testing.T) {
	f, err := os.Open("example/envelope.json")

	defer f.Close()

	if err != nil {
		t.Logf("could not open test data: %s", err)
		t.Fail()
	}

	var e Envelope

	d := json.NewDecoder(f)
	err = d.Decode(&e)
	if err != nil {
		t.Logf("error parsing envelope: %s", err)
		t.Fail()
	}
}

func TestItemDecode(t *testing.T) {
	types := []string{"ring", "gem", "sword"}
	for _, itemType := range types {
		t.Run(itemType, func(t *testing.T) {
			f, err := os.Open("example/" + itemType + ".json")
			defer f.Close()

			if err != nil {
				t.Logf("could not open test data: %s", err)
				t.Fail()
			}

			var item Item
			d := json.NewDecoder(f)

			err = d.Decode(&item)
			if err != nil {
				t.Logf("%#v", item)
				t.Logf("unmarshal error: %s", err)
				t.Fail()
			}
		})
	}
}

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

}
