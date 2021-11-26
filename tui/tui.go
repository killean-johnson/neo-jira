package tui

import (
	"log"

	"github.com/jroimartin/gocui"
	"github.com/killean-johnson/jira-tui/api"
	"github.com/killean-johnson/jira-tui/config"
)

func CreateTUI(client *api.JiraClient, conf *config.Config) {
	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer gui.Close()

    gui.InputEsc = true

    var bl *BoardLayout = new(BoardLayout)
    bl.client = client
    bl.gui = gui
    bl.config = conf
    bl.keymap = make(map[string]func(*gocui.Gui, *gocui.View) error)
    bl.keymap["blcursordown"] = cursorDown
    bl.keymap["blcursorup"] = cursorUp
    bl.keymap["blselect"] = bl.switchToIssueLayout
    bl.keymap["blquit"] = boardQuit

    gui.SetManagerFunc(bl.Layout)

	if err := bl.boardLayoutKeybindings(); err != nil {
		log.Panicln(err)
	}

	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()

		// Check to make sure the next line actually exists
		if l, err := v.Line(cy + 1); err == nil && l != "" {
			// Set the cursor to the next line down
			if err := v.SetCursor(cx, cy+1); err != nil {
				// Change the origin if we've hit the bottom
				if err := v.SetOrigin(ox, oy+1); err != nil {
					return err
				}
			}
		} else {
			return err
		}
	}
	return nil
}
