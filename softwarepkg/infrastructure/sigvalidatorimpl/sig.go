package sigvalidatorimpl

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/opensourceways/server-common-lib/utils"

	"github.com/opensourceways/software-package-server/softwarepkg/domain/sigvalidator"
)

type sigDetail = sigvalidator.Sig

// sigData
type sigData struct {
	Data []sigDetail `json:"data"`

	sigs map[string]bool
}

func (s *sigData) getAll() (info []sigDetail) {
	if s == nil {
		return nil
	}

	return s.Data
}

func (s *sigData) hasSig(sig string) bool {
	return s != nil && s.sigs != nil && s.sigs[sig]
}

func (s *sigData) init() {
	s.sigs = make(map[string]bool, len(s.Data))

	for i := range s.Data {
		s.sigs[s.Data[i].SigNames] = true
	}
}

// sigLoader
type sigLoader struct {
	cli utils.HttpClient
}

func (l *sigLoader) read(link string) (data []byte, err error) {
	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err == nil {
		data, _, err = l.cli.Download(req)
	}

	return
}

func (l *sigLoader) load(link, OldMD5Sum string) (s *sigData, md5sum string, err error) {
	b, err := l.read(link)
	if err != nil {
		return
	}

	if md5sum = fmt.Sprintf("%x", md5.Sum(b)); OldMD5Sum == md5sum {
		return
	}

	s = new(sigData)
	if err = json.Unmarshal(b, s); err == nil {
		s.init()
	}

	return
}
