package groupdomain

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

func (g *Group) ID() GroupID {
	return g.id
}

func (g *Group) GroupName() GroupName {
	return g.groupName
}
