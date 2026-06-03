package controllers

import (
	"regexp"
	"strconv"
	"strings"
)

type BaseController struct{}

var Base = BaseController{}

var subscriptionDatePattern = regexp.MustCompile(`^(0[1-9]|1[0-2])-\d{4}$`)

func (BaseController) isEndDateBeforeStartDate(startDate, endDate string) bool {
	startParts := strings.Split(startDate, "-")
	endParts := strings.Split(endDate, "-")
	startMonth, _ := strconv.Atoi(startParts[0])
	startYear, _ := strconv.Atoi(startParts[1])
	endMonth, _ := strconv.Atoi(endParts[0])
	endYear, _ := strconv.Atoi(endParts[1])
	return endYear*12+endMonth < startYear*12+startMonth
}
