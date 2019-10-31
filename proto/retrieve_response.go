package proto

import (
	"github.com/alexmorten/mhist/models"
)

// RetrieveResponseFromMeasurementMap constructs a retrieve response from an internal measurement map
func RetrieveResponseFromMeasurementMap(m map[string][]models.Measurement) *RetrieveResponse {
	histories := make(map[string]*MeasurementList)

	for name, measurements := range m {
		list := &MeasurementList{Measurements: make([]*Measurement, 0, len(measurements))}

		for _, measurement := range measurements {
			pM := MeasurementFromModel(measurement)
			if pM == nil {
				continue
			}

			list.Measurements = append(list.Measurements, pM)
		}
		histories[name] = list
	}

	return &RetrieveResponse{
		Histories: histories,
	}
}
