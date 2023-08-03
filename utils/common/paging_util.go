package common

import (
	"log"
	"math"
	"os"
	"strconv"

	"github.com/jutionck/golang-db-sinar-harapan-makmur-orm/model/dto"
)

func GetPaginationParams(params dto.PaginationParam) dto.PaginationQuery {
	var page int
	var take int
	var skip int
	if params.Page > 0 {
		page = params.Page
	} else {
		page = 1
	}

	if params.Limit == 0 {
		err := LoadEnv()
		if err != nil {
			log.Println(err)
		}
		n, _ := strconv.Atoi(os.Getenv("DEFAULT_ROWS_PER_PAGE"))
		take = n
	} else {
		take = params.Limit
	}

	skip = (page - 1) * take

	return dto.PaginationQuery{
		Page: page,
		Take: take,
		Skip: skip,
	}
}

func Paginate(page, limit, totalRows int) dto.Paging {
	return dto.Paging{
		Page:        page,
		TotalPages:  int(math.Ceil(float64(totalRows) / float64(limit))),
		TotalRows:   totalRows,
		RowsPerPage: limit,
	}
}
