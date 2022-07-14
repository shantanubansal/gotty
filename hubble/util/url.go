package util

import "fmt"

func BuildUrl(baseUrl string, queryParams map[string]string) string {
	rawQuery := ""
	appendAnd := false
	for k, v := range queryParams {
		if appendAnd {
			k = fmt.Sprintf("&%v", k)
		}
		rawQuery = fmt.Sprintf("%v%v=%v", rawQuery, k, v)
		appendAnd = true
	}
	if len(queryParams) > 0 {
		return fmt.Sprintf("%v?%v", baseUrl, rawQuery)
	}
	return baseUrl
}
