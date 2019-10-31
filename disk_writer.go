package mhist

import (
	"log"
	"os"
)

// DiskWriter handles writing the measurement index and value log
type DiskWriter struct {
	block          Block
	valueLogWriter *os.File
	currentPos     int64

	maxFileSize int64
	maxDiskSize int64
}

// NewDiskWriter returns a fully initialized DiskWriter
func NewDiskWriter(maxFileSize, maxDiskSize int) (*DiskWriter, error) {
	writer := &DiskWriter{
		block:       Block{},
		maxFileSize: int64(maxFileSize),
		maxDiskSize: int64(maxDiskSize),
	}

	fileList, err := GetSortedFileList()
	if err != nil {
		return nil, err
	}

	if len(fileList) == 0 {
		logWriter, err := os.Create(pathTo("current_value_log"))
		if err != nil {
			return nil, err
		}
		writer.valueLogWriter = logWriter
		return writer, nil
	}

	latestFile := fileList[len(fileList)-1]
	if latestFile.size < writer.maxFileSize {
		err := writer.setValueWriter(pathTo(latestFile.indexName()))
		if err != nil {
			return nil, err
		}
		return writer, nil
	}

	logWriter, err := os.Create(pathTo("current_value_log"))
	if err != nil {
		return nil, err
	}

	writer.valueLogWriter = logWriter
	return writer, nil
}

//Commit the buffered writes to actual disk
func (w *DiskWriter) commit() {
	if w.block.Size() == 0 {
		return
	}
	err := w.valueLogWriter.Sync()
	if err != nil {
		panic(err)
	}

	fileList, err := GetSortedFileList()
	if err != nil {
		log.Printf("couldn't get file List: %v", err)
		return
	}
	defer func() { w.block = w.block[:0] }()
	if len(fileList) == 0 {
		path, err := WriteBlockToFile(w.block)
		if err != nil {
			panic(err)
		}
		currentLogFilePath := w.valueLogWriter.Name()
		w.valueLogWriter.Close()

		os.Rename(currentLogFilePath, path)

		err = w.setValueWriter(path)
		if err != nil {
			panic(err)
		}

		return
	}
	latestFile := fileList[len(fileList)-1]
	if latestFile.size < w.maxFileSize {
		_, err := AppendBlockToFile(latestFile, w.block)
		if err != nil {
			panic(err)
		}
		return
	}
	path, err := WriteBlockToFile(w.block)
	if err != nil {
		panic(err)
	}
	currentLogFilePath := w.valueLogWriter.Name()
	w.valueLogWriter.Close()

	os.Rename(currentLogFilePath, path)

	err = w.setValueWriter(path)
	if err != nil {
		panic(err)
	}

	if fileList.TotalSize() > w.maxDiskSize {
		oldestFile := fileList[0]
		os.Remove(pathTo(oldestFile.indexName()))
		os.Remove(pathTo(oldestFile.valueLogName()))
	}
}

func (w *DiskWriter) handleAdd(m addMessage) {
	measuremnent := m.measurement
	if len(m.rawValue) > 0 {
		n, err := w.valueLogWriter.Write(m.rawValue)
		if err != nil {
			panic(err)
		}
		measuremnent.Value = float64(w.currentPos)
		measuremnent.Size = int64(n)
		w.currentPos += measuremnent.Size
	}

	w.block = append(w.block, m.measurement)

	if w.block.Size() > maxBuffer {
		w.commit()
	}
}

func (w *DiskWriter) setValueWriter(path string) error {
	f, err := os.OpenFile(path+"_values", os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	pos, err := f.Seek(0, 2)
	if err != nil {
		return err
	}
	w.valueLogWriter = f
	w.currentPos = pos
	return nil
}
