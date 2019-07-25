package webupgrade

import (
	"compress/bzip2"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var indexTmpl = template.Must(template.New("index").Parse(`<html>
<head>
<title>web upgrade</title>
<style>
</style>
</head>
<body>
web upgrade
</body>
</html>
`))

var errRollback = errors.New("roll back error")

func init() {
	http.HandleFunc("/internal/index", func(w http.ResponseWriter, r *http.Request) {
		if err := indexTmpl.Execute(w, nil); err != nil {
			log.Print(err)
		}
	})

	http.HandleFunc("/internal/upgrade", func(w http.ResponseWriter, r *http.Request) {
		isSuc := true
		defer func() {
			if isSuc {
				fmt.Fprintln(w, `{"code":0`)
			} else {
				fmt.Fprint(w, `{"code":1}`)
			}
		}()

		md5Str := r.URL.Query().Get("MD5")
		if len(md5Str) == 0 {
			isSuc = false
			return
		}

		file, _, err := r.FormFile("firmware")
		if err != nil {
			isSuc = false
			return
		}
		defer file.Close()
		if err := doUpdate(file, md5Str); err != nil {
			isSuc = false
		}
	})
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
