package controller

import (
	"github.com/cavaliercoder/go-rpm"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("example")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

// User
type RPMInfo struct {
	Repo    string `json:"repo" xml:"repo"`
	Name    string `json:"name" xml:"name"`
	Size    uint64 `json:"size" xml:"size"`
	Version string `json:"version" xml:"version"`
}

func Upload(c echo.Context) error {
	// Read form fields
	repo := c.FormValue("repo")

	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendLeveled := logging.AddModuleLevel(backend)
	backendLeveled.SetLevel(logging.INFO, format)
	logging.SetBackend(backendLeveled)
	//-----------
	// Read file
	//-----------

	// Source
	file, err := c.FormFile("data")
	if err != nil {
		log.Error(err)
		return err
	}
	src, err := file.Open()
	if err != nil {
		log.Error(err)
		return err
	}

	defer src.Close()

	// Crate directory
	path := filepath.Join(viper.GetString("UploadRpmPath"), repo)
	os.MkdirAll(path, os.ModeDir|os.ModePerm)

	// Destination
	dst, err := os.Create(path + string(os.PathSeparator) + file.Filename)
	if err != nil {
		log.Error(err)
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		log.Error(err)
		return err
	}

	p, err := rpm.OpenPackageFile(path + string(os.PathSeparator) + file.Filename)
	
	if err != nil {
		log.Error(err)
		return err
	}

	rpmi := &RPMInfo{
		Repo:    repo,
		Name:    p.Name(),
		Size:    p.Size(),
		Version: p.RPMVersion(),
	}
	return c.JSON(http.StatusOK, rpmi)
}
