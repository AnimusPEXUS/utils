package cliapp

import (
	"fmt"
	"os"
	"strings"
)

type (
	AppResult struct {
		Code             int
		Message          string
		DoNotPrintResult bool
	}

	AdditionalInfo struct {
		OptsAndArgs      []string
		AvailableOptions GetOptCheckList
		DepthLevel       []string
		Subnode          *AppCmdNode
		Rootnode         *AppCmdNode
		Arg0             string
		PassData         interface{}
	}

	AppCmdNode struct {
		Name             string
		ShortDescription string
		Version          string
		VersionInfo      string
		DevStatus        string
		Description      string
		License          string
		Date             string
		Developers       []string
		URIs             []string
		ManPages         []string
		InfoPages        []string
		SubCmds          []*AppCmdNode

		Callable func(
			getopt_result *GetOptResult,
			adds *AdditionalInfo,
		) *AppResult

		AvailableOptions GetOptCheckList

		CheckArgs bool
		MinArgs   int // -1 == allow 0 arguments
		MaxArgs   int // -1 == no limit to arguments; 0 = no arguments allowed
	}
)

func CheckAppCmdNode(in *AppCmdNode, depth []string) {

	depth = append(depth, in.Name)

	if len(in.Name) == 0 {
		panic(
			fmt.Sprintf(
				"no AppCmdNode Name specified! on depth %d (at %s)",
				len(depth),
				strings.Join(depth, "/"),
			),
		)
	}

	if in.Callable != nil && len(in.SubCmds) != 0 {
		panic(
			fmt.Sprintf(
				"invalid structure of AppCmdNode on depth %d (at %s)"+
					" - both has children and is callable\n",
				len(depth),
				strings.Join(depth, "/"),
			),
		)
	}

	if in.Callable == nil && len(in.SubCmds) == 0 {
		panic(
			fmt.Sprintf(
				"invalid structure of AppCmdNode on depth %d (at %s)"+
					" - no children and no callable\n",
				len(depth),
				strings.Join(depth, "/"),
			),
		)
	}

	if in.MinArgs < -1 {
		panic(
			"MinArgs can't be lesser than -1",
		)
	}

	if in.MaxArgs < -1 {
		panic(
			"MaxArgs can't be lesser than -1",
		)
	}

	if in.MinArgs > in.MaxArgs && in.MaxArgs != -1 {
		panic(
			"MinArgs can't be larger than MaxArgs",
		)
	}

	for _, i := range in.SubCmds {
		CheckAppCmdNode(i, depth)
	}
}

// NOTE: this function does not returns: os.Exit is called before the return.
//       use RunCmd to avoid os.Exit
func RunApp(
	rootnode *AppCmdNode,
	pass_data interface{},
) {

	res := RunCmd(os.Args[0], os.Args[1:], rootnode, pass_data)

	if res == nil {
		panic("cmd returned nil. this considered to be development error")
	}

	PrintAppResult(res)

	os.Exit(res.Code)
}

func RunCmd(
	arg0 string,
	opts_and_args []string,
	rootnode *AppCmdNode,
	pass_data interface{},
) *AppResult {

	var (
		ret           *AppResult
		depth_level   []string
		subtree       *AppCmdNode
		getopt_result *GetOptResult
	)

	ret = &AppResult{0, "No Error", false}

	{
		CheckAppCmdNode(rootnode, make([]string, 0))

		getopt_result = GetOpt(opts_and_args)

		subtree = rootnode

		depth_level = append(depth_level[:0], depth_level[:0]...)

	look_for_requested_cmd_or_subtree:
		for _, ii := range getopt_result.Args {
			for _, jj := range subtree.SubCmds {
				if jj.Name == ii {
					subtree = jj
					depth_level = append(depth_level, ii)

					if subtree.Callable != nil {
						break look_for_requested_cmd_or_subtree
					}

					continue look_for_requested_cmd_or_subtree
				}
			}

			ret = &AppResult{
				100,
				"arguments error: invalid command or subsection name",
				false,
			}

		}

	}

	if ret.Code == 0 {

		var (
			show_help bool = false
		)

		getopt_result.Args = getopt_result.Args[len(depth_level):]

		{
			help_option := getopt_result.GetLastNamedRetOptItem("--help")
			if help_option != nil {
				show_help = true
				if help_option.HaveValue {
					ret = &AppResult{
						1,
						"options error: option --help does not assume additional value",
						false,
					}
				}
			}
		}

		if show_help {

			PrintHelp(depth_level, rootnode, subtree)
			ret = &AppResult{0, "help text printed", false}

		} else {

			if subtree.Callable != nil {
				getopt_result.NodeInfo = subtree
				{
					errs := getopt_result.CheckOptResult()
					if len(errs) != 0 {
						for _, i := range errs {
							fmt.Println("error:", i.Error())
						}
						return &AppResult{
							1,
							"options and/or arguments errors detected",
							false,
						}
					}
				}

				ret = subtree.Callable(
					getopt_result,
					&AdditionalInfo{
						OptsAndArgs:      opts_and_args,
						AvailableOptions: subtree.AvailableOptions,
						DepthLevel:       depth_level,
						Subnode:          subtree,
						Rootnode:         rootnode,
						Arg0:             arg0,
						PassData:         pass_data,
					},
				)

				if ret == nil {
					ret = &AppResult{Code: 0}
				}

				if ret.Code == 0 && ret.Message == "" {
					ret.Message = "No Error"
				}

			} else {
				ret = &AppResult{
					100,
					"arguments error: invalid command or subsection name",
					false,
				}
			}
		}

	}

	return ret
}

