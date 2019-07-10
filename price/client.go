package price

type PricePoint struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Category  string    `json:"category"`
	Group     string    `json:"group"`
	Frame     int       `json:"frame"`
	StackSize int       `json:"stackSize,omitempty"`
	Icon      string    `json:"icon"`
	Mean      float64   `json:"mean"`
	Median    float64   `json:"median"`
	Mode      float64   `json:"mode"`
	Min       float64   `json:"min"`
	Max       float64   `json:"max"`
	Exalted   float64   `json:"exalted"`
	Total     int       `json:"total"`
	Daily     int       `json:"daily"`
	Current   int       `json:"current"`
	Accepted  int       `json:"accepted"`
	Change    float64   `json:"change"`
	History   []float64 `json:"history"`
	Type      string    `json:"type,omitempty"`
}

var IDMap = map[string]int{
	"orb-of-horizons":    418,
	"scr":                1340,
	"whe":                1339,
	"journeyman-sextant": 535,
	"orb-of-annulment":   1343,
	"apprentice-sextant": 114,
	"master-sextant":     1007,
	"mirror":             3282,
	"scour":              1340,
	"silver":             721,
	"regal":              222,
	"blessed":            801,
	"regret":             433,
	"chance":             421,
	"mir":                5498,
	"divine":             422,
	"chisel":             223,
	"jew":                720,
	"chrom":              221,
	"gcp":                113,
	"vaal":               224,
	"fuse":               301,
	"alt":                220,
	"alch":               225,
	"exa":                142,
	"chaos":              0, // omitted because chaos is the base currency
}
