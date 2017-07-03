package main

import (
	"bytes"
	"fmt"

	"github.com/docker/docker/api/types"
	ui "github.com/gizak/termui"
)

type Dashboard struct {
	row              int
	client           DockerClient
	lenContainers    int
	containersWidget *ui.Table
	docWidget        *ui.Par
	detailWidget     *ui.Par
	containers       []types.Container
}

func (dashboard *Dashboard) refreshContainers() {
	rows := [][]string{
		[]string{"ID", "Image", "Ports", "Status", "Name"},
	}
	containers := dashboard.client.GetContainers()
	dashboard.lenContainers = len(containers)
	dashboard.containers = containers

	for _, container := range containers {
		var ports bytes.Buffer
		if len(container.Ports) > 0 {
			for index, port := range container.Ports {
				if index > 0 {
					ports.WriteString(", ")
				}

				ports.WriteString(fmt.Sprintf("%v", port.PublicPort))
				ports.WriteString("-->")
				ports.WriteString(fmt.Sprintf("%v", port.PrivatePort))
			}
		}
		var names bytes.Buffer
		if len(container.Names) > 0 {
			for index, name := range container.Names {
				if index > 0 {
					names.WriteString(", ")
				}
				names.WriteString(name)
				names.WriteString(", ")
			}
		}
		rows = append(rows, []string{container.ID[0:10], container.Image, ports.String(), container.Status, names.String()})
	}

	dashboard.containersWidget.Rows = rows
	dashboard.containersWidget.FgColors = nil
	dashboard.containersWidget.BgColors = nil
	dashboard.containersWidget.Analysis()

	if dashboard.row > dashboard.lenContainers {
		dashboard.row = dashboard.lenContainers
	}
	if dashboard.row > 0 {
		dashboard.containersWidget.BgColors[dashboard.row] = ui.ColorRed
	}
}

func (dashboard *Dashboard) drawContainers() {
	rows := [][]string{
		[]string{"ID", "Image", "Ports", "Status", "Name"},
	}

	table := ui.NewTable()
	table.Rows = rows
	table.FgColor = ui.ColorWhite
	table.BgColor = ui.ColorDefault
	table.TextAlign = ui.AlignCenter
	table.Separator = true
	table.Height = 28
	table.Border = true

	dashboard.containersWidget = table
}

func (dashboard *Dashboard) reload() {
	ui.Clear()
	ui.Render(ui.Body)
}

func (dashboard *Dashboard) setEvents() {
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	dashboard.containersWidget.Handle("/timer/1s", func(e ui.Event) {
		go func() {
			dashboard.refreshContainers()

			// detail
			if dashboard.row >= 0 && dashboard.row < dashboard.lenContainers {

				container := dashboard.containers[dashboard.row-1]
				stats := dashboard.client.ContainerStats(container)
				dashboard.detailWidget.Text = fmt.Sprintf(`
    Name: %s
    CPU: %d %s
    Memory: %d MiB
    Process: %d
          `, stats.Name, stats.CPUStats.CPUUsage.UsageInKernelmode, "%", stats.MemoryStats.Usage/1024, stats.PidsStats.Current)

			} else {
				dashboard.detailWidget.Text = ""
			}

			dashboard.reload()
		}()

	})

	dashboard.containersWidget.Handle("/sys/kbd/<up>", func(e ui.Event) {
		if dashboard.row > 1 {
			dashboard.row--
			dashboard.containersWidget.BgColors[dashboard.row] = ui.ColorRed
			dashboard.containersWidget.BgColors[dashboard.row+1] = ui.ColorBlack
		}
		dashboard.reload()

	})

	dashboard.containersWidget.Handle("/sys/kbd/<down>", func(e ui.Event) {

		if dashboard.row < dashboard.lenContainers {
			dashboard.row++
			dashboard.containersWidget.BgColors[dashboard.row] = ui.ColorRed
			dashboard.containersWidget.BgColors[dashboard.row-1] = ui.ColorBlack
		}
		dashboard.reload()
	})

	dashboard.containersWidget.Handle("/sys/kbd/C-d", func(e ui.Event) {
		if dashboard.row > 0 {
			dashboard.client.ContainerRemove(dashboard.containers[dashboard.row-1])
		}
	})

	dashboard.containersWidget.Handle("/sys/kbd/C-s", func(e ui.Event) {
		if dashboard.row > 0 {
			container := dashboard.containers[dashboard.row-1]
			if container.State != "running" {
				dashboard.client.ContainerStart(container)
			} else {
				dashboard.client.ContainerStop(container)
			}
		}
	})
}

func (dashboard *Dashboard) drawDoc() {
	doc := ui.NewPar(`
  <up>/<down>: Select container
  C-s: Stop/Start container
  C-d: Delete (force) container
  q: Exit
  `)
	doc.BorderLabel = "Usage"
	doc.BorderFg = ui.ColorYellow
	doc.Height = 10
	dashboard.docWidget = doc
}

func (dashboard *Dashboard) drawDetail() {
	detail := ui.NewPar(`

  `)
	detail.BorderLabel = "Detail"
	detail.BorderFg = ui.ColorYellow
	detail.Height = 10
	dashboard.detailWidget = detail
}

func InitDashboard() Dashboard {

	client := CreateClient()

	dashboard := Dashboard{
		row:           1,
		client:        client,
		lenContainers: 0,
	}

	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	dashboard.drawContainers()
	dashboard.drawDoc()
	dashboard.drawDetail()

	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(12, 0, dashboard.containersWidget)),
		ui.NewRow(
			ui.NewCol(6, 0, dashboard.detailWidget),
			ui.NewCol(6, 0, dashboard.docWidget),
		),
	)
	ui.Body.Align()
	dashboard.setEvents()
	ui.Render(ui.Body)

	ui.Loop()

	return dashboard
}
