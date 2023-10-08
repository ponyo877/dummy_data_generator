package model

import (
	"time"

	"github.com/fatih/color"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

type Progress struct {
	*mpb.Progress
}

type Bar struct {
	*mpb.Bar
}

func NewProgress() (*Progress, func() /* wait func for mpb bug */) {
	return &Progress{
		mpb.New(
			mpb.WithOutput(color.Output),
			mpb.WithAutoRefresh(),
		),
	}, func() { time.Sleep(100 * time.Millisecond) }
}

func (p *Progress) AddBar(total int64, taskName string) Bar {
	var red, green = color.New(color.FgRed), color.New(color.FgGreen)
	pb := p.Progress.AddBar(total,
		mpb.PrependDecorators(
			decor.Name(taskName+": ", decor.WC{W: len(taskName) + 1, C: decor.DidentRight}),
			decor.CountersNoUnit("%d / %d ", decor.WCSyncWidth),
			decor.OnCompleteMeta(
				decor.OnComplete(
					decor.Meta(decor.Name("in progres", decor.WCSyncSpaceR), toMetaFunc(red)),
					"done!",
				),
				toMetaFunc(green),
			),
		),
		mpb.AppendDecorators(
			decor.OnComplete(decor.Percentage(decor.WC{W: 5}), ""),
		),
	)
	return Bar{pb}
}

func toMetaFunc(c *color.Color) func(string) string {
	return func(s string) string {
		return c.Sprint(s)
	}
}
