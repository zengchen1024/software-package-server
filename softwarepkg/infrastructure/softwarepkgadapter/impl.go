package softwarepkgadapter

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	commonrepo "github.com/opensourceways/software-package-server/common/domain/repository"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/repository"
)

func NewsoftwarePkgAdapter(dao dao) *softwarePkgAdapter {
	return &softwarePkgAdapter{dao}
}

type softwarePkgAdapter struct {
	dao dao
}

func (impl *softwarePkgAdapter) Add(pkg *domain.SoftwarePkg) error {
	do := new(softwarePkgDO)
	toSoftwarePkgDO(pkg, do)

	doc, err := do.toDoc()
	if err != nil {
		return err
	}
	doc[fieldVersion] = 0

	pkg.Id, err = impl.dao.InsertDocIfNotExists(do.docFilter(), doc)
	if err != nil && impl.dao.IsDocExists(err) {
		err = commonrepo.NewErrorDuplicateCreating(err)
	}

	return err
}

func (impl *softwarePkgAdapter) Find(pid string) (domain.SoftwarePkg, int, error) {
	return impl.find(pid, false)
}

func (impl *softwarePkgAdapter) FindAndIgnoreReview(pid string) (domain.SoftwarePkg, int, error) {
	return impl.find(pid, true)
}

func (impl *softwarePkgAdapter) find(pid string, ignoreReview bool) (pkg domain.SoftwarePkg, version int, err error) {
	filter, err := impl.dao.DocIdFilter(pid)
	if err != nil {
		return
	}

	var doc softwarePkgDO

	if ignoreReview {
		err = impl.dao.GetDoc(filter, bson.M{fieldReviews: 0}, &doc)
	} else {
		err = impl.dao.GetDoc(filter, nil, &doc)
	}

	if err != nil && impl.dao.IsDocNotExists(err) {
		err = commonrepo.NewErrorResourceNotFound(err)

		return
	}

	err = doc.toDomain(&pkg)
	version = doc.Version

	return
}

func (impl *softwarePkgAdapter) Save(pkg *domain.SoftwarePkg, version int) error {
	return impl.save(pkg, version, false)
}

func (impl *softwarePkgAdapter) SaveAndIgnoreReview(pkg *domain.SoftwarePkg, version int) error {
	return impl.save(pkg, version, true)
}

func (impl *softwarePkgAdapter) save(pkg *domain.SoftwarePkg, version int, ignoreReview bool) error {
	filter, err := impl.dao.DocIdFilter(pkg.Id)
	if err != nil {
		return err
	}

	if ignoreReview {
		// ignore the reviews
		pkg.Reviews = nil
	}

	do := new(softwarePkgDO)
	toSoftwarePkgDO(pkg, do)

	doc, err := do.toDoc()
	if err != nil {
		return err
	}

	if ignoreReview {
		delete(doc, fieldReviews)
	}

	err = impl.dao.UpdateDoc(filter, doc, version)
	if err != nil && impl.dao.IsDocNotExists(err) {
		err = commonrepo.NewErrorConcurrentUpdating(err)
	}

	return err
}

func (impl *softwarePkgAdapter) FindAll(opt *repository.OptToFindSoftwarePkgs) (
	[]repository.SoftwarePkgInfo, int64, error,
) {
	var err error
	var lastId primitive.ObjectID

	if opt.LastId != "" {
		if lastId, err = primitive.ObjectIDFromHex(opt.LastId); err != nil {
			return nil, 0, err
		}
	}

	filter := impl.listFilter(opt)

	// total
	var total int64
	if opt.Count {
		if total, err = impl.dao.Count(filter); err != nil {
			return nil, 0, err
		}
	}

	// list
	if opt.LastId != "" {
		filter[fieldIndex] = bson.M{mongodbCmdLt: lastId}
	}

	project := bson.M{
		fieldSig:          1,
		fieldPhase:        1,
		fieldCIStatus:     1,
		fieldImporter:     1,
		fieldAppliedAt:    1,
		fieldBasicDesc:    1,
		fieldPrimaryKey:   1,
		fieldRepoPlatform: 1,
	}

	var docs []softwarePkgDO

	err = impl.dao.Paginate(
		filter, project, bson.M{fieldIndex: -1},
		int64(opt.PageNum), int64(opt.CountPerPage), &docs,
	)
	if err != nil || len(docs) == 0 {
		return nil, total, err
	}

	r := make([]repository.SoftwarePkgInfo, len(docs))
	for i := range docs {
		if err = docs[i].toSoftwarePkgInfo(&r[i]); err != nil {
			return nil, 0, err
		}
	}

	return r, total, nil
}

func (impl *softwarePkgAdapter) listFilter(opt *repository.OptToFindSoftwarePkgs) bson.M {
	filter := bson.M{}
	if opt.Phase != nil {
		filter[fieldPhase] = opt.Phase.PackagePhase()
	}

	if opt.Platform != nil {
		filter[fieldRepoPlatform] = opt.Platform.PackagePlatform()
	}

	if opt.Importer != nil {
		filter[fieldImporter] = opt.Importer.Account()
	}

	if opt.PkgName != nil {
		filter[fieldPrimaryKey] = impl.dao.LikeFilter(opt.PkgName.PackageName(), true)
	}

	return filter
}
