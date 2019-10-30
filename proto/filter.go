package proto

import (
	"time"

	"github.com/alexmorten/mhist/models"
)

//ToModel converts the proto Filter to the internal representation
func (f *Filter) ToModel() models.FilterDefinition {
	return models.FilterDefinition{
		Names:       f.Names,
		Granularity: time.Duration(f.GranularityNanos),
	}
}
