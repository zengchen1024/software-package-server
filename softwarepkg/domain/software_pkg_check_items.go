package domain

import (
	"fmt"

	"github.com/opensourceways/software-package-server/softwarepkg/domain/dp"
)

func (entity *SoftwarePkg) CheckItemsMap() map[string]string {
	m := map[string]string{}

	for i := range commonCheckItems {
		item := &commonCheckItems[i]

		m[item.Id] = item.Name
	}

	others := entity.otherCheckItems()
	for i := range others {
		item := &others[i]

		m[item.Id] = item.Name
	}

	return m
}

func (entity *SoftwarePkg) CheckItems() []CheckItem {
	other := entity.otherCheckItems()

	r := make([]CheckItem, 0, len(other)+len(commonCheckItems))
	r = append(r, commonCheckItems...) // don't change commonCheckItems by copy it.
	r = append(r, other...)

	return r
}

func (entity *SoftwarePkg) otherCheckItems() []CheckItem {
	v := []CheckItem{
		{
			Id:            entity.Sig.ImportingPkgSig(),
			Name:          "Sig",
			Desc:          fmt.Sprintf("软件包被%s Sig接纳", entity.Sig.ImportingPkgSig()),
			Owner:         dp.CommunityRoleSigMaintainer,
			OnlyOwner:     true,
			Modifications: []string{pkgModificationSig},
		},
	}

	cs := entity.Repo.Committers
	for i := range cs {
		c := cs[i].Account.Account()

		v = append(v, CheckItem{
			Id:            c,
			Name:          "软件包维护人",
			Desc:          fmt.Sprintf("%s 同意作为此软件包的维护人", c),
			Owner:         dp.CommunityRoleCommitter,
			Keep:          true,
			OnlyOwner:     true,
			Modifications: []string{pkgModificationCommitter},
		})
	}

	return v
}
