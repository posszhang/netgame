package util

import (
	"strconv"
	"strings"
)

type RuleItem struct {
	Id  uint32
	Num uint32
	Key uint32
}

func NewRuleItem(str string) *RuleItem {

	vec := strings.Split(str, ",")
	if len(vec) != 3 {
		return nil
	}

	id, _ := strconv.Atoi(vec[0])
	num, _ := strconv.Atoi(vec[1])
	key, _ := strconv.Atoi(vec[2])

	item := &RuleItem{}
	item.Id = uint32(id)
	item.Num = uint32(num)
	item.Key = uint32(key)

	return item
}

type Rule struct {
	rule []*RuleItem
}

func NewRule(str string) *Rule {
	rule := &Rule{
		rule: make([]*RuleItem, 0),
	}

	vec := strings.Split(str, ";")
	for _, s := range vec {
		if len(s) == 0 {
			continue
		}

		rule.Add(s)
	}

	return rule
}

func (rule *Rule) Add(s string) {
	rule.rule = append(rule.rule, NewRuleItem(s))
}

func (rule *Rule) Do() (uint32, uint32) {

	if len(rule.rule) == 0 {
		return 0, 0
	}

	sum := 0
	for _, item := range rule.rule {
		if item == nil {
			continue
		}

		sum += int(item.Key)
	}

	k := RandBetween(1, sum)
	cur := 0
	for _, item := range rule.rule {

		cur += int(item.Key)

		if cur >= k {
			return item.Id, item.Num
		}
	}

	return 0, 0
}

type RuleList struct {
	rulelist []*Rule
}

func NewRuleList(str string) *RuleList {

	rule := &RuleList{
		rulelist: make([]*Rule, 0),
	}

	vec := strings.Split(str, "|")
	for _, s := range vec {
		rule.Add(s)
	}
	return rule
}

func (rulelist *RuleList) Add(s string) {

	rule := NewRule(s)
	rulelist.rulelist = append(rulelist.rulelist, rule)
}

func (rulelist *RuleList) Do() map[uint32]uint32 {

	result := make(map[uint32]uint32)
	for _, rule := range rulelist.rulelist {
		if rule == nil {
			continue
		}

		id, num := rule.Do()
		if id == 0 {
			continue
		}
		result[id] += num
	}

	return result
}
