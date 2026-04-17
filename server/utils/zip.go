package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func safeJoinUnderDir(destDir, name string) (string, error) {
	// Reject absolute paths and Windows volume paths early.
	// zip entries always use forward slashes, but callers might pass other forms.
	if name == "" {
		return "", fmt.Errorf("empty path")
	}
	if filepath.IsAbs(name) || filepath.VolumeName(name) != "" {
		return "", fmt.Errorf("%s file nameinvalid", name)
	}

	cleanDest := filepath.Clean(destDir)
	// Convert zip's forward slashes to OS separator and clean.
	cleanName := filepath.Clean(filepath.FromSlash(name))
	// Ensure the entry is treated as a relative path.
	cleanName = strings.TrimPrefix(cleanName, string(filepath.Separator))

	if cleanName == "." {
		return "", fmt.Errorf("%s file nameinvalid", name)
	}
	if strings.Contains(cleanName, "..") {
		return "", fmt.Errorf("%s file nameinvalid", name)
	}

	target := filepath.Join(cleanDest, cleanName)
	absDest, err := filepath.Abs(cleanDest)
	if err != nil {
		return "", err
	}
	absTarget, err := filepath.Abs(target)
	if err != nil {
		return "", err
	}
	sep := string(filepath.Separator)
	if !strings.HasPrefix(absTarget+sep, absDest+sep) {
		return "", fmt.Errorf("%s file nameinvalid", name)
	}
	return target, nil
}

// unzip
func Unzip(zipFile string, destDir string) ([]string, error) {
	zipReader, err := zip.OpenReader(zipFile)
	var paths []string
	if err != nil {
		return []string{}, err
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		fpath, err := safeJoinUnderDir(destDir, f.Name)
		if err != nil {
			return []string{}, err
		}
		paths = append(paths, fpath)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return []string{}, err
			}

			inFile, err := f.Open()
			if err != nil {
				return []string{}, err
			}
			defer inFile.Close()

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return []string{}, err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return []string{}, err
			}
		}
	}
	return paths, nil
}
