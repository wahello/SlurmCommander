package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pja237/slurmcommander-dev/internal/cmdline"
	"github.com/pja237/slurmcommander-dev/internal/command"
	"github.com/pja237/slurmcommander-dev/internal/config"
	"github.com/pja237/slurmcommander-dev/internal/logger"
	"github.com/pja237/slurmcommander-dev/internal/model"
	"github.com/pja237/slurmcommander-dev/internal/model/tabs/clustertab"
	"github.com/pja237/slurmcommander-dev/internal/model/tabs/jobdetailstab"
	"github.com/pja237/slurmcommander-dev/internal/model/tabs/jobfromtemplate"
	"github.com/pja237/slurmcommander-dev/internal/model/tabs/jobhisttab"
	"github.com/pja237/slurmcommander-dev/internal/model/tabs/jobtab"
	"github.com/pja237/slurmcommander-dev/internal/table"
	"github.com/pja237/slurmcommander-dev/internal/version"
)

func main() {

	var (
		debugSet bool = false
	)

	fmt.Printf("Welcome to Slurm Commander!\n\n")

	args, err := cmdline.NewCmdArgs()
	if err != nil {
		log.Fatalf("ERROR: parsing cmdline args: %s\n", err)
	}

	if *args.Version {
		version.DumpVersion()
		os.Exit(0)
	}

	cc := config.NewConfigContainer()
	err = cc.GetConfig()
	if err != nil {
		log.Printf("WARRNING: parsing config file: %s\n", err)
	}

	log.Printf("INFO: %s\n", cc.DumpConfig())
	command.NewCmdCC(*cc)
	jobtab.NewCmdCC(*cc)
	clustertab.NewCmdCC(*cc)
	jobhisttab.NewCmdCC(*cc)

	// TODO: move all this away to view somewhere...
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Background(lipgloss.Color("#0057B8")).
		Foreground(lipgloss.Color("#FFD700")).
		Bold(false)
	//s.Selected = s.Selected.
	//	Background(lipgloss.Color("#277BC0")).
	//	Foreground(lipgloss.Color("#FFCB42")).
	//	Bold(false)
	//
	//s.Selected = s.Selected.
	//	Foreground(lipgloss.Color("229")).
	//	Background(lipgloss.Color("57")).
	//	Bold(false)
	// Filter TextInput
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()
	ti.CharLimit = 20
	ti.Width = 20

	// logging
	debugSet, l := logger.SetupLogger()

	m := model.Model{
		Globals: model.Globals{
			Help:            help.New(),
			ActiveTab:       0,
			Log:             l,
			FilterSwitch:    -1,
			Debug:           debugSet,
			ConfigContainer: *cc,
			JobHistStart:    *args.HistDays,
			JobHistTimeout:  *args.HistTimeout,
		},
		JobTab: jobtab.JobTab{
			SqueueTable: table.New(table.WithColumns(jobtab.SqueueTabCols), table.WithRows(jobtab.TableRows{}), table.WithStyles(s)),
			Filter:      ti,
		},
		JobHistTab: jobhisttab.JobHistTab{
			SacctTable:    table.New(table.WithColumns(jobhisttab.SacctTabCols), table.WithRows(jobtab.TableRows{}), table.WithStyles(s)),
			Filter:        ti,
			HistFetched:   false,
			HistFetchFail: false,
		},
		JobDetailsTab: jobdetailstab.JobDetailsTab{
			SelJobIDNew: -1,
		},
		JobFromTemplateTab: jobfromtemplate.JobFromTemplateTab{
			EditTemplate: false,
			NewJobScript: "",
			TemplatesTable: table.New(
				table.WithColumns(jobfromtemplate.TemplatesListCols),
				table.WithRows(jobfromtemplate.TemplatesListRows{}),
				table.WithStyles(s),
			),
		},
		JobClusterTab: clustertab.JobClusterTab{
			SinfoTable: table.New(table.WithColumns(clustertab.SinfoTabCols), table.WithRows(jobtab.TableRows{}), table.WithStyles(s)),
			Filter:     ti,
		},
	}

	// OLD: bubbletea@v0.22.1
	//p := tea.NewProgram(tea.Model(m), tea.WithAltScreen())
	//if err := p.Start(); err != nil {
	//	log.Fatalf("ERROR: starting tea program: %q\n", err)
	//}

	// NEW: bubbletea@v0.23.0
	// Run returns the model as a tea.Model.
	p := tea.NewProgram(tea.Model(m), tea.WithAltScreen())
	ret, err := p.Run()
	if err != nil {
		fmt.Printf("Error starting program: %s", err)
		os.Exit(1)
	}

	if retMod, ok := ret.(model.Model); ok && retMod.Globals.SizeErr != "" {
		fmt.Printf("%s\n", retMod.Globals.SizeErr)
	}
	fmt.Printf("Goodbye!\n")
}
