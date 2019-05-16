package main

import (
	"bufio"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func (s *myWindow) setHeader(level int) {
	var cfmt = gui.NewQTextCharFormat()
	switch level {
	case 1:
		cfmt.SetFontPointSize(18)
	case 2:
		cfmt.SetFontPointSize(16)
	case 3:
		cfmt.SetFontPointSize(14)
	default:
		cfmt.SetFontPointSize(12)
	}

	cfmt.SetForeground(gui.NewQBrush3(gui.NewQColor2(core.Qt__blue), core.Qt__SolidPattern))
	s.mergeFormatOnLineOrSelection(cfmt)
}

func (s *myWindow) textStyle(styleIndex int) {
	var cursor = s.editor.TextCursor()

	if styleIndex != 0 {

		var style = gui.QTextListFormat__ListDisc

		switch styleIndex {
		case 1:
			{
				style = gui.QTextListFormat__ListDisc
			}

		case 2:
			{
				style = gui.QTextListFormat__ListCircle
			}

		case 3:
			{
				style = gui.QTextListFormat__ListSquare
			}

		case 4:
			{
				style = gui.QTextListFormat__ListDecimal
			}

		case 5:
			{
				style = gui.QTextListFormat__ListLowerAlpha
			}

		case 6:
			{
				style = gui.QTextListFormat__ListUpperAlpha
			}

		case 7:
			{
				style = gui.QTextListFormat__ListLowerRoman
			}

		case 8:
			{
				style = gui.QTextListFormat__ListUpperRoman
			}
		}

		cursor.BeginEditBlock()

		var (
			blockFmt = cursor.BlockFormat()
			listFmt  = gui.NewQTextListFormat()
		)

		if cursor.CurrentList().Pointer() != nil {
			listFmt = gui.NewQTextListFormatFromPointer(cursor.CurrentList().Format().Pointer())
		} else {
			listFmt.SetIndent(blockFmt.Indent() + 1)
			blockFmt.SetIndent(0)
			cursor.SetBlockFormat(blockFmt)
		}

		listFmt.SetStyle(style)
		cursor.CreateList(listFmt)

		cursor.EndEditBlock()

	} else {
		var bfmt = gui.NewQTextBlockFormat()
		bfmt.SetObjectIndex(-1)
		cursor.MergeBlockFormat(bfmt)
	}
}

func (s *myWindow) textColor() {
	var col = widgets.QColorDialog_GetColor(s.editor.TextColor(), s.editor, "", 0)
	if !col.IsValid() {
		return
	}
	var cfmt = gui.NewQTextCharFormat()
	cfmt.SetForeground(gui.NewQBrush3(col, core.Qt__SolidPattern))
	s.mergeFormatOnLineOrSelection(cfmt)
}

func (s *myWindow) textBold() {
	var afmt = gui.NewQTextCharFormat()
	var fw = gui.QFont__Normal
	if s.actionTextBold.IsChecked() {
		fw = gui.QFont__Bold
	}
	afmt.SetFontWeight(int(fw))
	s.mergeFormatOnLineOrSelection(afmt)
}

func (s *myWindow) textUnderline() {
	var afmt = gui.NewQTextCharFormat()
	afmt.SetFontUnderline(s.actionTextUnderline.IsChecked())
	s.mergeFormatOnLineOrSelection(afmt)
}

func (s *myWindow) textStrikeOut() {
	var afmt = gui.NewQTextCharFormat()
	afmt.SetFontStrikeOut(s.actionStrikeOut.IsChecked())
	s.mergeFormatOnLineOrSelection(afmt)
}

func (s *myWindow) textItalic() {
	var afmt = gui.NewQTextCharFormat()
	afmt.SetFontItalic(s.actionTextItalic.IsChecked())
	s.mergeFormatOnLineOrSelection(afmt)
}

func (s *myWindow) insertImage() {
	filename := widgets.QFileDialog_GetOpenFileName(s.window, "select a file", ".", "Image (*.png *.jpg)", "Image (*.png *.jpg)", widgets.QFileDialog__ReadOnly)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	img := gui.NewQImage()
	ok := img.LoadFromData(data, len(data), "")
	if !ok {
		return
	}
	uri := core.NewQUrl3("rc://"+filename, core.QUrl__TolerantMode)

	s.editor.Document().AddResource(int(gui.QTextDocument__ImageResource), uri, img.ToVariant())
	url := uri.Url(core.QUrl__None)
	cursor := s.editor.TextCursor()
	cursor.InsertImage4(img, url)
	s.document.Images[url] = data
}

func (s *myWindow) getImageList(html string) []string {
	r := strings.NewReader(html)
	bufr := bufio.NewReader(r)
	regex, err := regexp.Compile(`<img src="([^"]+)" />`)
	if err != nil {
		//fmt.Println(err)
		return nil
	}
	res := []string{}
	for line, _, err := bufr.ReadLine(); err == nil; line, _, err = bufr.ReadLine() {
		line1 := string(line)
		res1 := regex.FindAllStringSubmatch(line1, -1)
		for i := 0; i < len(res1); i++ {
			res = append(res, res1[i][1])
		}
	}
	return res
}

func (s *myWindow) insertTable() {
	dlg := widgets.NewQDialog(s.window, core.Qt__Dialog)
	dlg.SetWindowTitle(T("Table Rows and Columns"))

	grid := widgets.NewQGridLayout(dlg)

	row := widgets.NewQLabel2(T("Rows:"), dlg, core.Qt__Widget)
	grid.AddWidget(row, 0, 0, 0)

	rowInput := widgets.NewQLineEdit(dlg)
	rowInput.SetText("3")
	rowInput.SetValidator(gui.NewQIntValidator(dlg))
	grid.AddWidget(rowInput, 0, 1, 0)

	col := widgets.NewQLabel2(T("Columns:"), dlg, core.Qt__Widget)

	grid.AddWidget(col, 1, 0, 0)

	colInput := widgets.NewQLineEdit(dlg)
	colInput.SetText("3")
	colInput.SetValidator(gui.NewQIntValidator(dlg))
	grid.AddWidget(colInput, 1, 1, 0)

	btb := widgets.NewQGridLayout(nil)

	okBtn := widgets.NewQPushButton2(T("OK"), dlg)
	btb.AddWidget(okBtn, 0, 0, 0)

	cancelBtn := widgets.NewQPushButton2(T("Cancel"), dlg)
	btb.AddWidget(cancelBtn, 0, 1, 0)

	grid.AddLayout2(btb, 2, 0, 1, 2, 0)

	dlg.SetLayout(grid)

	okBtn.ConnectClicked(func(b bool) {
		cursor := s.editor.TextCursor()
		r, err := strconv.Atoi(rowInput.Text())
		if err != nil {
			return
		}
		c, err := strconv.Atoi(colInput.Text())
		if err != nil {
			return
		}
		tbl := cursor.InsertTable2(r, c)
		tbl.Format().SetBorderBrush(gui.NewQBrush2(core.Qt__SolidPattern))
		dlg.Hide()
		dlg.Destroy(true, true)
	})

	cancelBtn.ConnectClicked(func(b bool) {
		dlg.Hide()
		dlg.Destroy(true, true)
	})

	dlg.SetModal(true)
	dlg.Show()
}