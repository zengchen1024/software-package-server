package maintainerimpl

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/opensourceways/server-common-lib/utils"
)

// sigData
type sigData struct {
	Data []struct {
		Maintaines []string `json:"maintainers"`
	} `json:"data"`

	maintainers map[string]bool
	md5sum      string
}

func (s *sigData) hasMaintainer(v string) bool {
	return s != nil && s.maintainers != nil && s.maintainers[v]
}

func (s *sigData) init(md5sum string) {
	s.md5sum = md5sum

	if len(s.Data) == 0 {
		return
	}

	items := s.Data[0].Maintaines
	s.maintainers = make(map[string]bool, len(items))

	for i := range items {
		s.maintainers[items[i]] = true
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

func (l *sigLoader) load(old *sigData) (s *sigData, err error) {
	b, err := l.read()
	if err != nil {
		return
	}

	md5sum := fmt.Sprintf("%x", md5.Sum(b))
	if old != nil && old.md5sum == md5sum {
		return
	}

	s = new(sigData)
	if err = json.Unmarshal(b, s); err == nil {
		s.init(md5sum)
	}

	return
}
