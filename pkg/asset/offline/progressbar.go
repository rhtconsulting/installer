package offline

import (
	"log"
	"time"
	"fmt"
)

type Progbar struct {
        total int
}

const (
        maxbars  int           = 100
        interval time.Duration = 500 * time.Millisecond
        thebars  string        = "========================================================================================================"
)

func (p *Progbar) PrintComplete() {
        p.PrintProg(p.total)
        fmt.Print("\n")
}

func (p *Progbar) calcBars(portion int) int {
        if portion == 0 {
                return portion
        }

        return int(float32(maxbars) / (float32(p.total) / float32(portion)))
}

func check(err error) {
        if err != nil {
                log.Fatal(err)
        }
}

var src string

// func init() {
//         flag.Parse()
//         src = flag.Args()[0]
// }

func (p *Progbar) PrintProg(portion int) {
        bars := p.calcBars(portion)
        //spaces := maxbars - bars - 1
        percent := 100 * (float32(portion) / float32(p.total))

        fmt.Print("\033[G\033[K")
        fmt.Print("Progress [")

        fmt.Print(thebars[:bars])
        fmt.Print(">")
        fmt.Printf(" %3.2f%% (%d/%d) ]", percent, portion, p.total)
}
