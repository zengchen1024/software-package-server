package useradapterimpl

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/opensourceways/server-common-lib/utils"
)

type sigMaintainers map[string]bool

func (m sigMaintainers) isMaintainer(user string) bool {
	return m != nil && m[user]
}

func newSigMaintainers(maintainers []string) sigMaintainers {
	r := make(sigMaintainers, len(maintainers))

	for i := range maintainers {
		r[maintainers[i]] = true
	}

	return r
}

// sigData
type sigData struct {
	Data []struct {
		Maintainers []string `json:"maintainers"`
		SigName     string   `json:"sig_name"`
	} `json:"data"`

	maintainers map[string]sigMaintainers
	md5sum      string
}

func (s *sigData) isSigMaintainer(user, sig string) bool {
	if s != nil && s.maintainers != nil {
		v, ok := s.maintainers[sig]

		return ok && v.isMaintainer(user)
	}

	return false
}

func (s *sigData) init(md5sum string) {
	s.md5sum = md5sum

	if len(s.Data) == 0 {
		return
	}

	items := s.Data
	s.maintainers = make(map[string]sigMaintainers, len(items))

	for i := range items {
		s.maintainers[items[i].SigName] = newSigMaintainers(items[i].Maintainers)
	}
}

// sigLoader
type sigLoader struct {
	cli  utils.HttpClient
	link string
}

func (l *sigLoader) read() (data []byte, err error) {
	req, err := http.NewRequest(http.MethodGet, l.link, nil)
	if err == nil {
		data, _, err = l.cli.Download(req)
	}

	return
}

func (l *sigLoader) Load(old interface{}) (interface{}, error) {
	if old == nil {
		return l.load(nil)
	}

	return l.load(old.(*sigData))
}

func (l *sigLoader) load(old *sigData) (interface{}, error) {
	b, err := l.read()
	if err != nil || len(b) == 0 {
		return nil, err
	}

	md5sum := fmt.Sprintf("%x", md5.Sum(b))
	if old != nil && old.md5sum == md5sum {
		return nil, nil
	}

	s := new(sigData)
	if err = json.Unmarshal(b, s); err == nil {
		s.init(md5sum)
	}

	return s, nil
}
