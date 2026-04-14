package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Cutter Implement io.Writer API
// Used forLogslice, strings.Join([]string{director,layout, formats..., level+".log"}, os.PathSeparator)
type Cutter struct {
	level        string        // Log Level(debug, info, warn, error, dpanic, panic, fatal)
	layout       string        // time format 2006-01-02 15:04:05
	formats      []string      // CustomParameter([]string{Director,"2006-01-02", "business"(ThisParameterOptional), level+".log"}
	director     string        // log folder
	retentionDay int           //log retention days
	file         *os.File      // FileSentencehandle
	mutex        *sync.RWMutex // Read WriteLock
}

type CutterOption func(*Cutter)

// CutterWithLayout time format
func CutterWithLayout(layout string) CutterOption {
	return func(c *Cutter) {
		c.layout = layout
	}
}

// CutterWithFormats FormatParameter
func CutterWithFormats(format ...string) CutterOption {
	return func(c *Cutter) {
		if len(format) > 0 {
			c.formats = format
		}
	}
}

func NewCutter(director string, level string, retentionDay int, options ...CutterOption) *Cutter {
	rotate := &Cutter{
		level:        level,
		director:     director,
		retentionDay: retentionDay,
		mutex:        new(sync.RWMutex),
	}
	for i := 0; i < len(options); i++ {
		options[i](rotate)
	}
	return rotate
}

// Write satisfies the io.Writer interface. It writes to the
// appropriate file handle that is currently being used.
// If we have reached rotation time, the target file gets
// automatically rotated, and also purged if necessary.
func (c *Cutter) Write(bytes []byte) (n int, err error) {
	c.mutex.Lock()
	defer func() {
		if c.file != nil {
			_ = c.file.Close()
			c.file = nil
		}
		c.mutex.Unlock()
	}()
	length := len(c.formats)
	values := make([]string, 0, 3+length)
	values = append(values, c.director)
	if c.layout != "" {
		values = append(values, time.Now().Format(c.layout))
	}
	for i := 0; i < length; i++ {
		values = append(values, c.formats[i])
	}
	values = append(values, c.level+".log")
	filename := filepath.Join(values...)
	director := filepath.Dir(filename)
	err = os.MkdirAll(director, os.ModePerm)
	if err != nil {
		return 0, err
	}
	defer func() {
		err := removeNDaysFolders(c.director, c.retentionDay)
		if err != nil {
			fmt.Println("CleanExpiredLogfailed", err)
		}
	}()

	c.file, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return 0, err
	}
	return c.file.Write(bytes)
}

func (c *Cutter) Sync() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.file != nil {
		return c.file.Sync()
	}
	return nil
}

// IncreaseLogDirectoryFileClean SmallAtetc.AtzeroofValuedefaultIgnoreNotAgainHandle
func removeNDaysFolders(dir string, days int) error {
	if days <= 0 {
		return nil
	}
	cutoff := time.Now().AddDate(0, 0, -days)
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.ModTime().Before(cutoff) && path != dir {
			err = os.RemoveAll(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
