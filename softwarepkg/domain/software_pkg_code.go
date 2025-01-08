package domain

import (
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/utils"
)

// SoftwarePkgCode
type SoftwarePkgCode struct {
	Spec SoftwarePkgCodeFile
	SRPM SoftwarePkgCodeFile
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

func (code *SoftwarePkgCode) filesToDownload() []SoftwarePkgCodeSourceFile {
	r := []SoftwarePkgCodeSourceFile{}

	if v := code.Spec.fileToDownload(); v != nil {
		r = append(r, *v)
	}

	if v := code.SRPM.fileToDownload(); v != nil {
		r = append(r, *v)
	}

	return r
}

func (code *SoftwarePkgCode) saveDownloadedFiles(files []SoftwarePkgCodeSourceFile) (bool, bool) {
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

	return spec || srpm, code.isReady()
}

// SoftwarePkgCodeFile
type SoftwarePkgCodeFile struct {
	SoftwarePkgCodeSourceFile

	Dirty bool // if true, the code should be updated.

	//Reason string // the reason why can't download the code file
}

func (f *SoftwarePkgCodeFile) isReady() bool {
	return f.DownloadAddr != nil && !f.Dirty
}

func (f *SoftwarePkgCodeFile) update(src dp.URL) {
	f.Src = src
	f.Dirty = true
	f.UpdatedAt = utils.Now()
	f.DownloadAddr = nil
}

func (f *SoftwarePkgCodeFile) saveDownloadedFile(file *SoftwarePkgCodeSourceFile) bool {
	if !f.isReady() && f.isSame(file) {
		f.DownloadAddr = file.DownloadAddr
		f.Dirty = false

		return true
	}

	return false
}

func (f *SoftwarePkgCodeFile) fileToDownload() *SoftwarePkgCodeSourceFile {
	if !f.isReady() {
		return &f.SoftwarePkgCodeSourceFile
	}

	return nil
}

// SoftwarePkgCodeSourceFile
type SoftwarePkgCodeSourceFile struct {
	Src          dp.URL // Src is the url user inputed
	UpdatedAt    int64  // UpdatedAt is the time when user changes the Src or wants to reload
	DownloadAddr dp.URL
}

func (f *SoftwarePkgCodeSourceFile) FileName() string {
	return f.Src.FileName()
}

func (f *SoftwarePkgCodeSourceFile) IsSRPM() bool {
	return dp.IsSRPM(f.FileName())
}

func (f *SoftwarePkgCodeSourceFile) FormatedFileName(name dp.PackageName) string {
	if f.IsSRPM() {
		return name.PackageName() + dp.SRPMSuffix
	}

	return name.PackageName() + dp.SpecSuffix
}

func (f *SoftwarePkgCodeSourceFile) isSame(f1 *SoftwarePkgCodeSourceFile) bool {
	return f.UpdatedAt == f1.UpdatedAt && f.Src.URL() == f1.Src.URL()
}
