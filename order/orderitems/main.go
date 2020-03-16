package order

import (
	"container/list"
	"errors"

	uuid "github.com/satori/go.uuid"
)

type (
	OrderingRule uint
)

const (
	OrderingRuleUndefined OrderingRule = iota
	OrderingRuleFirst
	OrderingRuleLast
	OrderingRuleBefore
	OrderingRuleBeforeLast
	OrderingRuleAfter
	OrderingRuleAfterFirst
)

var (
	ErrTooManyFirst             = errors.New("too many first rules")
	ErrTooManyLast              = errors.New("too many last rules")
	ErrRulesSetSanityCheckError = errors.New("rules set sanity check error")
)

type OrderingRulesItem struct {
	Rule     OrderingRule
	TargetId uuid.UUID
}

func OrderingRulesSanityCheck(
	ides []uuid.UUID,
	get_rules func(id uuid.UUID) *OrderingRulesItem,
) (log []error, first_rules, last_rules []uuid.UUID, err error) {

	// first_rules = []uuid.UUID{}
	// last_rules = []uuid.UUID{}

	for _, i := range ides {
		ide_rule := get_rules(i)

		switch ide_rule.Rule {
		case OrderingRuleFirst:
			first_rules = append(first_rules, i)
		case OrderingRuleLast:
			last_rules = append(last_rules, i)
		}

	}

	if len(first_rules) > 1 {
		log = append(log, ErrTooManyFirst)
	}

	if len(last_rules) > 1 {
		log = append(log, ErrTooManyLast)
	}

	return
}

func OrderItems(
	ides []uuid.UUID,
	get_rules func(id uuid.UUID) *OrderingRulesItem,
	middle_loop_iterations uint,
) (ret []uuid.UUID, err error) {

	if middle_loop_iterations == 0 {
		middle_loop_iterations = 2
	}

	lst := list.New()

	for _, i := range ides {
		lst.PushBack(i)
	}

	// first
first_loop:
	for e := lst.Front(); e != nil; e = e.Next() {
		e_rule := get_rules(e.Value.(uuid.UUID))

		if e_rule.Rule == OrderingRuleFirst {
			lst.MoveToFront(e)
			break first_loop
		}

	}

	// last
last_loop:
	for e := lst.Front(); e != nil; e = e.Next() {
		e_rule := get_rules(e.Value.(uuid.UUID))

		if e_rule.Rule == OrderingRuleLast {
			lst.MoveToBack(e)
			break last_loop
		}

	}

	for z := middle_loop_iterations; z != 0; z-- {

		for e := lst.Front(); e != nil; e = e.Next() {
			e_rule := get_rules(e.Value.(uuid.UUID))
			err = _MoveTargetElementByRules(lst, e, e_rule)
			if err != nil {
				return
			}
		}

	}

	ret = []uuid.UUID{}
	for e := lst.Front(); e != nil; e = e.Next() {
		ret = append(ret, e.Value.(uuid.UUID))
	}

	return
}

func _MoveTargetElementByRules(lst *list.List, input_e *list.Element, rule *OrderingRulesItem) (err error) {

	input_e_uuid := input_e.Value.(uuid.UUID)

	switch rule.Rule {
	case OrderingRuleFirst:
		lst.MoveToFront(input_e)
	case OrderingRuleLast:
		lst.MoveToBack(input_e)
	case OrderingRuleAfter:
		for e := lst.Front(); e != nil; e = e.Next() {
			if e.Value.(uuid.UUID) == rule.TargetId {
				lst.MoveAfter(input_e, e)
				break
			}
		}
	case OrderingRuleBefore:
		for e := lst.Front(); e != nil; e = e.Next() {
			if e.Value.(uuid.UUID) == rule.TargetId {
				lst.MoveBefore(input_e, e)
				break
			}
		}
	case OrderingRuleAfterFirst:
		lst.MoveAfter(input_e, lst.Front())
	case OrderingRuleBeforeLast:
		lst.MoveBefore(input_e, lst.Back())
	}

	return
}
