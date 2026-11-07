package sidebar

import (
	"strings"
	"time"

	"github.com/killallgit/herdr-scope/internal/domain"
)

type Projection struct {
	Panes  map[domain.PaneID]PaneProjection
	Spaces map[domain.WorkspaceID]Tokens
	Next   time.Time
}

type PaneProjection struct {
	Agent  domain.AgentID
	Tokens Tokens
	// Restore keeps a carrier's title readable after a plain stop.
	Restore Tokens
}

func (p Presentation) Project(view domain.SidebarView) Projection {
	projected := Projection{
		Panes:  map[domain.PaneID]PaneProjection{},
		Spaces: map[domain.WorkspaceID]Tokens{},
		Next:   p.NextRepaint(view),
	}
	for _, group := range view.Groups {
		for i, row := range group.Rows {
			pane := PaneProjection{Agent: row.Agent, Tokens: p.PaneTokens(view, group, i)}
			if row.Fold != nil && row.Fold.Alone {
				pane.Restore = titles(p.FoldedTokens(view, row))
			}
			projected.Panes[row.Pane] = pane
		}
	}
	for _, row := range view.Folded {
		projected.Panes[row.Pane] = PaneProjection{Agent: row.Agent, Tokens: p.FoldedTokens(view, row)}
	}
	for _, space := range view.Spaces {
		projected.Spaces[space.Workspace] = p.SpaceTokens(space, view.At)
	}
	return projected
}

func titles(tokens Tokens) Tokens {
	out := Tokens{}
	for name, value := range tokens {
		if value != nil && strings.HasPrefix(name, titlePrefix) {
			out[name] = value
		}
	}
	return out
}
