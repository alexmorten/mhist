package mhist

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/alexmorten/mhist/models"
)

var metaFilePath = "meta.gob"

//DiskMeta holds the meta info for (Un)Marshalization
type DiskMeta struct {
	//sync maps would be better here, but are not easy to marshalize
	NameToID           map[string]int64
	IDToName           map[int64]string
	IDToType           map[int64]models.MeasurementType
	CategoricalMapping *CategoricalMapping
	HighestID          int64

	mutex sync.RWMutex
}

//MeasurementTypeInfo ...
type MeasurementTypeInfo struct {
	Name string                 `json:"name"`
	Type models.MeasurementType `json:"type"`
}

//InitMetaFromDisk ...
func InitMetaFromDisk() *DiskMeta {
	meta := &DiskMeta{}
	err := readGob(filepath.Join(dataPath, metaFilePath), meta)
	if err != nil {
		//assume no file exists
		return NewDiskMeta()
	}
	return meta
}

//NewDiskMeta with values initialized
func NewDiskMeta() *DiskMeta {
	return &DiskMeta{
		NameToID:           map[string]int64{},
		IDToName:           map[int64]string{},
		IDToType:           map[int64]models.MeasurementType{},
		CategoricalMapping: NewCategoricalMapping(),
	}
}

//GetOrCreateID for name, checks if MeasurementType is correct
func (m *DiskMeta) GetOrCreateID(name string, t models.MeasurementType) (int64, error) {
	m.mutex.RLock()
	id := m.NameToID[name]
	m.mutex.RUnlock()

	if id != 0 {
		if savedType := m.IDToType[id]; savedType != t {
			return 0, fmt.Errorf("had type %v but was provided %v", savedType, t)
		}
		return id, nil
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	//Check again to be sure we don't have multiple in parralel
	id = m.NameToID[name]
	if id != 0 {
		return id, nil
	}
	m.HighestID++
	m.NameToID[name] = m.HighestID
	m.IDToName[m.HighestID] = name
	m.IDToType[m.HighestID] = t
	m.sync()
	return m.HighestID, nil
}

//GetNameForID to translate back form csv to record
func (m *DiskMeta) GetNameForID(id int64) string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.IDToName[id]
}

//GetTypeForID to translate back form csv to record
func (m *DiskMeta) GetTypeForID(id int64) models.MeasurementType {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.IDToType[id]
}

//GetAllStoredInfos from meta
func (m *DiskMeta) GetAllStoredInfos() (infos []MeasurementTypeInfo) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	for name, id := range m.NameToID {
		info := MeasurementTypeInfo{
			Name: name,
			Type: m.IDToType[id],
		}
		infos = append(infos, info)
	}
	return
}

//GetValueIDForCategoricalValue ... also syncs the current meta to disk, if necessary
func (m *DiskMeta) GetValueIDForCategoricalValue(id int64, categoricalValue string) float64 {
	valueID, created := m.CategoricalMapping.GetOrCreateValueIDMap(id).GetOrCreateValueIDForValue(categoricalValue)
	if created {
		m.sync()
	}
	return valueID
}

func (m *DiskMeta) sync() {
	err := writeGob(filepath.Join(dataPath, metaFilePath), m)
	if err != nil {
		panic(fmt.Errorf("%v, couldn't write diskMeta %v, this shouldn't happen ", err, m))
	}
}

//ValueIDMapping is a bi-directional mapping between a categorical value and it's value ID
type ValueIDMapping struct {
	ValueToValueID map[string]float64
	ValueIDToValue map[float64]string

	HighestValueID float64

	mutex sync.RWMutex
}

//NewValueIDMapping with initialized maps
func NewValueIDMapping() *ValueIDMapping {
	return &ValueIDMapping{
		ValueToValueID: map[string]float64{},
		ValueIDToValue: map[float64]string{},
		HighestValueID: 1,
	}
}

//GetOrCreateValueIDForValue thread-safely
func (m *ValueIDMapping) GetOrCreateValueIDForValue(value string) (valueID float64, created bool) {
	m.mutex.RLock()
	valueID = m.ValueToValueID[value]
	m.mutex.RUnlock()

	if valueID != 0 {
		return valueID, false
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.HighestValueID++
	m.ValueToValueID[value] = m.HighestValueID
	m.ValueIDToValue[m.HighestValueID] = value

	return m.HighestValueID, true
}

//CategoricalMapping enables translating categorical values into a float that is easily dumpable to disk
type CategoricalMapping struct {
	IDToValueIDMap map[int64]*ValueIDMapping

	mutex sync.RWMutex
}

//NewCategoricalMapping with initialized map
func NewCategoricalMapping() *CategoricalMapping {
	return &CategoricalMapping{
		IDToValueIDMap: map[int64]*ValueIDMapping{},
	}
}

//GetOrCreateValueIDMap thread-safely
func (m *CategoricalMapping) GetOrCreateValueIDMap(id int64) *ValueIDMapping {
	m.mutex.RLock()
	valueIDMap := m.IDToValueIDMap[id]
	m.mutex.RUnlock()

	if valueIDMap != nil {
		return valueIDMap
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	newValueIDMap := NewValueIDMapping()
	m.IDToValueIDMap[id] = newValueIDMap

	return newValueIDMap
}
