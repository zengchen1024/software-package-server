package pkgciimpl

import (
	"fmt"
	"strings"

	"github.com/opensourceways/robot-gitee-lib/client"
	libutils "github.com/opensourceways/server-common-lib/utils"
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/yaml"

	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
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

func (impl *pkgCIImpl) closePR(id int) error {
	return impl.cli.ClosePR(impl.cfg.CIRepo.Org, impl.cfg.CIRepo.Repo, int32(id))
}

func (impl *pkgCIImpl) Download(files []domain.SoftwarePkgCodeSourceFile, name dp.PackageName) error {
	if len(files) == 0 {
		return nil
	}

	other := []string{"-", "-", "-", "-"}
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
		cfg.PRScript,
		impl.ciRepoDir,
		cfg.GitUser.Token,
		cfg.TargetBranch,
		name.PackageName(),
	}

	params = append(params, other...)

	out, err, _ := libutils.RunCmd(params...)
	if err != nil {
		return err
	}

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

		v, err := cfg.CIRepo.fileAddr(f, lfs)
		if err != nil {
			return err
		}

		item.DownloadAddr = v
	}

	return nil
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
