package api

type ApiError struct {
	Message string
}

func (e *ApiError) Error() string {
	return e.Message
}

type FavoritesResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Count int `json:"count"`
		List  []struct {
			ID         int    `json:"id"`
			FId        int    `json:"fid"`
			Title      string `json:"title"`
			MediaCount int    `json:"media_count"`
		} `json:"list"`
	}
}

type CollectionsResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Count   int  `json:"count"`
		HasMore bool `json:"has_more"`
		List    []struct {
			ID         int    `json:"id"`
			FId        int    `json:"fid"`
			Title      string `json:"title"`
			MediaCount int    `json:"media_count"`
			Type       int    `json:"type"`
		} `json:"list"`
	}
}

type MediasResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		HasMore bool `json:"has_more"`
		Info    struct {
			Title string `json:"title"`
		} `json:"info"`
		Medias []struct {
			ID           int    `json:"id"`
			Title        string `json:"title"`
			Cover        string `json:"cover"`
			Introduction string `json:"intro"`
			Upper        struct {
				Name     string `json:"name"`
				MemberID int    `json:"mid"`
			} `json:"upper"`
			BVID string `json:"bvid"`
			Attr int    `json:"attr"`
		} `json:"medias"`
	}
}

type PageResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    []struct {
		Page int    `json:"page"`
		Part string `json:"part"`
	}
}

type Media struct {
	Identity     string
	UpperName    string
	UpperID      int
	Cover        string
	Introduction string
	Folder       string
	Title        string
	Active       bool
}
