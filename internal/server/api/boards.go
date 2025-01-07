package api

type BoardShortInfo struct {
	ID        int    `json:"id"`
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

type AvailableBoardsResponse struct {
	Boards []BoardShortInfo `json:"boards"`
}

type GetBoardBySlugRequest struct {
	Slug string `uri:"slug" binding:"required"`
}

type GetBoardBySlugResponse struct {
}
