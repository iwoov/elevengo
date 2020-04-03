package elevengo

import (
	"errors"
	"github.com/deadblue/elevengo/core"
	"github.com/deadblue/elevengo/internal"
	"time"
)

const (
	apiOfflineSpace   = "https://115.com/"
	apiOfflineList    = "https://115.com/web/lixian/?ct=lixian&ac=task_lists"
	apiOfflineAddUrl  = "https://115.com/web/lixian/?ct=lixian&ac=add_task_url"
	apiOfflineAddUrls = "https://115.com/web/lixian/?ct=lixian&ac=add_task_urls"
	apiOfflineDelete  = "https://115.com/web/lixian/?ct=lixian&ac=task_del"
	apiOfflineClear   = "https://115.com/web/lixian/?ct=lixian&ac=task_clear"

	errOfflineCaptcha      = 911
	errOfflineTaskExisting = 10008
)

// Parameter for "Agent.OfflineClear()" method.
// Default value is to clear all done tasks without deleteing downloaded files.
type OfflineClearParam struct {
	flag int
}

// Clear all tasks, delete the downloaded files if "delete" is true.
func (p *OfflineClearParam) All(delete bool) *OfflineClearParam {
	if delete {
		p.flag = 5
	} else {
		p.flag = 1
	}
	return p
}

// Clear done tasks, delete the downloaded files if "delete" is true.
func (p *OfflineClearParam) Done(delete bool) *OfflineClearParam {
	if delete {
		p.flag = 4
	} else {
		p.flag = 0
	}
	return p
}

// Clear failed tasks.
func (p *OfflineClearParam) Failed() *OfflineClearParam {
	p.flag = 2
	return p
}

// Clean running tasks.
func (p *OfflineClearParam) Running() *OfflineClearParam {
	p.flag = 3
	return p
}

// Describe status of an offline task.
type OfflineTaskStatus int

// Return true if the task is still running.
func (s OfflineTaskStatus) IsRunning() bool {
	return s == 1
}

// Return true if the task has been done.
func (s OfflineTaskStatus) IsDone() bool {
	return s == 2
}

// Return true if the task has been failed.
func (s OfflineTaskStatus) IsFailed() bool {
	return s == -1
}

type OfflineTask struct {
	InfoHash string
	Name     string
	Url      string
	Status   OfflineTaskStatus
	Percent  int
	FileId   string
}

func (a *Agent) updateOfflineToken() (err error) {
	qs := core.NewQueryString().
		WithString("ct", "offline").
		WithString("ac", "space").
		WithInt64("_", time.Now().Unix())
	result := &internal.OfflineSpaceResult{}
	if err = a.hc.JsonApi(apiOfflineSpace, qs, nil, result); err != nil {
		return
	}
	// store to client
	if a.ot == nil {
		a.ot = &internal.OfflineToken{}
	}
	a.ot.Sign = result.Sign
	a.ot.Time = result.Time
	return nil
}

func (a *Agent) callOfflineApi(url string, form core.Form, result interface{}) (err error) {
	if a.ot == nil {
		if err = a.updateOfflineToken(); err != nil {
			return
		}
	}
	if form == nil {
		form = core.NewForm()
	}
	form.WithInt("uid", a.ui.UserId).
		WithString("sign", a.ot.Sign).
		WithInt64("time", a.ot.Time)
	err = a.hc.JsonApi(url, nil, form, result)
	// TODO: handle token expired error.
	return
}

// Get one page of offline tasks, the page size is hard-coded to 24 by upstream API.
func (a *Agent) OfflineList(page int) (tasks []*OfflineTask, next bool, err error) {
	if page < 1 {
		page = 1
	}
	form := core.NewForm().WithInt("page", page)
	result := &internal.OfflineListResult{}
	err = a.callOfflineApi(apiOfflineList, form, result)
	if err == nil && !result.State {
		err = errors.New(result.ErrorMsg)
	}
	if err != nil {
		return
	}
	tasks = make([]*OfflineTask, len(result.Tasks))
	for index, data := range result.Tasks {
		tasks[index] = &OfflineTask{
			InfoHash: data.InfoHash,
			Name:     data.Name,
			Url:      data.Url,
			Status:   OfflineTaskStatus(data.Status),
			Percent:  data.Precent,
			FileId:   data.FileId,
		}
	}
	next = result.Count-(result.Page*result.PageSize) > 0
	return
}

// Add one or more offline task by URL.
func (a *Agent) OfflineAdd(url ...string) (err error) {
	form, isSingle := core.NewForm(), len(url) == 1
	if isSingle {
		form.WithString("url", url[0])
		result := &internal.OfflineAddUrlResult{}
		err = a.callOfflineApi(apiOfflineAddUrl, form, result)
	} else {
		form.WithStrings("url", url)
		result := &internal.OfflineAddUrlsResult{}
		err = a.callOfflineApi(apiOfflineAddUrls, form, result)
	}
	// TODO: return add result
	return
}

// Delete some offline tasks.
// if "deleteFile" is true, the downloaded files will be deleted.
func (a *Agent) OfflineDelete(deleteFile bool, hash ...string) (err error) {
	form := core.NewForm().WithStrings("hash", hash)
	if deleteFile {
		form.WithInt("flag", 1)
	}
	result := &internal.OfflineBasicResult{}
	err = a.callOfflineApi(apiOfflineDelete, form, result)
	if err == nil && !result.State {
		err = errors.New(result.ErrorMsg)
	}
	return
}

// Clear offline tasks which specified by param, you can pass param as nil to
// clear done tasks without deleting downloaded files.
func (a *Agent) OfflineClear(param *OfflineClearParam) (err error) {
	if param == nil {
		param = (&OfflineClearParam{}).Done(false)
	}
	form := core.NewForm().
		WithInt("flag", param.flag)
	result := &internal.OfflineBasicResult{}
	err = a.callOfflineApi(apiOfflineClear, form, result)
	if err == nil && !result.State {
		err = errors.New(result.ErrorMsg)
	}
	return
}
