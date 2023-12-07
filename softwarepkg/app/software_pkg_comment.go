package app

import (
	"github.com/opensourceways/software-package-server/allerror"
	commonrepo "github.com/opensourceways/software-package-server/common/domain/repository"
	"github.com/opensourceways/software-package-server/softwarepkg/domain"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/repository"
	"github.com/opensourceways/software-package-server/softwarepkg/domain/translation"
)

type SoftwarePkgCommentAppService interface {
	NewReviewComment(*CmdToWriteSoftwarePkgReviewComment) error

	ListComments(pkgId string) ([]SoftwarePkgReviewCommentDTO, error)

	TranslateReviewComment(*CmdToTranslateReviewComment) (
		dto TranslatedReveiwCommentDTO, err error,
	)
}

func NewSoftwarePkgCommentAppService(
	repo repository.SoftwarePkg,
	translation translation.Translation,
	commentRepo repository.SoftwarePkgComment,
) *softwarePkgCommentAppService {
	return &softwarePkgCommentAppService{
		repo:        repo,
		translation: translation,
		commentRepo: commentRepo,
	}
}

type softwarePkgCommentAppService struct {
	repo        repository.SoftwarePkg
	translation translation.Translation
	commentRepo repository.SoftwarePkgComment
}

func (s *softwarePkgCommentAppService) NewReviewComment(cmd *CmdToWriteSoftwarePkgReviewComment) error {
	pkg, _, err := s.repo.FindAndIgnoreReview(cmd.PkgId)
	if err != nil {
		return parseErrorForFindingPkg(err)
	}

	if err := pkg.CanAddReviewComment(); err != nil {
		return err
	}

	// TODO: there is a critical case that the comment can't be added now
	comment := domain.NewSoftwarePkgReviewComment(cmd.Author, cmd.Content)
	return s.commentRepo.AddReviewComment(cmd.PkgId, &comment)
}

func (s *softwarePkgCommentAppService) ListComments(pkgId string) ([]SoftwarePkgReviewCommentDTO, error) {
	v, err := s.commentRepo.FindReviewComments(pkgId)
	if err != nil || len(v) == 0 {

	}

	return toSoftwarePkgReviewCommentDTOs(v), nil
}

func (s *softwarePkgCommentAppService) TranslateReviewComment(
	cmd *CmdToTranslateReviewComment,
) (dto TranslatedReveiwCommentDTO, err error) {
	v, err := s.commentRepo.FindTranslatedReviewComment(cmd)
	if err == nil {
		dto.Content = v.Content

		return
	}

	if !commonrepo.IsErrorResourceNotFound(err) {
		return
	}

	// translate it
	comment, err := s.commentRepo.FindReviewComment(cmd.PkgId, cmd.CommentId)
	if err != nil {
		if commonrepo.IsErrorResourceNotFound(err) {
			err = allerror.NewNotFound(allerror.ErrorCodePkgCommentNotFound, err.Error())
		}

		return
	}

	content, err := s.translation.Translate(
		comment.Content.ReviewComment(), cmd.Language,
	)
	if err != nil {
		return
	}

	dto.Content = content

	// save the translated one
	translated := domain.NewSoftwarePkgTranslatedReviewComment(
		&comment, content, cmd.Language,
	)
	_ = s.commentRepo.AddTranslatedReviewComment(cmd.PkgId, &translated)

	return
}
