package metrics

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	U "github.com/therealfakemoot/go-unidecode"

	"github.com/therealfakemoot/pom/poe"
)

func SanitizeName(name string) string {
	s := U.Unidecode(strings.ReplaceAll(name, `"`, ""))
	s = strings.ReplaceAll(s, `'`, "")
	s = strings.ReplaceAll(s, `,`, "")
	s = strings.ReplaceAll(s, `-`, "_")
	s = strings.ReplaceAll(s, ` `, "_")
	s = s + "_price_chaos"

	return strings.ToLower(s)
}

type GaugeSet struct {
	GaugeVec *prometheus.GaugeVec
	Gauges   map[poe.GaugeKey]prometheus.Gauge
}

func (gs GaugeSet) Add(i poe.Item) {
	gk := i.Key()
	_, ok := gs.Gauges[gk]
	if !ok {
		gs.Gauges[gk] = gs.GaugeVec.With(i.Labels())
	}
}
