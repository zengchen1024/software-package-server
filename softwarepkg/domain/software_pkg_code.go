package domain

import (
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/utils"
)

// SoftwarePkgCode
type SoftwarePkgCode struct {
	Spec SoftwarePkgCodeInfo
	SRPM SoftwarePkgCodeInfo
}

func (code *SoftwarePkgCode) isReady() bool {
	return code.Spec.isReady() && code.SRPM.isReady()
}

func (code *SoftwarePkgCode) update(spec, srpm dp.URL) {
	if spec != nil {
		code.Spec.update(spec)
	}

	if srpm != nil {
		code.SRPM.update(srpm)
	}
}

func (code *SoftwarePkgCode) filesToDownload() []SoftwarePkgCodeFile {
	r := []SoftwarePkgCodeFile{}

	if v := code.Spec.fileToDownload(); v != nil {
		r = append(r, *v)
	}

	if v := code.SRPM.fileToDownload(); v != nil {
		r = append(r, *v)
	}

	return r
}

func (code *SoftwarePkgCode) saveDownloadedFiles(files []SoftwarePkgCodeFile) bool {
	spec := false
	srpm := false

	for i := range files {
		item := &files[i]

		if !spec && code.Spec.saveDownloadedFile(item) {
			spec = true
		}

		if !srpm && code.SRPM.saveDownloadedFile(item) {
			srpm = true
		}
	}

	return spec || srpm
}

// SoftwarePkgCodeInfo
type SoftwarePkgCodeInfo struct {
	SoftwarePkgCodeFile

	Dirty bool // if true, the code should be updated.
	//Reason string // the reason why can't download the code file
}

func (f *SoftwarePkgCodeInfo) isReady() bool {
	return f.DownloadAddr != nil && !f.Dirty
}

func (f *SoftwarePkgCodeInfo) update(src dp.URL) {
	f.Src = src
	f.DownloadAddr = nil
	f.Dirty = true
	// f.Reason = ""
	f.UpdatedAt = utils.Now()
}

func (f *SoftwarePkgCodeInfo) saveDownloadedFile(file *SoftwarePkgCodeFile) bool {
	if !f.isReady() && f.isSame(file) {
		f.DownloadAddr = file.DownloadAddr
		f.Dirty = false

		return true
	}

	return false
}

func (f *SoftwarePkgCodeInfo) fileToDownload() *SoftwarePkgCodeFile {
	if !f.isReady() {
		return &f.SoftwarePkgCodeFile
	}

	return nil
}

// SoftwarePkgCodeFile
type SoftwarePkgCodeFile struct {
	Src          dp.URL // Src is the url user inputed
	UpdatedAt    int64  // UpdatedAt is the time when user changes the Src or wants to reload
	DownloadAddr dp.URL
}

func (f *SoftwarePkgCodeFile) FileName() string {
	return f.Src.FileName()
}

func (f *SoftwarePkgCodeFile) isSame(f1 *SoftwarePkgCodeFile) bool {
	return f.UpdatedAt == f1.UpdatedAt && f.Src.URL() == f1.Src.URL()
}
