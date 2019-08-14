// package anytool provider simple tool for gateway and other way
// config file fix anytool.yaml
// exec binary file same as
// memlog useful for debug,default adapter memory

package anytool

import (
	"compress/bzip2"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/thinkgos/memlog"
)

func init() {
	memlog.SetLogger(memlog.AdapterMemory)
	memlog.Info("anytool started")

	http.HandleFunc("/internal/tool", Toolhtml)
	http.HandleFunc("/internal/logs", LogsHtml)

	http.HandleFunc("/internal/api/reboot", Reboot)
	http.HandleFunc("/internal/api/config", Config)
	http.HandleFunc("/internal/api/upgrade", Upgrade)
	http.HandleFunc("/internal/api/logs", Logs)
}

var errRollback = errors.New("roll back error")

// Tool get tool html page
func Toolhtml(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		html404(w, r)
		return
	}

	if err := toolTpl.Execute(w, nil); err != nil {
		log.Println("temple execute failed", err)
	}
}

func Reboot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		html404(w, r)
		return
	}

	_ = exec.Command("reboot").Run()
}

func Config(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		html404(w, r)
		return
	}

	md5Str := r.URL.Query().Get("MD5")
	if len(md5Str) == 0 {
		response(w, CodeSysInvalidArguments)
		return
	}

	file, _, err := r.FormFile("config")
	if err != nil {
		response(w, CodeSysInvalidArguments)
		return
	}
	defer file.Close()
	err = saveConfigFile(file, md5Str)
	if err != nil {
		response(w, CodeSysOperationFailed)
	} else {
		response(w, CodeSuccess)
	}
}

func saveConfigFile(file io.ReadSeeker, md string) error {
	var err error

	// 校验文件的正确性
	h := md5.New()
	if _, err = io.Copy(h, file); err != nil {
		return err
	}
	mdStr := hex.EncodeToString(h.Sum(nil))
	if md != mdStr {
		return errors.New("invalid md5 check failed")
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return err
	}
	// 获取执行程序路径
	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	// 配置文件路径
	filePath := filepath.Join(filepath.Dir(execPath), "anytool.yaml")
	fp, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	defer fp.Close()

	_, err = io.Copy(fp, file)
	return err
}

// Upgrade upgrade firmware
func Upgrade(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		html404(w, r)
		return
	}
	md5Str := r.URL.Query().Get("MD5")
	if len(md5Str) == 0 {
		response(w, CodeSysInvalidArguments)
		return
	}

	file, _, err := r.FormFile("firmware")
	if err != nil {
		response(w, CodeSysException)
		return
	}
	defer file.Close()
	if err := doUpdate(file, md5Str); err != nil {
		response(w, CodeSysOperationFailed)
		return
	}
	response(w, CodeSuccess)
}

func doUpdate(file io.ReadSeeker, md string) error {
	var err error

	// 校验文件的正确性
	h := md5.New()
	if _, err = io.Copy(h, file); err != nil {
		return err
	}
	mdStr := hex.EncodeToString(h.Sum(nil))
	if md != mdStr {
		return errors.New("invalid md5 check failed")
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	// 获取执行程序路径
	execPath, err := os.Executable()
	if err != nil {
		return err
	}
	execDir := filepath.Dir(execPath)       // 程序路径
	execFileName := filepath.Base(execPath) // 文件名

	// 新文件放入新名字
	newPath := filepath.Join(execDir, fmt.Sprintf("%s.new", execFileName))
	fp, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer func() {
		_ = fp.Close()
		_ = os.Remove(newPath) // 防止预留在硬盘中
	}()

	if _, err = io.Copy(fp, bzip2.NewReader(file)); err != nil {
		return err
	}

	// 关闭文件,
	_ = fp.Sync()
	_ = fp.Close()

	// 旧文件名
	oldPath := filepath.Join(execDir, fmt.Sprintf("%s.old", execFileName))
	// 将原程序改为旧文件名
	if err = os.Rename(execPath, oldPath); err != nil {
		return err
	}

	// 将新文件改为执行文件名
	if err = os.Rename(newPath, execPath); err != nil {
		// 修改失败，进行回滚
		if rerr := os.Rename(oldPath, execPath); rerr != nil {
			return errRollback
		}
		return err
	}
	// 修改成功,删除旧文件名
	_ = os.Remove(oldPath)
	return nil
}
