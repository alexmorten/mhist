package mhist

import (
	"log"
	"os"
)

// DiskWriter handles writing the measurement index and value log
type DiskWriter struct {
	indexWriter                 *os.File
	valueLogWriter              *os.File
	firstWrittenTs              int64
	lastWrittenTs               int64
	bytesWrittenSinceLastCommit int64
	currentPos                  int64

	maxFileSize int64
	maxDiskSize int64
}

// NewDiskWriter returns a fully initialized DiskWriter
func NewDiskWriter(maxFileSize, maxDiskSize int) (*DiskWriter, error) {
	writer := &DiskWriter{
		maxFileSize: int64(maxFileSize),
		maxDiskSize: int64(maxDiskSize),
	}

	fileList, err := GetSortedFileList()
	if err != nil {
		return nil, err
	}

	if len(fileList) == 0 {
		err := writer.createWriters(pathTo("current"))
		if err != nil {
			return nil, err
		}
		return writer, nil
	}

	latestFile := fileList[len(fileList)-1]
	if latestFile.size < writer.maxFileSize {
		err := writer.createWriters(pathTo(latestFile.indexName()))
		if err != nil {
			return nil, err
		}
		return writer, nil
	}

	err = writer.createWriters(pathTo("current"))
	if err != nil {
		return nil, err
	}
	return writer, nil
}

//Commit the buffered writes to actual disk
func (w *DiskWriter) commit() {
	if w.bytesWrittenSinceLastCommit == 0 {
		return
	}
	w.bytesWrittenSinceLastCommit = 0

	info, err := os.Stat(w.indexWriter.Name())
	mustNotBeError(err)
	err = w.indexWriter.Sync()
	mustNotBeError(err)
	err = w.valueLogWriter.Sync()
	mustNotBeError(err)

	if info.Size() < w.maxFileSize {
		return
	}

	currentIndexPath := w.indexWriter.Name()
	w.valueLogWriter.Close()
	currentLogFilePath := w.valueLogWriter.Name()
	w.valueLogWriter.Close()

	newIndexPath := pathTo(fileNameFromTs(w.firstWrittenTs, w.lastWrittenTs))
	os.Rename(currentIndexPath, newIndexPath)
	os.Rename(currentLogFilePath, newIndexPath+"_values")

	err = w.createWriters(pathTo("current"))
	mustNotBeError(err)

	fileList, err := GetSortedFileList()
	mustNotBeError(err)

	if fileList.TotalSize() > w.maxDiskSize {
		oldestFile := fileList[0]
		err = os.Remove(oldestFile.indexName())
		if err != nil {
			log.Println(err)
		}
		err = os.Remove(oldestFile.valueLogName())
		if err != nil {
			log.Println(err)
		}
	}
}

func (w *DiskWriter) handleAdd(m addMessage) {
	measurement := m.measurement
	if len(m.rawValue) > 0 {
		n, err := w.valueLogWriter.Write(m.rawValue)
		mustNotBeError(err)
		measurement.Value = float64(w.currentPos)
		measurement.Size = int64(n)
		w.currentPos += measurement.Size
	}
	if w.lastWrittenTs < measurement.Ts {
		w.lastWrittenTs = measurement.Ts
	}
	if w.firstWrittenTs == 0 {
		w.firstWrittenTs = measurement.Ts
	}

	b := Block{measurement}.UnderlyingByteSlice()
	n, err := w.indexWriter.Write(b)
	mustNotBeError(err)

	w.bytesWrittenSinceLastCommit += int64(n)
	if w.bytesWrittenSinceLastCommit > maxBuffer {
		w.commit()
	}
}

func (w *DiskWriter) createWriters(path string) error {
	valueF, err := os.OpenFile(path+"_values", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	pos, err := valueF.Seek(0, 2)
	if err != nil {
		return err
	}
	w.valueLogWriter = valueF
	w.currentPos = pos
	w.firstWrittenTs = 0
	w.lastWrittenTs = 0
	indexF, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	w.indexWriter = indexF
	return nil
}

//getFilesInTimeRange gets the FileInfo list for data files in the time range
func (w *DiskWriter) getFilesInTimeRange(start, end int64) (FileInfoSlice, error) {
	allFiles, err := GetSortedFileList()
	if err != nil {
		return nil, err
	}

	currentInfo := &FileInfo{
		name:     w.indexWriter.Name(),
		oldestTs: w.firstWrittenTs,
		latestTs: w.lastWrittenTs,
	}

	allFiles = append(allFiles, currentInfo)

	filesInTimeRange := FileInfoSlice{}
	for _, fileInfo := range allFiles {
		if fileInfo.isInTimeRange(start, end) {
			filesInTimeRange = append(filesInTimeRange, fileInfo)
		}
	}
	return filesInTimeRange, nil
}

func mustNotBeError(err error) {
	if err != nil {
		panic(err)
	}
}
