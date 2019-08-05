package ui

import (
	"github.com/awnumar/memguard"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"

	. "github.com/fabmation-gmbh/oima/internal/log"
	"github.com/fabmation-gmbh/oima/pkg/registry"
)

var dockerReg registry.DockerRegistry
var tree []*TreeNode
var grid *ui.Grid
var repoImageTree *Tree
var stats *widgets.List
var imageInfo *widgets.List


type nodeValue string

func (nv nodeValue) String() string { return string(nv) }

func StartUI() {
	initRegistry()

	if err := ui.Init(); err != nil {
		Log.Fatalf("failed to initialize termui: %v", err);
		memguard.SafeExit(1)
	}
	defer ui.Close()

	// initialize Grid
	initGrid()

	// data
	tree = getTree()

	// widgets
	repoImageTree = NewTree()

	// draw UI
	drawFunction()

	previousKey := ""
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "j", "<Down>":
			repoImageTree.ScrollDown()
		case "k", "<Up>":
			repoImageTree.ScrollUp()
		case "<C-d>":
			repoImageTree.ScrollHalfPageDown()
		case "<C-u>":
			repoImageTree.ScrollHalfPageUp()
		case "<C-f>":
			repoImageTree.ScrollPageDown()
		case "<C-b>":
			repoImageTree.ScrollPageUp()
		case "g":
			if previousKey == "g" { repoImageTree.ScrollTop() }
		case "<Home>":
			repoImageTree.ScrollTop()
		case "<Enter>", "<Space>":
			repoImageTree.ToggleExpand()
		case "G", "<End>":
			repoImageTree.ScrollBottom()
		case "E":
			repoImageTree.ExpandAll()
		case "C":
			repoImageTree.CollapseAll()
		case "<Resize>":
			x, y := ui.TerminalDimensions()
			repoImageTree.SetRect(0, 0, x, y)
		case "i", "I":
			// check if selected Node is Image or Repository
			if repoImageTree.SelectedNode().isImage() {
				// update Image Info View
			}
		}

		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}

		// re-render UI
		ui.Render(grid)
	}
}

func drawFunction() {
	// TODO: add small Stats-Panel at the right side
	// TODO: move me to own Function
	stats = widgets.NewList()
	stats.Title = "Registry Stats"
	stats.TextStyle = ui.NewStyle(ui.ColorYellow)
	stats.WrapText = false
	stats.SetRect(0, 0, 25, 8)

	repoImageTree.TextStyle = ui.NewStyle(ui.ColorYellow)
	repoImageTree.WrapText = false
	repoImageTree.SetNodes(tree)

	imageInfo = widgets.NewList()
	imageInfo.Title = "Image Info"
	imageInfo.Rows = []string{""}
	imageInfo.TextStyle = ui.NewStyle(ui.ColorGreen)
	imageInfo.WrapText = false
	imageInfo.SetRect(0, 0, 25, 8)

	x, y := ui.TerminalDimensions()

	repoImageTree.SetRect(0, 0, x, y)

	// add Items to grid
	grid.Set(
		ui.NewRow(0.75,
			ui.NewCol(0.75, repoImageTree),	// Repo/ Image Tree	(75%)
			ui.NewCol(0.25, stats),			// Stats Panel		(25%)
		),
		ui.NewRow(0.25,
			ui.NewCol(1, imageInfo),			// Image Information
		),
	)

	// render and show the UI
	ui.Render(grid)
}

func getTree() []*TreeNode {
	// get List of Repositories
	repos := dockerReg.ListRepositories()
	Log.Debugf("Registry %s has %d Repositories", dockerReg.URI, len(repos))

	// create Tree
	var nodes []*TreeNode
	var _nodes []*TreeNode

	for _, v := range repos {
		images, err := v.ListImages()
		if err != nil {
			Log.Criticalf("Error while getting List of Images in Repo %s: %s", v.Name, err.Error())
			memguard.SafeExit(1)
		}

		var imageEntries []*TreeNode

		// fill 'treeEntry' with Images
		for _, img := range images {
			imgEntry := TreeNode{
				Value:    nodeValue(img.Name),
				Expanded: false,
				Image:    img,
				Nodes:    nil, // TODO: Implement Digest Information
			}

			imageEntries = append(imageEntries, &imgEntry)
		}

		// append Images to Repo Entry
		repoEntry := TreeNode{
			Value:    nodeValue(v.Name),
			Expanded: false,
			Nodes:    imageEntries,
		}

		// add Repo Entry to Registry Entry-Node
		_nodes = append(_nodes, &repoEntry)
	}

	nodes = []*TreeNode{
		{
			Value:    nodeValue(dockerReg.URI),
			Expanded: true,
			Nodes:    _nodes,
		},
	}

	return nodes
}

func initGrid() {
	// create grid
	grid = ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)
}

func getImageInfo(i *registry.Image) {

}

/// >>>>> internal Functions <<<<<

// Initialize Docker Registry and handle Errors
// TODO: better handle custom Errors!
func initRegistry() {
	err := dockerReg.Init()
	if err != nil {
		Log.Panicf("Error while Initialize DockerRegistry: %s", err.Error())
	}

	// Fetch All Informations from Docker Registry
	registryFetch()
}

// Fetch all Informations form the Docker Registry and handle Errors
// TODO: better handle custom Errors!
func registryFetch() {
	err := dockerReg.FetchAll()
	if err != nil {
		Log.Fatalf("Error while Fetching All Informations from Registry '%s': %s", dockerReg.URI, err.Error())
		memguard.SafeExit(1)
	}
}
