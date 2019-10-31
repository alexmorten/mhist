package mhist

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

//GetSortedFileList gets the FileInfo list for data files (not the meta file)
func GetSortedFileList() (FileInfoSlice, error) {
	infoList := FileInfoSlice{}
	files, err := ioutil.ReadDir(dataPath)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), "_values") {
			continue
		}

		info, err := timestampsFromFileName(f.Name())
		if err != nil {
			continue
		}
		info.size = f.Size()
		infoList = append(infoList, info)
	}
	sort.Sort(infoList)
	return infoList, nil
}

//GetFilesInTimeRange gets the FileInfo list for data files in the time range
func GetFilesInTimeRange(start, end int64) (FileInfoSlice, error) {
	allFiles, err := GetSortedFileList()
	if err != nil {
		return nil, err
	}
	filesInTimeRange := FileInfoSlice{}
	for _, fileInfo := range allFiles {
		if fileInfo.isInTimeRange(start, end) {
			filesInTimeRange = append(filesInTimeRange, fileInfo)
		}
	}
	return filesInTimeRange, nil
}

//FileInfo descibes file info
type FileInfo struct {
	name     string
	size     int64
	oldestTs int64
	latestTs int64
}

func (i *FileInfo) indexName() string {
	return i.name
}

func (i *FileInfo) valueLogName() string {
	return fmt.Sprintf("%s_values", i.name)
}

//FileInfoSlice ...
type FileInfoSlice []*FileInfo

//WriteBlockToFile ...
func WriteBlockToFile(b Block) (path string, err error) {
	path = pathTo(fileNameFromTs(b.OldestTs(), b.LatestTs()))
	err = ioutil.WriteFile(path, b.UnderlyingByteSlice(), os.ModePerm)
	return
}

//AppendBlockToFile ...
func AppendBlockToFile(info *FileInfo, block Block) (path string, err error) {
	f, err := os.OpenFile(pathTo(info.name), os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return "", err
	}
	_, err = f.Write(block.UnderlyingByteSlice())
	if err != nil {
		return "", err
	}
	path = pathTo(fileNameFromTs(info.oldestTs, block.LatestTs()))
	err = os.Rename(pathTo(info.name), path)
	return
}

func pathTo(filename string) string {
	return filepath.Join(dataPath, filename)
}

func fileNameFromTs(oldestTs, latestTs int64) string {
	return fmt.Sprintf("%v-%v", oldestTs, latestTs)
}

func timestampsFromFileName(name string) (info *FileInfo, err error) {
	var reg = regexp.MustCompile(`(\d+)-(\d+)`)
	matches := reg.FindStringSubmatch(name)
	if len(matches) != 3 {
		return nil, fmt.Errorf("file does not match regexp correctly: %v", matches)
	}

	oldestTs, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return nil, err
	}
	latestTs, err := strconv.ParseInt(matches[2], 10, 64)
	if err != nil {
		return nil, err
	}

	return &FileInfo{name: name, oldestTs: oldestTs, latestTs: latestTs}, nil
}

func (s FileInfoSlice) Len() int      { return len(s) }
func (s FileInfoSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s FileInfoSlice) Less(i, j int) bool {
	return s[i].latestTs < s[j].latestTs
}

//TotalSize of files
func (s FileInfoSlice) TotalSize() (size int64) {
	for _, info := range s {
		size += info.size
	}
	return
}

func (i *FileInfo) isInTimeRange(start, end int64) bool {
	return (i.latestTs > start && !(i.oldestTs > end))
}

func writeGob(filePath string, object interface{}) error {
	file, err := os.Create(filePath)
	if err == nil {
		encoder := gob.NewEncoder(file)
		err = encoder.Encode(object)
	}
	file.Close()
	return err
}

func readGob(filePath string, object interface{}) error {
	file, err := os.Open(filePath)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	return err
}
