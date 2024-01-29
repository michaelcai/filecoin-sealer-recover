package main

import (
	// "context"

	"context"
	"fmt"
	"os"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-padreader"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/storage/sealer/ffiwrapper"
	"github.com/filecoin-project/lotus/storage/sealer/ffiwrapper/basicfs"
	"github.com/filecoin-project/lotus/storage/sealer/partialfile"
	"github.com/filecoin-project/lotus/storage/sealer/storiface"
	"github.com/michaelcai/filecoin-sealer-recover/export"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/xerrors"

	// cliutil "github.com/filecoin-project/lotus/cli/util"
	logging "github.com/ipfs/go-log/v2"
)

// const (
// 	ss32GiB = 32 << 30
// )

var log = logging.Logger("recover")

func main() {
	sector := export.SectorInfo{
		SectorNumber: 34617,
		SealProof:    8,
		Ticket:       abi.Randomness("PP+UvrZn96434UWUJzP2zU58HXgNr8M++33LcgFLCHU="),
	}
	sealingTemp := "./tmp"
	sdir, err := homedir.Expand(sealingTemp)
	if err != nil {
		log.Errorf("Sector (%d) ,expands the path error: %v", sector.SectorNumber, err)
	}
	mkdirAll(sdir)
	tempDir, err := os.MkdirTemp(sdir, fmt.Sprintf("recover-%d", sector.SectorNumber))
	if err != nil {
		log.Errorf("Sector (%d) ,creates a new temporary directory error: %v", sector.SectorNumber, err)
	}

	if err := os.MkdirAll(tempDir, 0775); err != nil {
		log.Errorf("Sector (%d) ,creates a directory named path error: %v", sector.SectorNumber, err)
	}

	provider := basicfs.Provider{
		Root: tempDir,
	}
	sb, err := ffiwrapper.New(&provider)

	addr, err := address.NewFromString("f02359213")
	if err != nil {
		xerrors.Errorf("NewFromString err: %w", err)
	}
	actorID, err := address.IDFromAddress(addr)
	if err != nil {
		xerrors.Errorf("Getting IDFromAddress err: %w", err)
	}

	sid := storiface.SectorRef{
		ID: abi.SectorID{
			Miner:  abi.ActorID(actorID),
			Number: sector.SectorNumber,
		},
		ProofType: sector.SealProof,
	}

	stagedPath, done, err := provider.AcquireSector(context.Background(), sid, 0, storiface.FTUnsealed, storiface.PathSealing)
	defer done()

	ssize, err := sid.ProofType.SectorSize()
	maxPieceSize := abi.PaddedPieceSize(ssize)
	stagedFile, err := partialfile.CreatePartialFile(maxPieceSize, stagedPath.Unsealed)
	fmt.Println(stagedFile, stagedPath)

	sb.NewSector(context.Background(), sid)
	fmt.Println()
	unsealFile, err := os.Open("./tmp/s-t02359213-34617")
	if err != nil {
		log.Errorf(" read piece  error: %v", err)
	}
	defer unsealFile.Close()
	info, err := unsealFile.Stat()
	paddedReader, paddedSize := padreader.New(unsealFile, uint64(info.Size()))
	pi, err := sb.DataCid(context.TODO(), paddedSize, paddedReader)
	fmt.Println(pi)
	if err != nil {
		log.Errorf("Sector (%d) , AP  error: %v", sector.SectorNumber, err)
	}
	pc1o, err := sb.SealPreCommit1(context.TODO(), sid, abi.SealRandomness(sector.Ticket), []abi.PieceInfo{pi})
	if err != nil {
		log.Errorf("Sector (%d) , running PreCommit1  error: %v", sector.SectorNumber, err)
	}
	log.Infof("Complete PreCommit1, sector (%d)", sector.SectorNumber)
	fmt.Println(pc1o)
	// if err != nil {
	// 	log.Errorf("Sector (%d) , os.Stat  error: %v", sector.SectorNumber, err)
	// }
	// log.Info(commp.PayloadSize)
	// log.Info(commp.PieceCID)
	// log.Info(commp.PieceSize.Unpadded())
}

func mkdirAll(path string) {
	_, err := os.Stat(path)
	notexist := os.IsNotExist(err)

	if notexist {
		_ = os.MkdirAll(path, 0755) //nolint: gosec
	}
}
