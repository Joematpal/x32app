package main

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/grogersstephen/x32app/osc"
)

type homeScreen struct {
	title           *canvas.Text
	connectionInfo  *widget.Label
	lAddrEntry      line
	rAddrEntry      line
	faderResolution line
	connectB        *widget.Button
	channelBank     []*widget.Button
	dcaBank         []*widget.Button
	auxBank         []*widget.Button
	duration        line
	levelLabel      *widget.Label
	fadeTo          buttonLine
	fadeOutB        *widget.Button
	closeB          *widget.Button
	status          *widget.Label
	console         console
	selectedChannel int
}

func (h *homeScreen) setupDCABank() {
	dcaCount := 8
	h.dcaBank = make([]*widget.Button, dcaCount+1)
	for i := 1; i <= dcaCount; i++ {
		dca := i
		button := widget.NewButton(
			fmt.Sprintf("DCA%d", i),
			func() {
				var value float32
				var ok bool
				h.selectedChannel = 32 + dca
				// Make connection
				conn, err := connect()
				defer conn.Close()
				if err != nil {
					value = -1
				} else {
					osc.SendString(conn, fmt.Sprintf("/dca/%d/fader~~~~", dca))
					msg, err := osc.Listen(conn, 5*time.Second)
					if err != nil {
						value = -1
					}
					v := msg.DecodeArgument(0)
					value, ok = v.(float32)
					if !ok {
						value = -1
					}
				}
				h.levelLabel.SetText(
					fmt.Sprintf("DCA %d : %.2f\n", h.selectedChannel-32, value))
			})
		h.dcaBank[i] = button
	}
}
func (h *homeScreen) setupChannelBank() {
	channelCount := 32
	h.channelBank = make([]*widget.Button, channelCount+1)
	for i := 1; i <= channelCount; i++ {
		ch := i
		button := widget.NewButton(
			fmt.Sprintf("%02d", i),
			func() {
				var value float32
				var ok bool
				h.selectedChannel = ch
				// Make Connection
				conn, err := connect()
				defer conn.Close()
				if err != nil {
					value = -1
				} else {
					osc.SendString(conn, fmt.Sprintf("/ch/%02d/mix/fader~~~~", ch))
					msg, err := osc.Listen(conn, 5*time.Second)
					if err != nil {
						value = -1
					}
					//value, err = getChFader(conn, ch)
					v := msg.DecodeArgument(0)
					value, ok = v.(float32)
					if !ok {
						value = -1
					}
				}
				h.levelLabel.SetText(
					fmt.Sprintf("Channel %d : %.2f\n", h.selectedChannel, value))
			})
		h.channelBank[i] = button
	}
}

func (h *homeScreen) setup() {
	h.title = setupText("X32 App", color.White, 16)
	h.connectionInfo = setupLabel("")
	h.lAddrEntry = setupLine("local addr:", getLAddr(), "Set local address ip:port")
	h.rAddrEntry = setupLine("remote addr:", getRAddr(), "Set remote address ip:port")
	h.faderResolution = setupLine(
		"fader resolution:",
		fmt.Sprintf("%d", int(FADER_RESOLUTION)),
		"Set fader resolution")
	h.connectB = widget.NewButton("Connect", h.connectBPress)
	h.duration = setupLine("Duration: ", "2s", "")
	h.levelLabel = setupLabel("")
	h.fadeTo = setupButtonLine("\nFade To(0.00 to 1.00): \n", h.fadeToPress, "1", "")
	h.fadeOutB = widget.NewButton("\nFade Out\n", h.fadeOutPress)
	h.closeB = widget.NewButton("close", h.closeAppPress)
	h.status = setupLabel("")
	h.console.label = setupLabel("")
	h.console.scroller = container.NewVScroll(h.console.label)
	h.setupChannelBank()
	h.setupDCABank()

	win := App.NewWindow("Propres Ctrl")
	h.loadUI(win)
	win.Show()
}

func (h *homeScreen) loadUI(win fyne.Window) {
	h.console.log("")
	content :=
		container.New(layout.NewVBoxLayout(),
			h.title,
			container.NewGridWithColumns(1,
				container.NewGridWithColumns(2, h.lAddrEntry.label, h.lAddrEntry.entry),
				container.NewGridWithColumns(2, h.rAddrEntry.label, h.rAddrEntry.entry),
				container.NewGridWithColumns(2, h.faderResolution.label, h.faderResolution.entry),
			),
			container.NewGridWithColumns(1,
				h.connectB,
				h.status,
			),
			container.NewGridWithColumns(8,
				h.channelBank[1], h.channelBank[2], h.channelBank[3], h.channelBank[4],
				h.channelBank[5], h.channelBank[6], h.channelBank[7], h.channelBank[8],
				h.channelBank[9], h.channelBank[10], h.channelBank[11], h.channelBank[12],
				h.channelBank[13], h.channelBank[14], h.channelBank[15], h.channelBank[16],
			),
			container.NewGridWithColumns(4,
				h.dcaBank[1], h.dcaBank[2], h.dcaBank[3], h.dcaBank[4],
				h.dcaBank[5], h.dcaBank[6], h.dcaBank[7], h.dcaBank[8],
				//h.auxBank[1], h.auxBank[2], h.auxBank[3],
				//h.auxBank[4], h.auxBank[5], h.auxBank[6],
			),
			container.NewGridWithColumns(1,
				h.levelLabel,
				container.NewGridWithColumns(2,
					h.duration.label,
					h.duration.entry,
				),
			),
			container.NewGridWithColumns(2,
				h.fadeTo.button,
				h.fadeTo.entry,
			),
			h.fadeOutB,
			h.console.scroller,
			container.NewGridWithColumns(2,
				layout.NewSpacer(),
				h.closeB),
		)
	win.SetContent(content)
}
