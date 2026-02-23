package tdler

import (
	"github.com/gotd/td/telegram/downloader"
	"github.com/kiss2u/SaveAny-Bot/common/utils/dlutil"
	"github.com/kiss2u/SaveAny-Bot/config"
	"github.com/kiss2u/SaveAny-Bot/pkg/consts/tglimit"
	"github.com/kiss2u/SaveAny-Bot/pkg/tfile"
)

func NewDownloader(file tfile.TGFile) *downloader.Builder {
	return downloader.NewDownloader().WithPartSize(tglimit.MaxPartSize).
		Download(file.Dler(), file.Location()).WithThreads(dlutil.BestThreads(file.Size(), config.C().Threads))
}
