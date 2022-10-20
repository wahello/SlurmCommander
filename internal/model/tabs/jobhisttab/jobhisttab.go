package jobhisttab

import (
	"log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/pja237/slurmcommander/internal/keybindings"
	"github.com/pja237/slurmcommander/internal/slurm"
)

type JobHistTab struct {
	SacctTable        table.Model
	SacctHist         slurm.SacctJobHist
	SacctHistFiltered slurm.SacctJobHist
	Filter            textinput.Model
	Stats
}

type Stats struct {
	StateCnt map[string]uint
}

func (t *JobHistTab) GetStatsFiltered(l *log.Logger) {
	t.Stats.StateCnt = map[string]uint{}

	l.Printf("GetStatsFiltered start\n")
	for _, v := range t.SacctHistFiltered.Jobs {
		t.Stats.StateCnt[*v.State.Current]++
	}
	l.Printf("GetStatsFiltered end\n")
}

type Keys map[*key.Binding]bool

var KeyMap = Keys{
	&keybindings.DefaultKeyMap.Up:       true,
	&keybindings.DefaultKeyMap.Down:     true,
	&keybindings.DefaultKeyMap.PageUp:   true,
	&keybindings.DefaultKeyMap.PageDown: true,
	&keybindings.DefaultKeyMap.Slash:    true,
	&keybindings.DefaultKeyMap.Info:     false,
	&keybindings.DefaultKeyMap.Enter:    true,
}

func (k *Keys) SetupKeys() {
	for k, v := range KeyMap {
		k.SetEnabled(v)
	}
}
