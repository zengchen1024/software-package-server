package app

func (s *softwarePkgService) GetPkgReviewDetail(pid string) (SoftwarePkgIssueDTO, error) {
	v, err := s.repo.FindSoftwarePkg(pid)
	if err != nil {
		return SoftwarePkgIssueDTO{}, err
	}

	return toSoftwarePkgIssueDTO(&v), nil
}

func (s *softwarePkgService) NewReviewComment(cmd *CmdToWriteSoftwarePkgReviewComment) error {
	return nil
}
