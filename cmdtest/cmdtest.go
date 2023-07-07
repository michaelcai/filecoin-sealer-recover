package cmdtest

import (
	// "context"

	"io"
	"os"

	"github.com/filecoin-project/go-commp-utils/writer"

	// cliutil "github.com/filecoin-project/lotus/cli/util"
	logging "github.com/ipfs/go-log/v2"

	"github.com/urfave/cli/v2"
)

// const (
// 	ss32GiB = 32 << 30
// )

var log = logging.Logger("recover")

var CmdtestCmd = &cli.Command{
	Name:      "cmdtest",
	Usage:     "Recovery sector tools",
	ArgsUsage: "[sectorNum1 sectorNum2 ...]",
	Flags:     []cli.Flag{},
	Action: func(cctx *cli.Context) error {

		// ctx := cliutil.DaemonContext(cctx)
		// ctx, cancel := context.WithCancel(ctx)
		// defer cancel()
		p, err := os.Open("/home/xg/baga6ea4seaqlynigdram6ljmohj3nlyzano2viwqhxyzrcl2wh4xl3day26xqcq.car")
		if err != nil {
			log.Errorf(" read piece  error: %v", err)
		}
		defer p.Close()
		if _, err := p.Seek(0, io.SeekStart); err != nil {
			log.Errorf("seek to start: %w", err)
		}
		w := &writer.Writer{}
		_, err = io.CopyBuffer(w, p, make([]byte, writer.CommPBuf))
		if err != nil {
			log.Errorf("copy into commp writer: %w", err)
		}
		commp, err := w.Sum()
		if err != nil {
			log.Errorf("computing commP failed: %w", err)
		}

		log.Info(commp.PayloadSize)
		log.Info(commp.PieceCID)
		log.Info(commp.PieceSize.Unpadded())
		return nil
	},
}
