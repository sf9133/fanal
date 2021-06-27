package artifact

import (
	"context"

	"github.com/sf9133/fanal/types"
)

type Artifact interface {
	Inspect(ctx context.Context) (reference types.ArtifactReference, err error)
}
