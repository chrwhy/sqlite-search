package qparser

import (
	"log"
	"strings"
	"unicode"

	"github.com/chrwhy/open-pinyin/dict"
	pinyin "github.com/chrwhy/open-pinyin/parser"
	"github.com/chrwhy/open-pinyin/util"
)

const (
	SubPinyinStopSign = 3
)

func IsAllEn(query string) bool {
	for _, r := range []rune(query) {
		if !(r >= 'A' && r <= 'Z') && !(r >= 'a' && r <= 'z') {
			return false
		}
	}
	return true
}

func ParsePinyinClause(input string) string {
	pinyinGroups := pinyin.Parse(input)
	pinyinInitial := pinyin.ParseInitial(input)
	if len(pinyinInitial) > 0 {
		pinyinGroups = append(pinyinGroups, pinyinInitial)
	}
	clause := ""
	for i, pinyinGroup := range pinyinGroups {
		for j, _ := range pinyinGroup {
			if _, ok := dict.SUB_PINYIN[pinyinGroup[j]]; ok {
				if j != len(pinyinGroup)-1 && len(pinyinGroup[j]) > 1 {
					pinyinGroup[j] = "\"" + pinyinGroup[j] + string(rune(SubPinyinStopSign)) + "\""
				}
			}
		}
		clause += util.Concat(pinyinGroup, "+")
		if len(pinyinGroups) > 1 && i != len(pinyinGroups)-1 {
			clause += " OR "
		}
	}
	return clause
}

func ParseClause(query string) string {
	clause := ""
	spaceTokens := strings.Split(query, " ")

	regroupedTokens := make([]string, 0)
	for _, token := range spaceTokens {
		enCnTokens := splitCnEnToken(token)
		if enCnTokens == nil {
			regroupedTokens = append(regroupedTokens, token)
		} else {
			regroupedTokens = append(regroupedTokens, enCnTokens...)
		}
	}

	for _, token := range regroupedTokens {
		if IsAllEn(token) {
			log.Printf("Token: %s, Pinyin result: %v", token, pinyin.Parse(token))
			pinyinClause := ParsePinyinClause(token)
			partialSql := ""
			if len(clause) > 0 {
				partialSql = " AND "
			}
			if len(pinyinClause) > 0 {
				partialSql = partialSql + `(` + pinyinClause + " OR " + token + `)`
			} else {
				partialSql = partialSql + `("` + token + `")`
			}
			clause = clause + partialSql
			//log.Println(partialSql)
		} else {
			token = strings.Replace(token, "\"", "\"\"", -1)
			token = strings.Replace(token, "'", "''", -1)
			log.Printf("Token: %s", token)
			sql := `("` + token + `")`
			if len(clause) > 0 {
				clause = clause + " AND " + sql
			} else {
				clause = sql
			}
		}
	}

	return clause
}

func splitCnEnToken(input string) []string {
	var result []string
	var current string
	var currentType rune

	for _, r := range input {
		var charType rune
		if unicode.Is(unicode.Han, r) {
			charType = 'C' // Chinese
		} else if unicode.IsLetter(r) {
			charType = 'E' // English
		} else {
			charType = 'O' // Other
			return nil
		}

		if currentType == 0 {
			currentType = charType
		}

		if charType != currentType {
			if current != "" {
				result = append(result, current)
			}
			current = string(r)
			currentType = charType
		} else {
			current += string(r)
		}
	}

	if current != "" {
		result = append(result, current)
	}
	return result
}
