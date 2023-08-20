package common

import (
	"fmt"
	"strconv"

	"github.com/fajritsaniy/golang-SHM/model/dto"
	"github.com/gin-gonic/gin"
)

func ValidateRequestQueryParams(c *gin.Context) (dto.RequestQueryParams, error) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		return dto.RequestQueryParams{}, fmt.Errorf("Invalid page number")
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if err != nil || limit <= 0 {
		return dto.RequestQueryParams{}, fmt.Errorf("Invalid limit value")
	}

	order := c.DefaultQuery("order", "id")
	sort := c.DefaultQuery("sort", "ASC")

	return dto.RequestQueryParams{
		QueryParams: dto.QueryParams{
			Order: order,
			Sort:  sort,
		},
		PaginationParam: dto.PaginationParam{
			Page:  page,
			Limit: limit,
		},
	}, nil
}
