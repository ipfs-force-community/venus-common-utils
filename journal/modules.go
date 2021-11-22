package journal

import (
	"context"
	"go.uber.org/fx"
)

func OpenFilesystemJournal(lc fx.Lifecycle, journalPath string, component string, disabled DisabledEvents) (Journal, error) {
	jrnl, err := OpenFSJournal(journalPath, component, disabled)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error { return jrnl.Close() },
	})

	return jrnl, err
}
