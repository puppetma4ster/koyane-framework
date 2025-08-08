package analyzer

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/puppetma4ster/koyane-framework/internal/core/utils"
	"github.com/saintfish/chardet"
)

type GeneralAnalyzer struct {
	FileName    string
	FilePath    string
	FileSize    uint64
	Extension   string
	Encoding    string
	HashVal     string
	LastChanges string
}

func NewGeneralAnalyzer(path string) (*GeneralAnalyzer, error) {
	filePath, err := utils.ResolvePath(path)
	if err != nil {

	}
	filename := filepath.Base(filePath)
	size, err := FileSize(filePath)
	if err != nil {

	}

	var extension string = filepath.Ext(filePath)
	encodeing, err := detectFileEncoding(filePath)
	if err != nil {

	}

	hashString, err := FileMD5(filePath)
	if err != nil {

	}
	modified, err := os.Stat(filePath)
	if err != nil {
		log.Fatal(err)
	}

	return &GeneralAnalyzer{
		FileName:    filename,
		FilePath:    filePath,
		FileSize:    size,
		Extension:   extension,
		Encoding:    encodeing,
		HashVal:     hashString,
		LastChanges: modified.ModTime().String(),
	}, nil
}

func NewGeneralDummy() *GeneralAnalyzer {
	return &GeneralAnalyzer{
		FileName:    "",
		FilePath:    "",
		FileSize:    0,
		Extension:   "",
		Encoding:    "",
		HashVal:     "",
		LastChanges: "",
	}
}

func detectFileEncoding(path string) (string, error) {
	const sampleSize = 128 * 1024 // 64 KB

	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buf := make([]byte, sampleSize)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	}
	buf = buf[:n]

	detector := chardet.NewTextDetector()
	result, err := detector.DetectBest(buf)
	if err != nil {
		return "", err
	}

	return result.Charset, nil
}

func FileSize(path string) (uint64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return uint64(info.Size()), nil
}

func FileMD5(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
