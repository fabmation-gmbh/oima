package ui

// Some pieces of this Code where copied from the Original Widget "list"
// from https://github.com/gizak/termui/blob/master/v3/widgets/list.go
// The Code was modified by Emanuel Bennici <eb@fabmation.de>

import (
	"fmt"
	. "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	rw "github.com/mattn/go-runewidth"
	"image"

	"github.com/fabmation-gmbh/oima/internal"
	. "github.com/fabmation-gmbh/oima/internal/log"
	"github.com/fabmation-gmbh/oima/pkg/config"
	rt "github.com/fabmation-gmbh/oima/pkg/registry/interfaces"
)

var conf config.Configuration

type ImageInfo struct {
	Block

	Rows				*[]rt.Tag
	WrapText			bool
	TextStyle			Style
	SelectedRow			int
	topRow				int
	SelectedRowStyle	Style

	// Pointer to the List shown at the bottom of the UI, with Informations
	// about the selected Tag.
	// @ImageTagInfo can also be _nil_
	ImageTagInfo		*widgets.List

	// ImagePtr points back to the original Image,
	// this is needed because BaseImage contains
	// the Function to delete Signatures
	ImagePtr			*rt.BaseImage
}

func NewImageInfo() *ImageInfo {
	// init @conf
	conf = internal.GetConfig()

	return &ImageInfo{
		Block:            *NewBlock(),
		Rows:             nil,
		WrapText:         false,
		TextStyle:        Theme.List.Text,
		SelectedRowStyle: Theme.List.Text,
	}
}

func (ii *ImageInfo) Draw(buf *Buffer) {
	ii.Block.Draw(buf)

	point := ii.Inner.Min

	// adjusts view into widget
	if ii.SelectedRow >= ii.Inner.Dy()+ii.topRow {
		ii.topRow = ii.SelectedRow - ii.Inner.Dy() + 1
	} else if ii.SelectedRow < ii.topRow {
		ii.topRow = ii.SelectedRow
	}

	// draw rows
	for row := ii.topRow; row < len(*ii.Rows) && point.Y < ii.Inner.Max.Y; row++ {
		cells := ParseStyles(string((*ii.Rows)[row].Name), ii.TextStyle)
		if ii.WrapText {
			cells = WrapCells(cells, uint(ii.Inner.Dx()))
		}
		for j := 0; j < len(cells) && point.Y < ii.Inner.Max.Y; j++ {
			style := cells[j].Style
			if row == ii.SelectedRow {
				style = ii.SelectedRowStyle
			}
			if cells[j].Rune == '\n' {
				point = image.Pt(ii.Inner.Min.X, point.Y+1)
			} else {
				if point.X+1 == ii.Inner.Max.X+1 && len(cells) > ii.Inner.Dx() {
					buf.SetCell(NewCell(ELLIPSES, style), point.Add(image.Pt(-1, 0)))
					break
				} else {
					buf.SetCell(NewCell(cells[j].Rune, style), point)
					point = point.Add(image.Pt(rw.RuneWidth(cells[j].Rune), 0))
				}
			}
		}
		point = image.Pt(ii.Inner.Min.X, point.Y+1)
	}

	// draw UP_ARROW if needed
	if ii.topRow > 0 {
		buf.SetCell(
			NewCell(UP_ARROW, NewStyle(ColorWhite)),
			image.Pt(ii.Inner.Max.X-1, ii.Inner.Min.Y),
		)
	}

	// draw DOWN_ARROW if needed
	if len(*ii.Rows) > int(ii.topRow)+ii.Inner.Dy() {
		buf.SetCell(
			NewCell(DOWN_ARROW, NewStyle(ColorWhite)),
			image.Pt(ii.Inner.Max.X-1, ii.Inner.Max.Y-1),
		)
	}

	// update Image Tag Info
	ii.updateTagInfo()
}

// ScrollAmount scrolls by amount given. If amount is < 0, then scroll up.
// There is no need to set ii.topRow, as this will be set automatically when drawn,
// since if the selected item is off screen then the topRow variable will change accordingly.
func (ii *ImageInfo) ScrollAmount(amount int) {
	if len(*ii.Rows)-int(ii.SelectedRow) <= amount {
		ii.SelectedRow = len(*ii.Rows) - 1
	} else if int(ii.SelectedRow)+amount < 0 {
		ii.SelectedRow = 0
	} else {
		ii.SelectedRow += amount
	}
}

func (ii *ImageInfo) ScrollUp() {
	ii.ScrollAmount(-1)

	// update shown Tag Info
	ii.updateTagInfo()
}

func (ii *ImageInfo) ScrollDown() {
	ii.ScrollAmount(1)

	// update shown Tag Info
	ii.updateTagInfo()
}

func (ii *ImageInfo) ScrollPageUp() {
	// If an item is selected below top row, then go to the top row.
	if ii.SelectedRow > ii.topRow {
		ii.SelectedRow = ii.topRow
	} else {
		ii.ScrollAmount(-ii.Inner.Dy())
	}

	// update shown Tag Info
	ii.updateTagInfo()
}

func (ii *ImageInfo) ScrollPageDown() {
	ii.ScrollAmount(ii.Inner.Dy())

	// update shown Tag Info
	ii.updateTagInfo()
}

func (ii *ImageInfo) ScrollHalfPageUp() {
	ii.ScrollAmount(-int(FloorFloat64(float64(ii.Inner.Dy()) / 2)))

	// update shown Tag Info
	ii.updateTagInfo()
}

func (ii *ImageInfo) ScrollHalfPageDown() {
	ii.ScrollAmount(int(FloorFloat64(float64(ii.Inner.Dy()) / 2)))

	// update shown Tag Info
	ii.updateTagInfo()
}

func (ii *ImageInfo) ScrollTop() {
	ii.SelectedRow = 0

	// update shown Tag Info
	ii.updateTagInfo()
}

func (ii *ImageInfo) ScrollBottom() {
	ii.SelectedRow = len(*ii.Rows) - 1

	// update shown Tag Info
	ii.updateTagInfo()
}

func (ii *ImageInfo) DeleteSignature() {
	// delete Signature
	(*ii.ImagePtr).DeleteSignature(&(*ii.Rows)[ii.SelectedRow])

	// update Tag Informations
	(*ii.Rows)[ii.SelectedRow].S3SignFound = false
	ii.updateTagInfo()
}

/// >>>>> Internal Function <<<<<

// updateTagInfo Updates the Info Box at the Bottom with
//  new Informations about the (new) selected Tag
func (ii *ImageInfo) updateTagInfo() {
	if ii.ImageTagInfo != nil {
		// get S3 Signature Status
		var s3SignatureStatus string
		if !conf.S3.Enabled {
			s3SignatureStatus = "[S3 Component Disabled](fg:red)"
		} else {
			if (*ii.Rows)[ii.SelectedRow].S3SignFound {
				s3SignatureStatus = "[Signature found](fg:green)"
			} else {
				s3SignatureStatus = "[Signature not found](fg:red)"
			}
		}

		ii.ImageTagInfo.Rows = []string{
			"",
			fmt.Sprintf("[Tag Name:](mod:bold,fg:clear)              %s", (*ii.Rows)[ii.SelectedRow].Name),
			fmt.Sprintf("[Content Digest:](mod:bold,fg:clear)        %s", (*ii.Rows)[ii.SelectedRow].ContentDigest),
			fmt.Sprintf("[Signature found in S3:](mod:bold,fg:clear) %s", s3SignatureStatus),
			"[](fg:clear)",
		}
	} else {
		Log.Warningf("ImageInfo.ImageTagInfo is nil!")
	}
}