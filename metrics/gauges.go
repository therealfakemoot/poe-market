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

func NewGauge(i poe.Item) prometheus.Gauge {
	sanitized := SanitizeName(i.TypeLine)
	return prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "market",
		Name:      sanitized,
	})
}

type GaugeSet struct {
	Gauges map[string]prometheus.Gauge
}

func (gs GaugeSet) RegisterItem(i poe.Item) {

	if i.Note != "" {
		_, ok := gs.Gauges[i.TypeLine]
		if !ok {
			gs.Gauges[i.TypeLine] = NewGauge(i)
			prometheus.MustRegister(gs.Gauges[i.TypeLine])
		}
	}
}
