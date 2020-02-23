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
)

// API URL
const (
	URLAPIReboot  = "/api/anytool/reboot"
	URLAPIConfig  = "/api/anytool/config"
	URLAPIUpgrade = "/api/anytool/upgrade"
)

var errRollback = errors.New("roll back error")

// ToolHTML get tool html page
func ToolHTML(w http.ResponseWriter, _ *http.Request) {
	if err := toolTpl.Execute(w, nil); err != nil {
		log.Println("temple execute failed", err)
	}
}

// Reboot 重启命令
func Reboot(_ http.ResponseWriter, _ *http.Request) {
	_ = exec.Command("reboot").Run()
}

// UploadConfigFile 配置命令 method post
func UploadConfigFile(w http.ResponseWriter, r *http.Request) {
	md5Str := r.FormValue("md5")
	if md5Str == "" {
		response(w, http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("config")
	if err != nil {
		response(w, CodeSysInvalidArguments)
		return
	}
	defer file.Close()
	err = doConfigFile(file, md5Str)
	if err != nil {
		response(w, CodeSysOperationFailed)
	} else {
		response(w, http.StatusOK)
	}
}

func doConfigFile(file io.ReadSeeker, md string) error {
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

// Upgrade upgrade firmware ( method post )
func Upgrade(w http.ResponseWriter, r *http.Request) {
	md5Str := r.FormValue("md5")
	if md5Str == "" {
		response(w, CodeSysInvalidArguments)
		return
	}

	file, _, err := r.FormFile("firmware")
	if err != nil {
		response(w, http.StatusInternalServerError)
		return
	}
	defer file.Close()
	if err := doUpdate(file, md5Str); err != nil {
		response(w, CodeSysOperationFailed)
		return
	}
	response(w, http.StatusOK)
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
