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

type SummarySet struct {
	SummaryVec *prometheus.SummaryVec
	Summaries  map[poe.SummaryKey]prometheus.Observer
}

func (hs SummarySet) Add(i poe.Item) {
	hk := i.Key()
	_, ok := hs.Summaries[hk]
	if !ok {
		hs.Summaries[hk] = hs.SummaryVec.With(i.Labels())
	}
}
