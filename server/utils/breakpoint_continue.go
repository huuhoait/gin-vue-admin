package utils

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

// FrontendTransmitComeFileSliceAndCurrentSliceForWhatFileofNo.FewSlice
// AfterEndtakeToByAfterCompareTimePartSliceYesNoUpload OrYesNoForNotCompleteAllSlice
// FrontendsendEverySliceMultipleLarge
// FrontendinformYesNoForFinalOneSliceAndYesNoComplete

const (
	breakpointDir = "./breakpointDir/"
	finishDir     = "./fileDir/"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: BreakPointContinue
//@description: Resumable Upload
//@param: content []byte, fileName string, contentNumber int, contentTotal int, fileMd5 string
//@return: error, string

func BreakPointContinue(content []byte, fileName string, contentNumber int, contentTotal int, fileMd5 string) (string, error) {
	if strings.Contains(fileName, "..") || strings.Contains(fileMd5, "..") {
		return "", errors.New("file name or path invalid")
	}
	path := breakpointDir + fileMd5 + "/"
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return path, err
	}
	pathC, err := makeFileContent(content, fileName, path, contentNumber)
	return pathC, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: CheckMd5
//@description: CheckMd5
//@param: content []byte, chunkMd5 string
//@return: CanUpload bool

func CheckMd5(content []byte, chunkMd5 string) (CanUpload bool) {
	fileMd5 := MD5V(content)
	if fileMd5 == chunkMd5 {
		return true // CanByContinueUpload
	} else {
		return false // cutSliceNotCompleteArrange, deprecated
	}
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: makeFileContent
//@description: CreatecutSlicecontent
//@param: content []byte, fileName string, FileDir string, contentNumber int
//@return: string, error

func makeFileContent(content []byte, fileName string, FileDir string, contentNumber int) (string, error) {
	if strings.Contains(fileName, "..") || strings.Contains(FileDir, "..") {
		return "", errors.New("file name or path invalid")
	}
	path := FileDir + fileName + "_" + strconv.Itoa(contentNumber)
	f, err := os.Create(path)
	if err != nil {
		return path, err
	}
	defer f.Close()
	_, err = f.Write(content)
	if err != nil {
		return path, err
	}

	return path, nil
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: makeFileContent
//@description: CreatecutSliceFile
//@param: fileName string, FileMd5 string
//@return: error, string

func MakeFile(fileName string, FileMd5 string) (string, error) {
	if strings.Contains(fileName, "..") || strings.Contains(FileMd5, "..") {
		return "", errors.New("file name or path invalid")
	}
	rd, err := os.ReadDir(breakpointDir + FileMd5)
	if err != nil {
		return finishDir + fileName, err
	}
	_ = os.MkdirAll(finishDir, os.ModePerm)
	fd, err := os.OpenFile(finishDir+fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return finishDir + fileName, err
	}
	defer fd.Close()
	for k := range rd {
		content, _ := os.ReadFile(breakpointDir + FileMd5 + "/" + fileName + "_" + strconv.Itoa(k))
		_, err = fd.Write(content)
		if err != nil {
			_ = os.Remove(finishDir + fileName)
			return finishDir + fileName, err
		}
	}
	return finishDir + fileName, nil
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: RemoveChunk
//@description: RemovecutSlice
//@param: FileMd5 string
//@return: error

func RemoveChunk(FileMd5 string) error {
	if strings.Contains(FileMd5, "..") {
		return errors.New("pathinvalid")
	}
	err := os.RemoveAll(breakpointDir + FileMd5)
	return err
}
