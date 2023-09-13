package tui

import (
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type filePicker struct {
	flex *tview.Flex
	*tview.TreeView
}

func newFilePicker(pickFile func(path string)) *filePicker {
	rootDir := "."
	root := tview.NewTreeNode(rootDir).
		SetColor(tcell.ColorRed)
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)
	fp := &filePicker{TreeView: tree}

	add := func(target *tview.TreeNode, path string) {
		fileInfo, err := os.Stat(path)
		if err != nil {
			messager.setError(err.Error())
			pickFile("")
			return
		}
		if !fileInfo.IsDir() {
			pickFile(path)
			return
		}
		files, err := os.ReadDir(path)
		if err != nil {
			messager.setError(err.Error())
			pickFile("")
			return
		}
		for _, file := range files {
			node := tview.NewTreeNode(file.Name()).
				SetReference(filepath.Join(path, file.Name())).
				SetSelectable(true)
			if file.IsDir() {
				node.SetColor(tcell.ColorGreen)
			}
			target.AddChild(node)
		}
	}

	add(root, rootDir)

	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return
		}
		children := node.GetChildren()
		if len(children) == 0 {
			path := reference.(string)
			add(node, path)
		} else {
			node.SetExpanded(!node.IsExpanded())
		}
	})

	fp.flex = tview.NewFlex().
		AddItem(fp, 0, 1, true).
		AddItem(messager, 0, 1, false)
	return fp
}
