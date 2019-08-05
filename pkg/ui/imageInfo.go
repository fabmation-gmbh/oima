package ui

// Some pieces of this Code where copied from the Original Widget "list"
// from https://github.com/gizak/termui/blob/master/v3/widgets/list.go
// The Code was modified by Emanuel Bennici <eb@fabmation.de>

import (
	"github.com/awnumar/memguard"
	. "github.com/fabmation-gmbh/oima/internal/log"
	"github.com/fabmation-gmbh/oima/pkg/registry"

	ui "github.com/gizak/termui/v3"
	rw "github.com/mattn/go-runewidth"
	"image"
)

type ImageInfo struct {
	ui.Block

	// Registry Image
	Image				*registry.Image

	// Text Style
	TextStyle        	ui.Style

	// Current Row over which the Text Cursor "stands"
	SelectedRow      	int

	// Style of the @SelectedRow
	SelectedRowStyle 	ui.Style

	topRow           	int
}

// NewImageInfo returns a new empty *NewImageInfo
func NewImageInfo() *ImageInfo {
	return &ImageInfo{
		Block:            *ui.NewBlock(),
		Image:            nil,
		TextStyle:        ui.Theme.List.Text,
		SelectedRowStyle: ui.Theme.List.Text,
	}
}

// Original Source Code was edited/ changed!!!
// Source: https://github.com/gizak/termui/blob/master/v3/widgets/list.go#L33
func (ii *ImageInfo) Draw(buf *ui.Buffer) {
	// check if Tags are Empty
	if len(ii.Image.Tags) == 0 {
		Log.Warningf("[Internal Warning] Image Tags are empty! Trying to fetch Tags")
		err := ii.Image.FetchAllTags()
		if err != nil {
			Log.Fatalf("Error while Fetch Tags from Image %s: %s", ii.Image.Name, err.Error())
			memguard.SafeExit(1)
		}

		Log.Warningf("Fetched %d Images from %s", len(ii.Image.Tags), ii.Image.Name)
	}

	ii.Block.Draw(buf)

	point := ii.Inner.Min

	// adjusts view into widget
	if ii.SelectedRow >= (ii.Inner.Dy() + ii.topRow) {
		ii.topRow = ii.SelectedRow - ii.Inner.Dy() + 1
	} else if ii.SelectedRow < ii.topRow {
		ii.topRow = ii.SelectedRow
	}

	// draw rows
	for row := ii.topRow; row < len(ii.Image.Tags) && point.Y < ii.Inner.Max.Y; row++ {
		cells := ui.ParseStyles(string(ii.Image.Tags[row].TagName), ii.TextStyle)
		for j := 0; j < len(cells) && point.Y < ii.Inner.Max.Y; j++ {
			style := cells[j].Style
			if row == ii.SelectedRow {
				style = ii.SelectedRowStyle
			}
			if cells[j].Rune == '\n' {
				point = image.Pt(ii.Inner.Min.X, point.Y+1)
			} else {
				if (point.X + 1 == ii.Inner.Max.X + 1) && (len(cells) > ii.Inner.Dx()) {
					buf.SetCell(ui.NewCell(ui.ELLIPSES, style), point.Add(image.Pt(-1, 0)))
					break
				} else {
					buf.SetCell(ui.NewCell(cells[j].Rune, style), point)
					point = point.Add(image.Pt(rw.RuneWidth(cells[j].Rune), 0))
				}
			}
		}
		point = image.Pt(ii.Inner.Min.X, point.Y+1)
	}

	// draw UP_ARROW if needed
	if ii.topRow > 0 {
		buf.SetCell(
			ui.NewCell(ui.UP_ARROW, ui.NewStyle(ui.ColorWhite)),
			image.Pt(ii.Inner.Max.X - 1, ii.Inner.Min.Y),
		)
	}

	// draw DOWN_ARROW if needed
	if len(ii.Image.Tags) > (int(ii.topRow) + ii.Inner.Dy()) {
		buf.SetCell(
			ui.NewCell(ui.DOWN_ARROW, ui.NewStyle(ui.ColorWhite)),
			image.Pt(ii.Inner.Max.X - 1, ii.Inner.Max.Y - 1),
		)
	}
}

/// >>>>> Internal Functions <<<<<
