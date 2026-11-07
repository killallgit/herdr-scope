# Configuration and keyboard shortcuts

Scope’s settings and Herdr’s keyboard bindings go in **different configuration files**.

## Keyboard shortcuts

`--keys` belongs to the `configure` subcommand:

```sh
bin/herdr-scope configure --keys
```

This **prints a template only**; it does not install the shortcuts. Scope binds no keys itself.

Paste the following into **Herdr’s `config.toml`**, not Scope’s plugin configuration:

```toml
[[keys.command]]
key = "prefix+a"
type = "plugin_action"
command = "killallgit.herdr-scope.view-flip"
description = "Agents: active <-> recent"

[[keys.command]]
key = "prefix+u"
type = "plugin_action"
command = "killallgit.herdr-scope.jump"
description = "Agents: next one waiting on you"

[[keys.command]]
key = "prefix+f"
type = "plugin_action"
command = "killallgit.herdr-scope.view-here"
description = "Agents: only this workspace <-> all"

[[keys.command]]
key = "prefix+m"
type = "plugin_action"
command = "killallgit.herdr-scope.fold"
description = "Agents: fold <-> unfold this folder"

[[keys.command]]
key = "prefix+comma"
type = "plugin_action"
command = "killallgit.herdr-scope.settings"
description = "herdr-scope settings"
```

With Herdr’s default prefix, press **Ctrl+B, then the second key**:

| Second key | What happens |
| --- | --- |
| `a` | Switch between grouped activity order and a flat recent-activity list |
| `u` | Jump to an agent needing attention: questions first, then unseen results |
| `f` | Toggle this workspace versus all workspaces |
| `m` | Fold/unfold the current folder; agents needing attention remain visible |
| `,` | Open Scope settings |

The README describes the comma shortcut as `prefix+,`, but the generated configuration uses **`prefix+comma`**. Use the latter in the file.

This fork uses the plugin ID `killallgit.herdr-scope`. Existing upstream installations
use a different ID and keep separate configuration and state directories; these are
not migrated automatically.

## Where to change Scope’s settings

Scope’s configuration file is:

```sh
$(herdr plugin config-dir killallgit.herdr-scope)/config.toml
```

The easiest interface is the settings popup:

- **↑/↓:** select an option
- **←/→:** change it
- **Enter:** edit
- **r:** restore its default
- **s:** save and apply
- **q:** close

Manual file edits are picked up within about two seconds.

## Panel and organization

| Option | Default | Meaning |
| --- | --- | --- |
| `agents_panel` | `plugin` | Use Scope’s Agents panel or Herdr’s native one |
| `order` | `active` | `active`: grouped, busiest first; `recent`: flat activity list; `here`: current workspace and family; `off`: Herdr’s order |
| `nest_folders` | `true` | Nest workspaces whose folders are inside another workspace; also moves them together in Herdr’s list |
| `reorder_workspaces` | `false` | Make Herdr’s workspace numbering follow Scope’s activity order |
| `group_indent` | `2` | Indentation per tree level; `0` removes tree lines |
| `group_gap` | `true` | Blank line between groups; defaults to `false` with compact layout |
| `split_corner` | `false` | Visually branch additional split panes beneath the first |

**Code detail:** `agents_panel` and `order` are special live controls in the popup, not ordinary saved settings. The file’s `order`/`sort` supplies an initial ordering choice; `agents_panel` is not parsed as a normal file setting.

## Appearance

| Option | Default | Meaning |
| --- | --- | --- |
| `layout` | `classic` | Separate logo/status cells; `compact` makes denser rows |
| `variant` | `auto` | Icons: automatically choose, force `font`, use Unicode `text`, or `none` |
| `space_mark` | `logo` | Mark workspaces with an agent logo or a `state` symbol |
| `logo_color` | `active` | Vendor colors only while working/waiting, or `always` |
| `row_label` | `title` | Show the session `title`, `tab` name, or `both` |
| `trim_group_prefix` | `true` | Remove redundant workspace names from titles |
| `worktree_mark` | `U+F418` | Worktree-header symbol; empty disables it; requires a compatible font |
| `follow_appearance` | `false` | Follow desktop light/dark appearance; requires Herdr’s `[theme]` table |
| `colors.active_row_bg_light` | `#b9cdf2` | Selected-row background for light themes |
| `colors.active_row_bg_dark` | `#232636` | Selected-row background for dark themes |

Set either color override to `""` to retain the theme’s own background.

## Status and attention

| Option | Default | Meaning |
| --- | --- | --- |
| `done_hold` | `until_seen` | Keep ✓ until you focus the pane; alternatively specify seconds or `"30s"` |
| `blocked_hold` | `true` | Keep ? until the agent resumes working |
| `idle_grace_seconds` | `2.5` | Wait this long before treating idle as a completed turn |
| `activity_fresh_minutes` | `15` | How long an inactive agent looks “fresh” |
| `activity_stale_minutes` | `120` | When an inactive agent starts looking dim/stale |
| `show_age` | `waiting` | Show elapsed time only for attention-needed rows, for `all`, or `off` |
| `mark_tabs` | `true` | Append ? or ✓ to tabs needing attention |

## Advanced, file-only settings

```toml
log_level = "info"  # debug, info, warn, error
led_glyph = "●"     # Shared blocked/idle symbol

[observer]
heartbeat_interval = "1s"   # Agent-list refresh; 250ms–1m
reconcile_interval = "30s"  # Full-session refresh; 1s–1h

[glyphs]
done = "✔"                 # Override individual status symbols
```

Other glyph names: `working`, `blocked`, `idle_fresh`, `idle`, `idle_stale`, `unknown`.

## Sources

- [Shortcut generation and configure command](internal/bootstrap/cli.go)
- [Settings popup fields](internal/adapters/config/fields.go)
- [Configuration parsing and saving](internal/adapters/config/plugin.go)
- [Full configuration documentation](docs/configuration.md)
