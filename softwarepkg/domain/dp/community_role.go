package dp

const (
	communityRoleTC            = "tc"
	communityRoleCommitter     = "committer"
	communityRoleSigMaintainer = "sig_maintainer"
)

var (
	CheckItemTC            = communityRole(communityRoleTC)
	CheckItemCommitter     = communityRole(communityRoleCommitter)
	CheckItemSigMaintainer = communityRole(communityRoleSigMaintainer)
)

type CommunityRole interface {
	CommunityRole() string
	IsTC() bool
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
