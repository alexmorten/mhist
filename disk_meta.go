package mhist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

var metaFilePath = "meta.json"

//DiskMeta holds the meta info for (Un)Marshalization
type DiskMeta struct {
	//sync maps would be better here, but are not easy to marshalize
	NameToID map[string]int64          `json:"name_to_id"`
	IDToName map[int64]string          `json:"id_to_name"`
	IDToType map[int64]MeasurementType `json:"id_to_type"`

	HighestID int64 `json:"highest_id"`

	sync.RWMutex
}

//InitMetaFromDisk ...
func InitMetaFromDisk() *DiskMeta {
	byteSlice, err := ioutil.ReadFile(filepath.Join(dataPath, metaFilePath))
	if err != nil {
		//assume no file exists
		return NewDiskMeta()
	}
	meta := &DiskMeta{}
	err = json.Unmarshal(byteSlice, meta)
	if err != nil {
		//assume file invalid/corrupted
		meta = NewDiskMeta()
		meta.sync() //Overwrite bad file
		return meta
	}
	return meta
}

//NewDiskMeta with values initialized
func NewDiskMeta() *DiskMeta {
	return &DiskMeta{
		NameToID: map[string]int64{},
		IDToName: map[int64]string{},
		IDToType: map[int64]MeasurementType{},
	}
}

//GetOrCreateID for name, checks if MeasurementType is correct
func (m *DiskMeta) GetOrCreateID(name string, t MeasurementType) (int64, error) {
	m.RLock()
	id := m.NameToID[name]
	m.RUnlock()

	if id != 0 {
		if savedType := m.IDToType[id]; savedType != t {
			return 0, fmt.Errorf("had type %v but was provided %v", savedType, t)
		}
		return id, nil
	}

	m.Lock()
	defer m.Unlock()

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
	m.RLock()
	defer m.RUnlock()
	return m.IDToName[id]
}

//GetTypeForID to translate back form csv to record
func (m *DiskMeta) GetTypeForID(id int64) MeasurementType {
	m.RLock()
	defer m.RUnlock()
	return m.IDToType[id]
}

//GetAllStoredNames from meta
func (m *DiskMeta) GetAllStoredNames() (names []string) {
	m.RLock()
	defer m.RUnlock()
	for name := range m.NameToID {
		names = append(names, name)
	}
	return
}

func (m *DiskMeta) sync() {
	byteSlice, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Errorf("%v ,couldn't marshal diskMeta %v, this shouldn't happen ", err, m))
	}
	err = ioutil.WriteFile(filepath.Join(dataPath, metaFilePath), byteSlice, os.ModePerm)
	if err != nil {
		panic(fmt.Errorf("%v, couldn't write diskMeta %v, this shouldn't happen ", err, m))
	}
}
