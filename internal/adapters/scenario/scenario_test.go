package scenario_test

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/killallgit/herdr-scope/internal/adapters/scenario"
	"github.com/killallgit/herdr-scope/internal/domain"
	"github.com/killallgit/herdr-scope/internal/ports"
)

func TestReplayKeepsFocusWhenAStepOmitsIt(t *testing.T) {
	s, err := scenario.Load(filepath.Join("..", "..", "..", "tests", "fixtures", "scenarios", "scenario-10.json"))
	require.NoError(t, err)
	s.Steps = s.Steps[:3]
	s.Steps[2].Focused = ""
	for _, agent := range s.Steps[2].Agents {
		agent["focused"] = json.RawMessage("false")
	}

	result, err := scenario.Run(s)

	require.NoError(t, err)
	assert.Empty(t, result.View.Folded, "the visit stays open at the last known focus")
	assert.Contains(t, result.Panes["w1:p1"]["fold"], "▾ acme")
}

func TestReplayFollowsFoldVisits(t *testing.T) {
	for _, layout := range []ports.Layout{ports.LayoutClassic, ports.LayoutCompact} {
		for _, tc := range []struct {
			name   string
			steps  int
			open   bool
			hidden int
		}{
			{name: "folded", steps: 1, hidden: 2},
			{name: "carrier focused", steps: 2, open: true},
			{name: "sibling focused", steps: 3, open: true},
			{name: "focus left", steps: 4, hidden: 2},
		} {
			t.Run(string(layout)+"/"+tc.name, func(t *testing.T) {
				s, err := scenario.Load(filepath.Join("..", "..", "..", "tests", "fixtures", "scenarios", "scenario-10.json"))
				require.NoError(t, err)
				s.Layout, s.Steps = layout, s.Steps[:tc.steps]

				result, err := scenario.Run(s)

				require.NoError(t, err)
				var folds []*domain.Fold
				for _, group := range result.View.Groups {
					for _, row := range group.Rows {
						if row.Fold != nil {
							folds = append(folds, row.Fold)
						}
					}
				}
				require.Len(t, folds, 2)
				assert.Equal(t, domain.WorkspaceID("w1"), folds[0].Root)
				assert.Equal(t, tc.open, folds[0].Open)
				assert.Equal(t, tc.hidden, folds[0].Hidden)
				assert.Equal(t, !tc.open, folds[0].Alone)
				assert.Equal(t, domain.WorkspaceID("w3"), folds[1].Root)
				assert.False(t, folds[1].Alone, "a lone agent keeps its own row")
				assert.NotEmpty(t, result.Panes["w3:p1"]["title_idle_fresh"])
				if tc.open {
					assert.Empty(t, result.View.Folded)
					assert.NotContains(t, result.Panes["w2:p1"], "folded")
					assert.Contains(t, result.Panes["w1:p1"]["fold"], "▾ acme")
				} else {
					require.Len(t, result.View.Folded, 1)
					assert.Equal(t, domain.PaneID("w2:p1"), result.View.Folded[0].Pane)
					assert.Equal(t, "1", result.Panes["w2:p1"]["folded"])
					assert.NotContains(t, result.Panes["w1:p1"], "title_idle_fresh")
				}
			})
		}
	}
}
