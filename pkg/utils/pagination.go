package utils

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Pagination struct {
	Page   int `json:"page,omitempty"`
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

func GetPagination(c *fiber.Ctx) *Pagination {
	pageParam := c.Query("page", "1")
	limitParam := c.Query("limit", "10")

	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	return &Pagination{
		Page:   page,
		Limit:  limit,
		Offset: offset,
	}
}

func ParseSort(raw string, allowedSortFields []string, fallback string) string {

	if fallback == "" {
		fallback = "id desc"
	}

	if raw == "" {
		return fallback
	}

	parts := strings.SplitN(raw, "-", 2)
	field := parts[0]
	direction := "asc"
	if len(parts) == 2 {
		if parts[1] == "desc" || parts[1] == "asc" {
			direction = parts[1]
		}
	}

	if !Includes(allowedSortFields, field) {
		return fallback
	}

	return field + " " + direction
}
