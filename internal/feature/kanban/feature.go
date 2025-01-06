package kanban

type Feature struct {
	boardRepo       boardRepo
	boardMemberRepo boardMemberRepo
}

func NewFeature(
	boardRepo boardRepo,
	boardMemberRepo boardMemberRepo,
) Feature {
	return Feature{
		boardRepo:       boardRepo,
		boardMemberRepo: boardMemberRepo,
	}
}
