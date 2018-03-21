package cliapp

import (
	"errors"
	"fmt"
	"strings"
)

// result after arguments processd
type GetOptResult struct {
	NodeInfo *AppCmdNode
	Args     []string
	Opts     []*RetOptItem
}

func (self *GetOptResult) String() string {
	ret := ""

	for _, i := range self.Opts {
		ret += fmt.Sprintf("%s", i.Name)
		if i.HaveValue {
			ret += fmt.Sprintf("=%s", i.Value)
		}
		ret += "\n"
	}

	for _, i := range self.Args {
		ret += fmt.Sprintf("'%s'\n", i)
	}

	return ret
}

// If arguments contained many options with name 'name', then return only last
// of them
func (self *GetOptResult) GetLastNamedRetOptItem(name string) *RetOptItem {
	var ret *RetOptItem = nil
	for _, i := range self.Opts {
		if i.Name == name {
			ret = i
		}
	}
	return ret
}

// return all options with name 'name'
func (self *GetOptResult) GetAllNamedRetOptItems(name string) []RetOptItem {
	var ret []RetOptItem
	for _, i := range self.Opts {
		if i.Name == name {
			ret = append(ret, *i)
		}
	}
	return ret
}

// check result accordingly with AvailableOptions passed to
func (self *GetOptResult) CheckOptResult() []error {
	ret := make([]error, 0)
	for _, i := range self.NodeInfo.AvailableOptions {
		if i.IsRequired {
			if !self.DoesHaveNamedRetOptItem(i.Name) {
				ret = append(
					ret,
					errors.New("option "+i.Name+" is required but absent"),
				)
			}
		} else {
			if !self.DoesHaveNamedRetOptItem(i.Name) && i.HaveDefault {
				self.Opts = append(
					self.Opts,
					&RetOptItem{
						Name:      i.Name,
						HaveValue: true,
						Value:     i.Default,
					},
				)
			}
		}

		if i.MustHaveValue {
			item := self.GetLastNamedRetOptItem(i.Name)
			if item != nil && !item.HaveValue {
				if i.HaveDefault {
					item.HaveValue = true
					item.Value = i.Default
				} else {
					ret = append(
						ret,
						errors.New("option "+i.Name+" is required to have value"),
					)
				}
			}
		}
	}

	for _, i := range self.Opts {
		found := false
		for _, ii := range self.NodeInfo.AvailableOptions {
			if ii.Name == i.Name {
				found = true
				break
			}
		}
		if !found {
			ret = append(
				ret,
				errors.New("unsupported option passed: "+i.Name),
			)
		}
	}

	if self.NodeInfo.CheckArgs {
		len_args := len(self.Args)

		if len_args < self.NodeInfo.MinArgs {
			ret = append(ret, errors.New("given too few arguments"))
		}

		if self.NodeInfo.MaxArgs > -1 && len_args > self.NodeInfo.MaxArgs {
			ret = append(ret, errors.New("given too many arguments"))
		}
	}

	return ret
}

// check if option with name 'name' has been supplied
func (self *GetOptResult) DoesHaveNamedRetOptItem(name string) bool {
	for _, i := range self.Opts {
		if i.Name == name {
			return true
		}
	}
	return false
}

// check if option with name 'name' has been supplied and it have some value
func (self *GetOptResult) DoesNamedRetOptItemHaveValue(name string) bool {
	if self.DoesHaveNamedRetOptItem(name) {
		return self.GetLastNamedRetOptItem(name).HaveValue
	}
	return false
}

func (self *GetOptResult) GetOptionByName(name string) *GetOptCheckListItem {
	return self.NodeInfo.AvailableOptions.GetByName(name)
}

type RetOptItem struct {
	Name      string
	HaveValue bool
	Value     string
}

// NOTE: spaces, around option names, theyr values and arguments are not
//       stripped.
//       spaces considered to be stripped already by some other lexer
func GetOpt(args []string) *GetOptResult {

	var (
		last_opt_switch bool = false
		len_args        int
		i               int
		args_i_len      int
		ret             *GetOptResult
		eq_pos          int
	)

	ret = new(GetOptResult)

	len_args = len(args)

	i = 0

	for true {

		if i == len_args {
			break
		}

		if last_opt_switch {
			ret.Args = append(ret.Args, args[i])

		} else {

			args_i_len = len(args[i])

			if args_i_len == 0 {

				ret.Args = append(ret.Args, args[i])

			} else {

				if args[i] == "--" {
					last_opt_switch = true
				} else {

					if strings.HasPrefix(args[i], "-") {

						eq_pos = strings.Index(args[i], "=")
						_t := new(RetOptItem)

						if eq_pos != -1 {
							_t.Name = args[i][:eq_pos]
							_t.HaveValue = true
							_t.Value = args[i][eq_pos+1:]
						} else {
							_t.Name = args[i]
							_t.HaveValue = false
							_t.Value = ""
						}

						ret.Opts = append(ret.Opts, _t)

					} else {
						ret.Args = append(ret.Args, args[i])
					}
				}

			}

		}

		i++
	}

	return ret
}

type GetOptCheckListItem struct {
	/*
		must be with leading '-' (minuses) and without trailing '=' equals
	*/
	Name string

	/*
		set this to true if Default have meaningful value
	*/
	HaveDefault bool

	/*
		default value in case if defined without value
	*/
	Default string

	/*
		put here true, if user MUST set this flag
	*/
	IsRequired bool

	/* insist on value to flag.
	if this is true, and value not passed to option,
	then if HaveDefault is true then Default is used as value.
	*/
	MustHaveValue bool

	Description string // this text will be printed beside option on --help parameter
}

func (self *GetOptCheckListItem) HelpString() string {
	ret := ""

	if !self.IsRequired {
		ret += "["
	}
	ret += self.Name
	if self.MustHaveValue {
		ret += "="
		ret += self.Default
	}
	if !self.IsRequired {
		ret += "]"
	}

	return ret
}

type GetOptCheckList []*GetOptCheckListItem

func (self GetOptCheckList) GetByName(name string) *GetOptCheckListItem {
	for _, i := range self {
		if i.Name == name {
			return i
		}
	}
	return nil
}

func (self GetOptCheckList) HelpString() string {
	ret := ""
	self_len := len(self)
	for ii, i := range self {

		ret += i.HelpString()

		if ii != self_len-1 {
			ret += " "
		}
	}
	return ret
}
