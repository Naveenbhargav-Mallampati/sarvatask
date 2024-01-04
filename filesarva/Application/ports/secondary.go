package ports

import "context"

type RaftPort interface {
	RaftConsensus(ctx context.Context, fileInfo *FileInfo) (string, error)
}
