package upload

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/utils"
	"go.uber.org/zap"
)

var mu sync.Mutex

type Local struct{}

//@author: [piexlmax](https://github.com/piexlmax)
//@author: [ccfish86](https://github.com/ccfish86)
//@author: [SliverHorn](https://github.com/SliverHorn)
//@object: *Local
//@function: UploadFile
//@description: upload file
//@param: file *multipart.FileHeader
//@return: string, string, error

func (*Local) UploadFile(file *multipart.FileHeader) (string, string, error) {
	// ReadFileAfterSuffix
	ext := filepath.Ext(file.Filename)
	// Readfile nameAndEncrypt
	name := strings.TrimSuffix(file.Filename, ext)
	name = utils.MD5V([]byte(name))
	// ConcatenateNewfile name
	filename := name + "_" + time.Now().Format("20060102150405") + ext
	// TryCreateThispath
	mkdirErr := os.MkdirAll(global.GVA_CONFIG.Local.StorePath, os.ModePerm)
	if mkdirErr != nil {
		global.GVA_LOG.Error("function os.MkdirAll() failed", zap.Any("err", mkdirErr.Error()))
		return "", "", errors.New("function os.MkdirAll() failed, err:" + mkdirErr.Error())
	}
	// ConcatenatepathAndfile name
	p := global.GVA_CONFIG.Local.StorePath + "/" + filename
	filepath := global.GVA_CONFIG.Local.Path + "/" + filename

	f, openError := file.Open() // ReadFile
	if openError != nil {
		global.GVA_LOG.Error("function file.Open() failed", zap.Any("err", openError.Error()))
		return "", "", errors.New("function file.Open() failed, err:" + openError.Error())
	}
	defer f.Close() // create file, defer close

	out, createErr := os.Create(p)
	if createErr != nil {
		global.GVA_LOG.Error("function os.Create() failed", zap.Any("err", createErr.Error()))

		return "", "", errors.New("function os.Create() failed, err:" + createErr.Error())
	}
	defer out.Close() // create file, defer close

	_, copyErr := io.Copy(out, f) // TransmitInput(Copy)File
	if copyErr != nil {
		global.GVA_LOG.Error("function io.Copy() failed", zap.Any("err", copyErr.Error()))
		return "", "", errors.New("function io.Copy() failed, err:" + copyErr.Error())
	}
	return filepath, filename, nil
}

//@author: [piexlmax](https://github.com/piexlmax)
//@author: [ccfish86](https://github.com/ccfish86)
//@author: [SliverHorn](https://github.com/SliverHorn)
//@object: *Local
//@function: DeleteFile
//@description: deleteFile
//@param: key string
//@return: error

func (*Local) DeleteFile(key string) error {
	// Check key YesNoEmpty
	if key == "" {
		return errors.New("keycannot be empty")
	}

	// Verify key YesNoPackageIncludeNon-methodCharacterOrTryAccessSaveStorepathOfoutsideofFile
	if strings.Contains(key, "..") || strings.ContainsAny(key, `\/:*?"<>|`) {
		return errors.New("Non-methodofkey")
	}

	p := filepath.Join(global.GVA_CONFIG.Local.StorePath, key)
	// Defense-in-depth: resolve both paths and verify p stays inside StorePath.
	// Protects against edge cases where the key check above is bypassed.
	absStore, absErr := filepath.Abs(global.GVA_CONFIG.Local.StorePath)
	absP, absErr2 := filepath.Abs(p)
	if absErr != nil || absErr2 != nil || !strings.HasPrefix(absP+string(filepath.Separator), absStore+string(filepath.Separator)) {
		return errors.New("Non-methodofkey")
	}

	// CheckFileYesNoExists
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return errors.New("FileDoes not exist")
	}

	// UseFileLockPreventConcurrencydelete
	mu.Lock()
	defer mu.Unlock()

	err := os.Remove(p)
	if err != nil {
		return errors.New("file deletion failed: " + err.Error())
	}

	return nil
}
