package main

import (
	"math/rand"

	"github.com/bennicholls/burl-E/burl"
	"github.com/veandco/go-sdl2/sdl"
)

type MainMenu struct {
	burl.BaseState

	menu *burl.List
}

func NewMainMenu() (mm *MainMenu) {
	mm = new(MainMenu)
	mm.menu = burl.NewList(10, 5, 10, 10, 1, true, "")
	mm.menu.CenterInConsole()

	mm.menu.Append("New Game", "Load Game", "Ship Designer", "Options", "Quit")
	return
}

func (mm *MainMenu) HandleKeypress(key sdl.Keycode) {
	switch key {
	case sdl.K_UP:
		mm.menu.Prev()
	case sdl.K_DOWN:
		mm.menu.Next()
	case sdl.K_RETURN:
		switch mm.menu.GetSelection() {
		case 0: //New Game
			burl.ChangeState(NewCreateGalaxyMenu())
		case 1: //Load Game
			//Load Game dialog
		case 2: //Ship Designer
			//Not sure if this one stays in.
		case 3: //Options
			//Options Dialog
		case 4: //Quit
			burl.PushEvent(burl.NewEvent(burl.QUIT_EVENT, ""))
		}
	}
}

func (mm *MainMenu) Render() {
	mm.menu.Render()
}

type CreateGalaxyMenu struct {
	burl.BaseState

	window         *burl.Container
	nameInput      *burl.Inputbox
	densityChoice  *burl.ChoiceBox //choice between some pre-defined densities
	shapeChoice    *burl.ChoiceBox //choice between blob, spiral, maybe more esoteric shapes??
	explainText    *burl.Textbox
	randomButton   *burl.Button
	generateButton *burl.Button
	cancelButton   *burl.Button

	focusedField burl.UIElem
	dialog       Dialog

	galaxyMap *burl.TileView

	galaxy *Galaxy
}

func NewCreateGalaxyMenu() (cgm *CreateGalaxyMenu) {
	cgm = new(CreateGalaxyMenu)

	cgm.window = burl.NewContainer(78, 43, 1, 1, 0, true)
	cgm.window.SetTitle("CREATE A WHOLE GALAXY WHY NOT")

	cgm.window.Add(burl.NewTextbox(5, 1, 2, 2, 1, false, false, "Name:"))
	cgm.nameInput = burl.NewInputbox(20, 1, 10, 2, 1, true)
	cgm.window.Add(burl.NewTextbox(5, 1, 2, 5, 1, false, false, "Density:"))
	cgm.densityChoice = burl.NewChoiceBox(20, 1, 10, 5, 2, true, burl.CHOICE_HORIZONTAL, "Sparse", "Normal", "Dense")
	cgm.window.Add(burl.NewTextbox(5, 1, 2, 8, 1, false, false, "Shape:"))
	cgm.shapeChoice = burl.NewChoiceBox(20, 1, 10, 8, 1, true, burl.CHOICE_HORIZONTAL, "Disk", "Spiral")

	cgm.explainText = burl.NewTextbox(30, 10, 2, 28, 1, true, false, "explanations")

	cgm.randomButton = burl.NewButton(15, 1, 58, 30, 2, true, true, "Randomize Galaxy")
	cgm.generateButton = burl.NewButton(15, 1, 58, 34, 1, true, true, "Generate the Galaxy as Shown!")
	cgm.generateButton.Register(burl.NewEvent(burl.BUTTON_PRESS, "generate"))
	cgm.cancelButton = burl.NewButton(15, 1, 58, 38, 2, true, true, "Return to Main Menu")
	cgm.cancelButton.Register(burl.NewEvent(burl.BUTTON_PRESS, "cancel"))

	cgm.galaxyMap = burl.NewTileView(25, 25, 53, 0, 0, true)

	cgm.window.Add(cgm.nameInput, cgm.densityChoice, cgm.shapeChoice, cgm.generateButton, cgm.explainText, cgm.cancelButton, cgm.galaxyMap, cgm.randomButton)

	cgm.nameInput.SetTabID(1)
	cgm.densityChoice.SetTabID(2)
	cgm.shapeChoice.SetTabID(3)
	cgm.randomButton.SetTabID(4)
	cgm.generateButton.SetTabID(5)
	cgm.cancelButton.SetTabID(6)

	cgm.focusedField = cgm.nameInput
	cgm.focusedField.ToggleFocus()
	cgm.UpdateExplanation()

	cgm.dialog = nil

	return
}

