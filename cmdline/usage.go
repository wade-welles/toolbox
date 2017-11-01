package cmdline

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kardianos/osext"
	"github.com/richardwilkes/gokit/i18n"
	"github.com/richardwilkes/gokit/xio/term"
)

var (
	// AppCmdName holds the application's name as specified on the command
	// line.
	AppCmdName string
	// AppName holds the name of the application. By default, this is the same
	// as AppCmdName.
	AppName string
	// CopyrightYears holds the years to place in the copyright banner.
	CopyrightYears string
	// CopyrightHolder holds the name of the copyright holder.
	CopyrightHolder string
	// License holds the license the software is being distributed under. This
	// is intended to be a simple one line description, such as "Mozilla
	// Public License 2.0" and not the full license itself.
	License string
	// AppVersion holds the application's version information. Typically set
	// by the build system.
	AppVersion string
	// GitVersion holds the git revision and clean/dirty status and should be
	// set by the build system.
	GitVersion string
)

func init() {
	if path, err := osext.Executable(); err == nil {
		path = filepath.Base(path)
		if path != "." {
			AppCmdName = path
		}
	}
	if AppCmdName == "" {
		AppCmdName = "<unknown>"
	}
	AppName = AppCmdName
}

// DisplayUsage displays the program usage information.
func (cl *CmdLine) DisplayUsage() {
	version := NewVersionFromString(AppVersion)
	term.WrapText(cl, "", AppName)
	term.WrapText(cl, "  ", version.Format(true, false))
	var dot string
	if !strings.HasSuffix(CopyrightHolder, ".") {
		dot = "."
	}
	term.WrapText(cl, "  ", fmt.Sprintf(i18n.Text("Copyright (c) %[1]s by %[2]s%s All rights reserved."), CopyrightYears, CopyrightHolder, dot))
	if License != "" {
		term.WrapText(cl, "  ", fmt.Sprintf(i18n.Text("License: %s"), License))
	}
	if !version.IsDevelopment() && !version.IsWhenUnset() {
		when := version.When().Local()
		term.WrapText(cl, "  ", fmt.Sprintf(i18n.Text("Built on %[1]s at %[2]s"), when.Format("Monday, January 2, 2006"), when.Format("3:04:05 PM MST")))
	}
	if GitVersion != "" {
		term.WrapText(cl, "  ", "git: "+GitVersion)
	}
	fmt.Fprintln(cl)
	if cl.Description != "" {
		term.WrapText(cl, "", cl.Description)
		fmt.Fprintln(cl)
	}
	usage := fmt.Sprintf(i18n.Text("%s [options]"), AppCmdName)
	opts := cl
	var stack []*CmdLine
	for opts != nil {
		stack = append(stack, opts)
		opts = opts.parent
	}
	for i := len(stack) - 1; i >= 0; i-- {
		one := stack[i]
		if one.cmd == nil {
			if i == 0 && len(cl.cmds) > 0 {
				usage += i18n.Text(" <command> [command options]")
			}
		} else {
			usage += fmt.Sprintf(i18n.Text(" %[1]s [%[1]s options]"), one.cmd.Name())
		}
	}
	if cl.UsageSuffix != "" {
		usage += " " + cl.UsageSuffix
	}
	term.WrapText(cl, i18n.Text("Usage: "), usage)
	for i := len(stack) - 1; i >= 0; i-- {
		one := stack[i]
		fmt.Fprintln(one)
		if one.cmd == nil {
			if i == 0 {
				usage += i18n.Text(" <command> [command options]")
			}
			fmt.Fprintln(one, i18n.Text("Options:"))
		} else {
			fmt.Fprintf(one, i18n.Text("%s options:\n"), one.cmd.Name())
		}
		fmt.Fprintln(one)
		one.displayOptions()
	}
	cl.displayCommands(2)
}

func (cl *CmdLine) displayOptions() {
	sort.Sort(cl.options)
	hasShort := false
	largest := 0
	for _, option := range cl.options {
		if option.usage != "" {
			if option.single != 0 {
				hasShort = true
			}
			length := len([]rune(option.name))
			if length > 0 {
				length += 2
			}
			if !option.isBool() {
				if length > 0 {
					length++
				}
				length += 2 + len([]rune(option.arg))
			}
			if length > largest {
				largest = length
			}
		}
	}
	largest += 2
	for _, option := range cl.options {
		if option.usage != "" {
			var sn string
			if hasShort {
				if option.single != 0 {
					sn = "-" + string(option.single)
					if option.name != "" {
						sn += ", "
					} else {
						sn += "  "
					}
				} else {
					sn = "    "
				}
			}
			var ln string
			if option.name != "" {
				ln = "--" + option.name
			}
			if !option.isBool() {
				if ln != "" {
					ln += " "
				}
				ln += "<" + option.arg + ">"
			}
			prefix := "  " + sn + ln + strings.Repeat(" ", largest-len([]rune(ln)))
			usage := option.usage
			if !strings.HasSuffix(usage, ".") {
				usage += "."
			}
			if !option.isBool() && option.def != "" {
				usage += i18n.Text(" Default: ")
				usage += option.def
			}
			term.WrapText(cl, prefix, usage)
		}
	}
}

func (cl *CmdLine) displayCommands(indent int) {
	if len(cl.cmds) > 0 {
		fmt.Fprintln(cl)
		term.WrapText(cl, "", i18n.Text("Available commands:"))
		fmt.Fprintln(cl)
		var all []string
		largest := 0
		for key := range cl.cmds {
			all = append(all, key)
			length := len(key)
			if length > largest {
				largest = length
			}
		}
		sort.Strings(all)
		format := fmt.Sprintf("%s%%-%ds  ", strings.Repeat(" ", indent), largest)
		for _, cmd := range all {
			term.WrapText(cl, fmt.Sprintf(format, cmd), cl.cmds[cmd].Usage())
		}
		fmt.Fprintln(cl)
		term.WrapText(cl, "", fmt.Sprintf(i18n.Text("Use '%s help <command>' to see command options"), AppCmdName))
	}
}