func PrintAppResult(in *AppResult) int {
	var ret int = 0
	if !in.DoNotPrintResult {
		fmt.Printf("Exit Code: %d (%s)\n", in.Code, in.Message)
	}
	return ret
}

func PrintHelp(
	depth_level []string,
	root *AppCmdNode,
	render_for *AppCmdNode,
) {
	var text string

	text += root.Name

	if len(root.Version) != 0 {
		text += " v" + root.Version
		if len(root.VersionInfo) != 0 {
			text += " (ver.info: " + root.VersionInfo + ")"
		}
	}

	if len(root.DevStatus) != 0 {
		text += " status: " + root.DevStatus
	}

	if len(root.License) != 0 {
		text += " license: " + root.License
	}

	if len(root.Date) != 0 {
		text += " date: " + root.Date
	}

	text += "\n\n"

	text += "Usage: "

	text += root.Name

	if len(depth_level) != 0 {
		text += " " + strings.Join(depth_level, " ")
	}

	ao_h := render_for.AvailableOptions.HelpString()
	text += " " + ao_h
	if len(ao_h) != 0 {
		text += " "
	}
	text += "[--] [args]"

	text += "\n\n"

	if len(render_for.ShortDescription) != 0 {
		new_deskr := render_for.ShortDescription
		new_deskr = strings.Replace(new_deskr, "\n", "\n  ", -1)
		text += "Short Info:\n  " + new_deskr
		text += "\n\n"
	}

	subsections := make([]*AppCmdNode, 0)
	subcommands := make([]*AppCmdNode, 0)

	for _, i := range render_for.SubCmds {
		if i.Callable == nil {
			subsections = append(subsections, i)
		} else {
			subcommands = append(subcommands, i)
		}
	}

	if len(subsections) != 0 {
		text += "SubSections:\n\n"
		for _, i := range subsections {
			text += " " + i.Name + "\n"
			if len(i.ShortDescription) != 0 {
				text += "  " + i.ShortDescription + "\n"
			}
			text += "\n"
		}
	}

	if len(subcommands) != 0 {
		text += "Commands:\n\n"
		for _, i := range subcommands {
			text += " " + i.Name + "\n"
			new_deskr := i.ShortDescription
			new_deskr = strings.Replace(new_deskr, "\n", "\n  ", -1)
			if len(i.ShortDescription) != 0 {
				text += "  " + new_deskr + "\n"
			}
			text += "\n"
		}
	}

	if render_for.Callable != nil {
		text += "Arguments: "
		if !render_for.CheckArgs {
			text += "unchecked"
		} else {

			texts := make([]string, 0)

			if render_for.MinArgs == -1 {
				texts = append(texts, "no minimum limit")
			} else {
				texts = append(texts,
					fmt.Sprintf("minimum count: %d", render_for.MinArgs))
			}

			if render_for.MaxArgs == -1 {
				texts = append(texts, "no maximum limit")
			} else {
				texts = append(texts,
					fmt.Sprintf("maximum count: %d", render_for.MaxArgs))
			}

			if (render_for.MaxArgs == render_for.MinArgs) && render_for.MinArgs != -1 {
				texts = append(
					texts,
					fmt.Sprintf("exact count: %d", render_for.MinArgs),
				)
			}

			text += strings.Join(texts, ", ") + "."

		}

		text += "\n\n"
	}

	opts := append(
		render_for.AvailableOptions,
		&GetOptCheckListItem{"--help", false, "", false, false, "Print help text"},
	)

	if len(opts) != 0 {
		text += "Options:\n\n"
		for _, i := range opts {

			text += " "

			if !i.IsRequired {
				text += "["
			}

			text += i.Name

			if i.MustHaveValue {
				text += "=value"
			}

			if !i.IsRequired {
				text += "]"
			}

			text += "\n"

			if len(i.Description) != 0 {
				text += "  " + i.Description + "\n"
			}

			text += "\n"
		}
	}

	if len(render_for.Description) != 0 {
		text += "Description:\n\n" + render_for.Description
		text += "\n\n"
	}

	if len(render_for.ManPages) != 0 {
		text += "Man Page"

		if len(render_for.ManPages) > 1 {
			text += "s"
		}
		text += ":\n"

		text += " " + strings.Join(render_for.ManPages, ", ") + "\n"

		text += "\n"
	}

	if len(render_for.InfoPages) != 0 {
		text += "Info Page"

		if len(render_for.InfoPages) > 1 {
			text += "s"
		}
		text += ":\n"

		text += " " + strings.Join(render_for.InfoPages, ", ") + "\n"

		text += "\n"
	}

	if len(render_for.Developers) != 0 {
		text += "Developer"

		if len(render_for.Developers) > 1 {
			text += "s"
		}
		text += ":\n"

		for _, i := range render_for.Developers {
			text += " " + i + "\n"
		}
		text += "\n"
	}

	if len(render_for.URIs) != 0 {
		text += "URI"

		if len(render_for.URIs) > 1 {
			text += "s"
		}
		text += ":\n"

		for _, i := range render_for.URIs {
			text += " " + i + "\n"
		}
		text += "\n"
	}

	fmt.Printf(text)
}
