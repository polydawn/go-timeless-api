package ingest

import (
	"context"

	"go.polydawn.net/go-timeless-api"
)

type IngestTool func(
	ctx context.Context,
	ingestRef api.ImportRef_Ingest,
) (*api.WareID, *api.WareSourcing, error)
