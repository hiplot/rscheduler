package task

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"rscheduler/global"
	"rscheduler/rslog"
)

var once sync.Once

func NewHiPlotTask(b []byte) (t *Task) {
	ht := new(HiPlotTask)
	err := json.Unmarshal(b, ht)
	if err != nil {
		global.Logger.Error("Decode HiPlotTask failed, err: " + err.Error())
		return nil
	}
	ht.TryParseURL()
	t = &Task{
		Name:      ht.Name,
		ID:        ht.ID,
		CreatedAt: time.Now(),
		Runner:    ht,
		Logger:    rslog.NewTaskLogger(ht.Name, ht.ID),
		StopLog:   make(chan struct{}),
	}
	return t
}

func (t *HiPlotTask) CommendList() []string {
	commends := make([]string, 2)
	commends[0] = `taskID = "` + t.ID + `"`
	commends[1] = fmt.Sprintf(`hiFunc("%s", "%s", "%s", "%s", "%s")`, t.InputFile, t.ConfFile, t.OutputFilePrefix, t.Tool, t.Module)
	return commends
}

func (t *HiPlotTask) TryParseURL() {
	isUrl := func(str string) bool {
		str = strings.ToLower(str)
		return strings.HasPrefix(str, "http://") || strings.HasPrefix(str, "https://")
	}
	downloadFileFromUrl := func(urlStr string) (localFilepath string) {
		localFilepath = urlStr
		resp, err := http.Get(urlStr)
		if err != nil {
			global.Logger.Error(err)
			return
		}
		defer func() {
			_ = resp.Body.Close()
		}()
		url1, _ := url.Parse(urlStr)
		baseName := filepath.Base(url1.Path)

		currentPath, err := os.Getwd()
		if err != nil {
			global.Logger.Error(err)
			return
		}
		once.Do(func() {
			_ = os.Mkdir("./tmp", 0777)
		})
		tmpLocalFilepath := currentPath + string(os.PathSeparator) + "tmp/" + baseName
		out, err := os.Create(tmpLocalFilepath)
		if err != nil {
			global.Logger.Error(err)
			return
		}
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			global.Logger.Error(err)
			return
		}
		localFilepath = tmpLocalFilepath
		return
	}
	if isUrl(t.InputFile) {
		t.InputFile = downloadFileFromUrl(t.InputFile)
	}
	if isUrl(t.ConfFile) {
		t.ConfFile = downloadFileFromUrl(t.ConfFile)
	}
}
