package repositoryimpl

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/plugin/optimisticlock"

	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
	"github.com/opensourceways/software-package-server/utils"
)

const (
	fieldId              = "uuid"
	fieldPhase           = "phase"
	fieldVersion         = "version"
	fieldImporter        = "importer"
	fieldAppliedAt       = "applied_at"
	fieldUpdatedAt       = "updated_at"
	fieldApprovedby      = "approvedby"
	fieldRejectedby      = "rejectedby"
	fieldPackageName     = "package_name"
	fieldPackagePlatform = "package_platform"
)

func (s softwarePkgBasic) toSoftwarePkgBasicDO(pkg *domain.SoftwarePkg, do *SoftwarePkgBasicDO) (err error) {
	email, err := toEmailDO(pkg.Importer.Email)
	if err != nil {
		return err
	}

	code := &pkg.Code
	basic := &pkg.Basic

	*do = SoftwarePkgBasicDO{
		Id:              uuid.New(),
		PackageName:     basic.Name.PackageName(),
		Importer:        pkg.Importer.Account.Account(),
		ImporterEmail:   email,
		Phase:           pkg.Phase.PackagePhase(),
		CIPRNum:         pkg.CI.PRNum,
		CIStatus:        pkg.CI.Status.PackageCIStatus(),
		SpecURL:         code.Spec.Src.URL(),
		SrcRPMURL:       code.SRPM.Path.URL(),
		Upstream:        basic.Upstream.URL(),
		PackageDesc:     basic.Desc.PackageDesc(),
		PackagePlatform: pkg.Repo.Platform.PackagePlatform(),
		RepoLink:        pkg.Repo.Link.URL(),
		Sig:             pkg.Sig.ImportingPkgSig(),
		ReasonToImport:  basic.Reason.ReasonToImportPkg(),
		AppliedAt:       pkg.AppliedAt,
		UpdatedAt:       pkg.AppliedAt,
		ApprovedBy:      toStringArray(pkg.ApprovedBy),
		RejectedBy:      toStringArray(pkg.RejectedBy),
	}

	if pkg.CommunityPR != nil {
		do.RelevantPR = pkg.CommunityPR.URL()
	}

	return
}

type SoftwarePkgBasicDO struct {
	// must set "uuid" as the name of column
	Id              uuid.UUID              `gorm:"column:uuid;type:uuid"                           json:"-"`
	Sig             string                 `gorm:"column:sig"                                      json:"sig"`
	Phase           string                 `gorm:"column:phase"                                    json:"phase"`
	SpecURL         string                 `gorm:"column:spec_url"                                 json:"spec_url"`
	Importer        string                 `gorm:"column:importer"                                 json:"importer"`
	RepoLink        string                 `gorm:"column:repo_link"                                json:"repo_link"`
	CIStatus        string                 `gorm:"column:ci_status"                                json:"ci_status"`
	Upstream        string                 `gorm:"column:upstream"                                 json:"upstream"`
	SrcRPMURL       string                 `gorm:"column:src_rpm_url"                              json:"src_rpm_url"`
	RelevantPR      string                 `gorm:"column:relevant_pr"                              json:"relevant_pr"`
	PackageName     string                 `gorm:"column:package_name"                             json:"package_name"`
	PackageDesc     string                 `gorm:"column:package_desc"                             json:"package_desc"`
	ImporterEmail   string                 `gorm:"column:importer_email"                           json:"importer_email"`
	ReasonToImport  string                 `gorm:"column:reason_to_import"                         json:"reason_to_import"`
	PackagePlatform string                 `gorm:"column:package_platform"                         json:"package_platform"`
	CIPRNum         int                    `gorm:"column:ci_pr_num"                                json:"ci_pr_num"`
	AppliedAt       int64                  `gorm:"column:applied_at"                               json:"applied_at"`
	UpdatedAt       int64                  `gorm:"column:updated_at"                               json:"updated_at"`
	Version         optimisticlock.Version `gorm:"column:version"                                  json:"-"`
	ApprovedBy      pq.StringArray         `gorm:"column:approvedby;type:text[];default:'{}'"      json:"-"`
	RejectedBy      pq.StringArray         `gorm:"column:rejectedby;type:text[];default:'{}'"      json:"-"`
}

