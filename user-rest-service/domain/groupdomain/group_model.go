package groupdomain

import "golang.org/x/xerrors"

type Group struct {
	id        GroupID
	groupName GroupName
}

func NewGroup(id GroupID, groupName GroupName) *Group {
	return &Group{
		id:        id,
		groupName: groupName,
	}
}

func NewGroupWithoutID(groupName GroupName) *Group {
	return &Group{
		groupName: groupName,
	}
}

func (g *Group) ID() (GroupID, error) {
	if g.id == 0 {
		return 0, xerrors.Errorf("group id value is 0")
	}

	return g.id, nil
}

func (g *Group) GroupName() GroupName {
	return g.groupName
}

func (g *Group) UpdateGroupName(groupName GroupName) {
	g.groupName = groupName
}
