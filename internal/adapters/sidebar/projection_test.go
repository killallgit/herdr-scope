package sidebar

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/killallgit/herdr-scope/internal/domain"
	"github.com/killallgit/herdr-scope/internal/ports"
)

func TestProjectionSchedulesVisibleContent(t *testing.T) {
	for _, layout := range []ports.Layout{ports.LayoutClassic, ports.LayoutCompact} {
		for _, tc := range []struct {
			name    string
			fold    *domain.Fold
			folded  bool
			blocked bool
			want    time.Time
		}{
			{name: "hidden carrier age", fold: &domain.Fold{Label: "acme", Hidden: 2, Alone: true}},
			{name: "fold spinner", fold: &domain.Fold{Label: "acme", Hidden: 2, Alone: true, Working: true}, want: at0().Add(150 * time.Millisecond)},
			{name: "lone agent age", fold: &domain.Fold{Label: "acme"}, want: at0().Add(30 * time.Second)},
			{name: "hidden working row", folded: true},
			{name: "Spaces pulse", blocked: true, want: at0().Add(750 * time.Millisecond)},
		} {
			t.Run(string(layout)+"/"+tc.name, func(t *testing.T) {
				settings := ports.DefaultSettings()
				settings.Presentation.Layout = layout
				settings.Presentation.Variant = ports.VariantText
				settings.Presentation.ShowAge = ports.AgeAll
				p := NewPresentation(settings, nil)
				row := domain.SidebarRow{
					Pane: "w1:p1", Kind: "claude", Display: domain.DisplayIdle,
					Title: "Draft release notes", Since: at0().Add(-90 * time.Second), Fold: tc.fold,
				}
				view, _ := rowView(row, true)
				if tc.folded {
					row.Display = domain.DisplayWorking
					view.Groups, view.Folded = nil, []domain.SidebarRow{row}
				}
				if tc.blocked {
					view.Groups = nil
					view.Spaces = []domain.SpaceView{{Workspace: "w1", Mark: domain.DisplayBlocked, MarkKind: "claude"}}
				}

				projected := p.Project(view)

				assert.Equal(t, tc.want, projected.Next)
			})
		}
	}
}

func TestPresentationProjectsFoldTokensAndRestoration(t *testing.T) {
	p := textPresentation()
	view, _, _ := foldedView(domain.Fold{Root: "w1", Label: "acme", Hidden: 2, Working: true, WorkingKind: "claude"})
	view.Spaces = []domain.SpaceView{{Workspace: "w1", Label: "acme", Mark: domain.DisplayWorking, MarkKind: "claude"}}

	projected := p.Project(view)

	require.Len(t, projected.Panes, 3)
	assert.Equal(t, "▸ acme", val(t, projected.Panes["w1:p1"].Tokens, "fold"))
	assert.Empty(t, val(t, projected.Panes["w1:p1"].Tokens, "title_working"))
	assert.Equal(t, "1", val(t, projected.Panes["w3:p1"].Tokens, "folded"))
	assert.Equal(t, "○", val(t, projected.Spaces["w1"], "space_working_claude"))
	assert.Equal(t, view.At.Add(spinStep), projected.Next)
	assert.Equal(t, "⣷ Draft release notes", val(t, projected.Panes["w1:p1"].Restore, "title_working"))
	assert.Len(t, projected.Panes["w1:p1"].Restore, 1)
	assert.Empty(t, projected.Panes["w2:p1"].Restore)
	assert.Empty(t, projected.Panes["w3:p1"].Restore)
}