func (cgm *CreateGalaxyMenu) UpdateExplanation() {
	switch cgm.focusedField {
	case cgm.nameInput:
		cgm.explainText.ChangeText("GALAXY NAME:/n/nIt is believed that one of the main ways in which all sentient races of the galaxy are similar is a universal desire to name and label the universe. No Galaxy is complete without a name!")
	case cgm.densityChoice:
		cgm.explainText.ChangeText("GALAXY DENSITY:/n/nGalaxies come in all shapes, sizes and consistencies. Some are small and dense, with stars but a stone's throw away from each. Others have stars so spread out that many sentient species decide to never even attempt inter-system travel, instead deciding to focus efforts on art and philosophy and creating better and better tofu-based meat substitutes.")
	case cgm.shapeChoice:
		cgm.explainText.ChangeText("GALAXY SHAPE:/n/nGalaxies, like cookies, come in many different shapes. Some are globular, some are spirals, some are simple disks, and during certain times of year some are shaped like Christmas trees. (Note: currently only disk galaxies are created).")
	case cgm.randomButton:
		cgm.explainText.ChangeText("RANDOMIZE:/n/nIndecisive? Stunned by the marvelous array of choices before you? Let me do the work!")
	case cgm.generateButton:
		cgm.explainText.ChangeText("GENERATE:/n/n If this galaxy looks good, we can then generate the galaxy and move on to Ship Selection.")
	case cgm.cancelButton:
		cgm.explainText.ChangeText("CANCEL:/n/n Return to the main menu, discarding everything here.")
	}
}

func (cgm *CreateGalaxyMenu) Randomize() {
	names := []string{"The Biggest Galaxy", "The Galaxy of Terror", "The Lactose Blob", "The Thing Fulla Stars", "Andromeda 2", "Home"}

	cgm.nameInput.ChangeText(names[rand.Intn(len(names))])
	cgm.densityChoice.RandomizeChoice()
	cgm.shapeChoice.RandomizeChoice()
}

func (cgm *CreateGalaxyMenu) HandleKeypress(key sdl.Keycode) {
	if cgm.dialog != nil {
		cgm.dialog.HandleInput(key)
		return
	}

	switch key {
	case sdl.K_UP:
		cgm.focusedField.ToggleFocus()
		cgm.focusedField = cgm.window.FindPrevTab(cgm.focusedField)
		cgm.focusedField.ToggleFocus()
		cgm.UpdateExplanation()
	case sdl.K_DOWN, sdl.K_TAB:
		cgm.focusedField.ToggleFocus()
		cgm.focusedField = cgm.window.FindNextTab(cgm.focusedField)
		cgm.focusedField.ToggleFocus()
		cgm.UpdateExplanation()
	default:
		switch cgm.focusedField {
		case cgm.nameInput:
			switch key {
			case sdl.K_BACKSPACE:
				cgm.nameInput.Delete()
			case sdl.K_SPACE:
				cgm.nameInput.Insert(" ")
			default:
				cgm.nameInput.InsertText(rune(key))
			}
		case cgm.densityChoice:
			switch key {
			case sdl.K_LEFT:
				cgm.densityChoice.Prev()
			case sdl.K_RIGHT:
				cgm.densityChoice.Next()
			}
		case cgm.shapeChoice:
			switch key {
			case sdl.K_LEFT:
				cgm.shapeChoice.Prev()
			case sdl.K_RIGHT:
				cgm.shapeChoice.Next()
			}
		case cgm.randomButton:
			if key == sdl.K_RETURN {
				cgm.randomButton.Press()
				cgm.Randomize()
			}
		case cgm.generateButton:
			if key == sdl.K_RETURN {
				cgm.generateButton.Press()
			}
		case cgm.cancelButton:
			if key == sdl.K_RETURN {
				cgm.cancelButton.Press()
			}
		}
	}
}

func (cgm *CreateGalaxyMenu) HandleEvent(e *burl.Event) {
	switch e.ID {
	case burl.ANIMATION_DONE:
		if e.Message == "generate" {
			if cgm.nameInput.GetText() == "" {
				cgm.dialog = NewCommDialog("", "", "", "You must give your galaxy a name before you can continue!")
			}
		} else if e.Message == "cancel" {
			burl.ChangeState(NewMainMenu())
		}
	}
}

func (cgm *CreateGalaxyMenu) Update() {
	if cgm.dialog != nil && cgm.dialog.Done() {
		cgm.dialog.ToggleVisible()
		cgm.dialog = nil
	}
}

func (cgm *CreateGalaxyMenu) Render() {
	cgm.window.Render()

	if cgm.dialog != nil {
		cgm.dialog.Render()
	}
}
