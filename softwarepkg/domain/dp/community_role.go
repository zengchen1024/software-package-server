package dp

import (
	"errors"
	"strings"
)

const (
	communityRoleTC            = "tc"
	communityRoleCommitter     = "committer"
	communityRoleRepoMember    = "repo_member"
	communityRoleSigMaintainer = "sig_maintainer"
)

var (
	CommunityRoleTC            = communityRole(communityRoleTC)
	CommunityRoleCommitter     = communityRole(communityRoleCommitter)
	CommunityRoleRepoMember    = communityRole(communityRoleRepoMember)
	CommunityRoleSigMaintainer = communityRole(communityRoleSigMaintainer)
)

type CommunityRole interface {
	CommunityRole() string
	IsTC() bool
}

func NewCommunityRole(v string) (CommunityRole, error) {
	switch strings.ToLower(v) {
	case communityRoleTC:
		return CommunityRoleTC, nil

	case communityRoleCommitter:
		return CommunityRoleCommitter, nil

	case communityRoleRepoMember:
		return CommunityRoleRepoMember, nil

	case communityRoleSigMaintainer:
		return CommunityRoleSigMaintainer, nil

	default:
		return nil, errors.New("unknown community role")
	}
}

type communityRole string

func (v communityRole) CommunityRole() string {
	return string(v)
}

func (v communityRole) IsTC() bool {
	return string(v) == communityRoleTC
}

func IsSameCommunityRole(r1, r2 CommunityRole) bool {
	return r1 != nil && r2 != nil && r1.CommunityRole() == r2.CommunityRole()
}
