package models

import (
	"time"
)

//FilterDefinition is the definition of what measurments to forward in what intervals
type FilterDefinition struct {
	Names       []string      `json:"names"`
	Granularity time.Duration `json:"granularity"`
}

//IsInNames checks if the provided name is allowed according to the filterDefiniton
func (d FilterDefinition) IsInNames(nameToCheck string) bool {
	if len(d.Names) == 0 {
		return true
	}
	for _, name := range d.Names {
		if name == nameToCheck {
			return true
		}
	}
	return false
}

//TimestampFilter filters measurements by their timestamp in a stateful manner
type TimestampFilter struct {
	Granularity     time.Duration
	latestTimestamp int64
}

//Passes the measurement through the filter?
func (f *TimestampFilter) Passes(measurement Measurement) bool {
	if f.latestTimestamp == 0 || f.latestTimestamp+f.Granularity.Nanoseconds() <= measurement.Timestamp() {
		f.latestTimestamp = measurement.Timestamp()
		return true
	}
	return false
}

//FilterCollection is the running state of the filter
type FilterCollection struct {
	Definition             FilterDefinition
	timestampFilterPerName map[string]*TimestampFilter
}

//NewFilterCollection creates a new filterState and initializes the map
func NewFilterCollection(definition FilterDefinition) *FilterCollection {
	return &FilterCollection{
		Definition:             definition,
		timestampFilterPerName: make(map[string]*TimestampFilter),
	}
}

//Passes checks if this measurement passes the filter. If it does, it updates the filter accordingly (passes one time max)
func (c *FilterCollection) Passes(name string, measurement Measurement) bool {
	if !c.Definition.IsInNames(name) {
		return false
	}
	if c.Definition.Granularity == 0 {
		return true
	}

	timestampFilter := c.timestampFilterPerName[name]
	if timestampFilter == nil {
		timestampFilter = &TimestampFilter{Granularity: c.Definition.Granularity}
		c.timestampFilterPerName[name] = timestampFilter
	}

	return timestampFilter.Passes(measurement)
}
