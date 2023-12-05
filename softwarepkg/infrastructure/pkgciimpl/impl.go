package pkgciimpl

import (
	"fmt"
	"strings"

	"github.com/opensourceways/robot-gitee-lib/client"
	libutils "github.com/opensourceways/server-common-lib/utils"

	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

const codeChanged = "code_changed!!"

var instance *pkgCIImpl

func Init(cfg *Config) error {
	if err := cloneRepo(cfg); err != nil {
		return err
	}

	instance = &pkgCIImpl{
		cli: client.NewClient(func() []byte {
			return []byte(cfg.GitUser.Token)
		}),
		cfg:         *cfg,
		ciRepoDir:   cfg.WorkDir + "/" + cfg.CIRepo.Repo,
		pkgInfoFile: cfg.WorkDir + "/pkginfo.yaml",
	}

	return nil
}

func cloneRepo(cfg *Config) error {
	user := &cfg.GitUser

	params := []string{
		cfg.InitScript,
		cfg.WorkDir,
		user.User,
		user.Email,
		cfg.CIRepo.Repo,
		cfg.CIRepo.cloneURL(user),
	}

	if out, err, _ := libutils.RunCmd(params...); err != nil {
		return fmt.Errorf("%s, %s", string(out), err.Error())
	}

	return nil
}

func PkgCI() *pkgCIImpl {
	return instance
}

// pkgCIImpl
type pkgCIImpl struct {
	cli         client.Client
	cfg         Config
	ciRepoDir   string
	pkgInfoFile string
}

func (impl *pkgCIImpl) StartNewCI(pkg *domain.SoftwarePkg) (int, error) {
	if v := pkg.CIId(); v > 0 {
		impl.closePR(v)
	}

	name := pkg.PackageName().PackageName()
	cfg := &impl.cfg.CIRepo

	pr, err := impl.cli.CreatePullRequest(
		cfg.Org, cfg.Repo,
		fmt.Sprintf("test for new package: %s", name), pkg.Id,
		name, cfg.MainBranch, true,
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

	other := []string{"-", "-", "-", "-", codeChanged}
	specIndex, srpmIndex := 0, 2
	for _, item := range files {
		i := specIndex
		v := item.FileName()
		if dp.IsSRPM(v) {
			i = srpmIndex
		}

		other[i] = item.Src.URL()
		other[i+1] = v
	}

	cfg := &impl.cfg
	params := []string{
		cfg.DownloadScript,
		impl.ciRepoDir,
		cfg.GitUser.Token,
		cfg.CIRepo.MainBranch,
		name.PackageName(),
	}

	params = append(params, other...)

	out, err, _ := libutils.RunCmd(params...)
	if err != nil {
		return false, err
	}

	changed := strings.Contains(string(out), codeChanged)

	// fetch download addr
	for _, item := range files {
		f := ""
		lfs := false
		if dp.IsSRPM(item.FileName()) {
			f = name.PackageName() + ".src.rpm"
			lfs = strings.Contains(string(out), f)
		} else {
			f = name.PackageName() + ".spec"
		}

		v, err := cfg.CIRepo.fileAddr(name, f, lfs)
		if err != nil {
			return changed, err
		}

		item.DownloadAddr = v
	}

	return changed, nil
}
