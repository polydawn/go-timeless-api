package ingest

import (
	"context"

	api "github.com/polydawn/go-timeless-api"
)

type IngestTool func(
	ctx context.Context,
	ingestRef api.ImportRef_Ingest,
) (*api.WareID, *api.WareSourcing, error)
