package models

type (
	PaginateReq struct {
		Start string `query:"start" validate:"max=64"`
		Limit int    `query:"limit" validate:"min=2,max=20"`
	}

	PaginateRes struct {
		Data  any           `json:"data"`
		Limit int           `json:"limit"`
		Tolal int64         `json:"total"`
		First FirstPaginate `json:"first"`
		Next  NextPaginate  `json:"last"`
	}

	FirstPaginate struct {
		Herf  string `json:"herf"`
		Start string `json:"start"`
	}

	NextPaginate struct {
		Herf string `json:"herf"`
	}

	KafkaOffset struct {
		Offset int64 `json:"offset" bson:"offset"`
	}
)
