package pkgciimpl

import (
	"fmt"
	"strings"

	"github.com/opensourceways/robot-gitee-lib/client"
	libutils "github.com/opensourceways/server-common-lib/utils"

	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

const (
	codeChangedTag = "code_changed!!"
	srpmFileLFSTag = "srpm_file_is_lfs!!"
)

func Init(cfg *Config) (*pkgCIImpl, error) {
	if err := cloneRepo(cfg); err != nil {
		return nil, err
	}

	return &pkgCIImpl{
		cli: client.NewClient(func() []byte {
			return []byte(cfg.CIRepo.Token)
		}),
		cfg:       *cfg,
		ciRepoDir: cfg.WorkDir + "/" + cfg.CIRepo.Repo,
	}, nil
}

func cloneRepo(cfg *Config) error {
	repo := &cfg.CIRepo

	params := []string{
		cfg.InitScript,
		cfg.WorkDir,
		repo.Owner,
		repo.Email,
		cfg.CIRepo.Repo,
		cfg.CIRepo.cloneURL(),
	}

	if out, err, _ := libutils.RunCmd(params...); err != nil {
		return fmt.Errorf("%s, %s", string(out), err.Error())
	}

	return nil
}

// pkgCIImpl
type pkgCIImpl struct {
	cli       client.Client
	cfg       Config
	ciRepoDir string
}

func (impl *pkgCIImpl) StartNewCI(pkg *domain.SoftwarePkg) (int, error) {
	if v := pkg.CIId(); v > 0 {
		impl.closePR(v)
	}

	name := pkg.PackageName().PackageName()
	repo := &impl.cfg.CIRepo

	pr, err := impl.cli.CreatePullRequest(
		repo.Org, repo.Repo,
		fmt.Sprintf("test for new package: %s", name), pkg.Id,
		repo.Owner+":"+name, repo.MainBranch, true,
	)
	if err != nil {
		return 0, err
	}

	return int(pr.Number), nil
}

func (impl *pkgCIImpl) Clear(pkg *domain.SoftwarePkg) error {
	if v := pkg.CIId(); v > 0 {
		impl.closePR(v)
	}

	// clear branch

	cfg := &impl.cfg
	params := []string{
		cfg.DownloadScript,
		impl.ciRepoDir,
		cfg.CIRepo.MainBranch,
		pkg.PackageName().PackageName(),
	}

	_, err, _ := libutils.RunCmd(params...)

	return err
}

func (impl *pkgCIImpl) closePR(id int) error {
	return impl.cli.ClosePR(impl.cfg.CIRepo.Org, impl.cfg.CIRepo.Repo, int32(id))
}

func (impl *pkgCIImpl) Download(files []domain.SoftwarePkgCodeSourceFile, name dp.PackageName) (bool, error) {
	if len(files) == 0 {
		return false, nil
	}

	other := []string{"-", "-", "-", "-", codeChangedTag, srpmFileLFSTag}
	specIndex, srpmIndex := 0, 2
	for i := range files {
		item := &files[i]

		i := specIndex
		if item.IsSRPM() {
			i = srpmIndex
		}

		other[i] = item.Src.URL()
		other[i+1] = item.FileName()
	}

	repo := &impl.cfg.CIRepo
	params := []string{
		impl.cfg.DownloadScript,
		impl.ciRepoDir,
		repo.Token,
		repo.MainBranch,
		name.PackageName(),
	}

	params = append(params, other...)

	out, err, _ := libutils.RunCmd(params...)
	if err != nil {
		return false, err
	}

	outStr := string(out)
	changed := strings.Contains(outStr, codeChangedTag)

	// fetch download addr
	for i := range files {
		item := &files[i] // need update files item by pointer.

		v, err := repo.fileAddr(
			name, item.FormatedFileName(name),
			item.IsSRPM() && strings.Contains(outStr, srpmFileLFSTag),
		)
		if err != nil {
			return changed, err
		}

		item.DownloadAddr = v
	}

	return changed, nil
}
