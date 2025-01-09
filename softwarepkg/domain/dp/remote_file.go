package dp

import (
	"errors"
	"strings"

	"github.com/opensourceways/software-package-server/allerror"
	"github.com/opensourceways/software-package-server/utils"
)

const (
	srpmSuffix = ".src.rpm"
	specSuffix = ".spec"
)

// RemoteFile
type RemoteFile interface {
	URL
	IsSRPM() bool
	Suffix() string
	CheckFile() error
}

type remoteFile struct {
	dpURL
}

func (r remoteFile) IsSRPM() bool {
	return r.Suffix() == srpmSuffix
}

func (r remoteFile) Suffix() string {
	if strings.HasSuffix(strings.ToLower(r.FileName()), specSuffix) {
		return specSuffix
	}

	return srpmSuffix
}

func (r remoteFile) CheckFile() error {
	return utils.CheckFile(r.URL(), "", 0)
}

func newRemoteFile(url, suffix string) (RemoteFile, error) {
	if _, err := NewURL(url); err != nil {
		return nil, err
	}

	if err := utils.CheckFile(url, "", 0); err != nil {
		return nil, allerror.New(allerror.ErrorCodeRemoteFileInvalid, err.Error())
	}

	v := dpURL(url)

	if !strings.HasSuffix(strings.ToLower(v.FileName()), suffix) {
		return nil, errors.New("unknown file")
	}

	return remoteFile{v}, nil
}

func NewSpecFile(url string) (RemoteFile, error) {
	return newRemoteFile(url, specSuffix)
}

func NewSRPMFile(url string) (RemoteFile, error) {
	return newRemoteFile(url, srpmSuffix)
}

func CreateRemoteFile(url string) (RemoteFile, error) {
	v := dpURL(url)

	name := strings.ToLower(v.FileName())
	if !strings.HasSuffix(name, specSuffix) && !strings.HasSuffix(name, srpmSuffix) {
		return nil, errors.New("unknown file")
	}

	return remoteFile{v}, nil
}
