package board

type Feature struct {
	boardRepo       boardRepo
	boardMemberRepo boardMemberRepo
}

func New(
	boardRepo boardRepo,
	boardMemberRepo boardMemberRepo,
) *Feature {
	return &Feature{
		boardRepo:       boardRepo,
		boardMemberRepo: boardMemberRepo,
	}
}