func (do *SoftwarePkgBasicDO) toMap() (map[string]any, error) {
	v, err := json.Marshal(do)
	if err != nil {
		return nil, err
	}

	var r map[string]any
	if err = json.Unmarshal(v, &r); err != nil {
		return nil, err
	}

	// fieldVersion
	r[fieldVersion] = gorm.Expr(fieldVersion+" + ?", 1)

	// fieldApprovedby
	s, err := marshalStringArray(do.ApprovedBy)
	if err != nil {
		return nil, err
	}
	r[fieldApprovedby] = s

	// fieldRejectedby
	s, err = marshalStringArray(do.RejectedBy)
	if err != nil {
		return nil, err
	}
	r[fieldRejectedby] = s

	return r, err
}

func (do *SoftwarePkgBasicDO) toSoftwarePkg() (info domain.SoftwarePkg, err error) {
	info.Id = do.Id.String()

	if err = do.toSoftwarePkgApplication(&info); err != nil {
		return
	}

	if do.RelevantPR != "" {
		if info.CommunityPR, err = dp.NewURL(do.RelevantPR); err != nil {
			return
		}
	}

	if info.Importer.Account, err = dp.NewAccount(do.Importer); err != nil {
		return
	}

	if info.Importer.Email, err = toEmail(do.ImporterEmail); err != nil {
		return
	}

	if info.Phase, err = dp.NewPackagePhase(do.Phase); err != nil {
		return
	}

	info.AppliedAt = do.AppliedAt

	if info.ApprovedBy, err = do.toAccounts(do.ApprovedBy); err != nil {
		return
	}

	if info.CI.Status, err = dp.NewPackageCIStatus(do.CIStatus); err != nil {
		return
	}

	info.CI.PRNum = do.CIPRNum

	info.RejectedBy, err = do.toAccounts(do.RejectedBy)

	return
}

func (do *SoftwarePkgBasicDO) toAccounts(v []string) (r []domain.SoftwarePkgApprover, err error) {
	if len(v) == 0 {
		return
	}

	r = make([]domain.SoftwarePkgApprover, len(v))
	for i := range v {
		if r[i], err = domain.StringToSoftwarePkgApprover(v[i]); err != nil {
			return
		}
	}

	return
}

func (do *SoftwarePkgBasicDO) toSoftwarePkgApplication(pkg *domain.SoftwarePkg) (err error) {
	basic := &pkg.Basic

	if basic.Name, err = dp.NewPackageName(do.PackageName); err != nil {
		return
	}

	if basic.Reason, err = dp.NewReasonToImportPkg(do.ReasonToImport); err != nil {
		return
	}

	if basic.Desc, err = dp.NewPackageDesc(do.PackageDesc); err != nil {
		return
	}

	if basic.Upstream, err = dp.NewURL(do.Upstream); err != nil {
		return
	}

	if pkg.Repo.Platform, err = dp.NewPackagePlatform(do.PackagePlatform); err != nil {
		return
	}

	if pkg.Repo.Link, err = dp.NewURL(do.RepoLink); err != nil {
		return
	}

	if pkg.Sig, err = dp.NewImportingPkgSig(do.Sig); err != nil {
		return
	}

	code := &pkg.Code
	if code.SRPM.Src, err = dp.NewURL(do.SrcRPMURL); err != nil {
		return
	}

	code.Spec.Src, err = dp.NewURL(do.SpecURL)

	return
}

func toStringArray(v []domain.SoftwarePkgApprover) (arr pq.StringArray) {
	arr = make(pq.StringArray, len(v))
	for i := range v {
		arr[i] = v[i].String()
	}

	return
}

func marshalStringArray(sa pq.StringArray) (string, error) {
	v, err := sa.Value()
	if err != nil {
		return "", err
	}

	if v != nil {
		if s, ok := v.(string); ok {
			return s, nil
		}
	}

	return "{}", nil
}

func toEmailDO(e dp.Email) (string, error) {
	return utils.Encryption.Encrypt([]byte(e.Email()))
}

func toEmail(e string) (dp.Email, error) {
	v, err := utils.Encryption.Decrypt(e)
	if err != nil {
		return nil, err
	}

	return dp.NewEmail(string(v))
}
