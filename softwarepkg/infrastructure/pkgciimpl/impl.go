package pkgciimpl

import (
	"fmt"

	"github.com/opensourceways/robot-gitee-lib/client"
	libutils "github.com/opensourceways/server-common-lib/utils"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/yaml"

	"github.com/opensourceways/software-package-server/softwarepkg/domain"
)

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
		cfg.CloneScript,
		cfg.WorkDir,
		user.User,
		user.Email,
		cfg.CIRepo.Repo,
		cfg.CIRepo.cloneURL(user),
	}

	_, err, _ := libutils.RunCmd(params...)

	return err
}

func PkgCI() *pkgCIImpl {
	return instance
}

type softwarePkgInfo struct {
	PkgId   string `json:"pkg_id"`
	PkgName string `json:"pkg_name"`
	Service string `json:"service"`
}

func (s *softwarePkgInfo) toYaml() ([]byte, error) {
	return yaml.Marshal(s)
}

// pkgCIImpl
type pkgCIImpl struct {
	cli         client.Client
	cfg         Config
	ciRepoDir   string
	pkgInfoFile string
}

func (impl *pkgCIImpl) StartNewCI(pkg *domain.SoftwarePkg) (int, error) {
	if pkg.CI.Id > 0 {
		impl.closePR(pkg.CI.Id)
	}

	name := pkg.Basic.Name.PackageName()

	pr, err := impl.cli.CreatePullRequest(
		impl.cfg.CIRepo.Org,
		impl.cfg.CIRepo.Repo,
		fmt.Sprintf("test for new package: %s", name), pkg.Id,
		name,
		impl.cfg.TargetBranch,
		true,
	)
	if err != nil {
		return 0, err
	}

	return int(pr.Number), nil
}

func (impl *pkgCIImpl) ClearCI(pkg *domain.SoftwarePkg) error {
	if pkg.CI.Id > 0 {
		impl.closePR(pkg.CI.Id)
	}

	// clear branch

	return nil
}

func (impl *pkgCIImpl) DownloadPkgCode(pkg *domain.SoftwarePkg) error {
	branch := pkg.Basic.Name.PackageName()
	if err := impl.createBranch(pkg, branch); err != nil {
		return err
	}

	if v := &pkg.Code.Spec; v.Dirty {
		v.Dirty = false
		v.DownloadAddr = nil // TODO
	}

	if v := &pkg.Code.SRPM; v.Dirty {
		v.Dirty = false
		v.DownloadAddr = nil // TODO
	}

	return nil
}

func (impl *pkgCIImpl) closePR(id int) error {
	return impl.cli.ClosePR(impl.cfg.CIRepo.Org, impl.cfg.CIRepo.Repo, int32(id))
}

func (impl *pkgCIImpl) createBranch(info *domain.SoftwarePkg, branch string) error {
	cfg := &impl.cfg

	code := &info.Code
	spec := "-"
	if code.Spec.Dirty {
		spec = code.Spec.Src.URL()
	}

	srpm := "-"
	if code.SRPM.Dirty {
		srpm = code.SRPM.Src.URL()
	}

	params := []string{
		cfg.PRScript,
		impl.ciRepoDir,
		cfg.GitUser.Token,
		cfg.TargetBranch,
		branch,
		spec,
		srpm,
	}

	return impl.runcmd(params)
}

func (impl *pkgCIImpl) runcmd(params []string) error {
	out, err, _ := libutils.RunCmd(params...)
	if err != nil {
		logrus.Errorf(
			"run create pull request shell, err=%s, out=%s, params=%v",
			err.Error(), out, params[:len(params)-1],
		)
	}

	return err
}
