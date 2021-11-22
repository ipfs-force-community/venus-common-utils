package journal

import (
	"context"
	"go.uber.org/fx"
)

func OpenFilesystemJournal(journalPath string, component string, lc fx.Lifecycle, disabled DisabledEvents) (Journal, error) {
	jrnl, err := OpenFSJournal(journalPath, component, disabled)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error { return jrnl.Close() },
	})

	return jrnl, err
}
