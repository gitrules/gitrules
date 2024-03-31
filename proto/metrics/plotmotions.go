package metrics

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"

	"github.com/gitrules/gitrules/lib/must"
	"github.com/gitrules/gitrules/runtime"
)

func plotDailyMotionsPNG(
	ctx context.Context,
	series *Series,

) []byte {

	// prepare py plotting program

	var w bytes.Buffer

	w.WriteString(`
import matplotlib.pyplot as plt
import numpy as np
from datetime import date

x = `)
	writePyDateArray(&w, series.DailyNumMotionOpen.X)
	fmt.Fprintln(&w)

	w.WriteString(`y1 = np.array(`)
	writePyIntArray(&w, series.DailyNumMotionOpen.Y)
	fmt.Fprintln(&w, ")")

	w.WriteString(`y2 = np.array(`)
	writePyIntArray(&w, series.DailyNumMotionClose.Y)
	fmt.Fprintln(&w, ")")

	w.WriteString(`y3 = np.array(`)
	writePyIntArray(&w, series.DailyNumMotionCancel.Y)
	fmt.Fprintln(&w, ")")

	w.WriteString(
		`fig, ax = plt.subplots(figsize=(9, 5))
ax.bar(x, y1, color='#55cc88')
ax.bar(x, y2, bottom=y1, color='#eeaa77')
ax.bar(x, y3, bottom=y1+y2, color='#cccccc')
`)

	w.WriteString(
		`ax.set_xlabel("Days")
ax.set_ylabel("Count")
ax.legend(["Opened", "Closed", "Cancelled", ])
ax.set_title("Daily issue and PR open/close/cancel counts")
`)

	n := series.DailyNumMotionOpen.Len()
	fmt.Fprintf(&w, "ax.set_xticks(x[0::%d])\n", xTickSkipDates(n))

	fp := filepath.Join(os.TempDir(), generateRandomID()+".png")
	fmt.Fprintf(&w, "fig.savefig(%q, dpi=200, bbox_inches = 'tight')\n", fp)

	py := w.String()
	fmt.Println(py)

	// call python
	outerr, err := runtime.RunPython(ctx, py)
	fmt.Println(string(outerr))
	must.NoError(ctx, err)

	// retrieve png plot
	pngData, err := os.ReadFile(fp)
	must.NoError(ctx, err)

	return pngData
}

func xTickSkipDates(n int) int {
	switch {
	case n > 31:
		k := n / 5
		if k < 1 {
			return 1
		}
		return int(math.Log2(float64(k)) + 3)
	case n >= 22:
		return 6
	case n >= 15:
		return 3
	case n >= 3:
		return 2
	default:
		return 1
	}
}
