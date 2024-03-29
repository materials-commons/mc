package globus

import (
	"context"

	"github.com/materials-commons/mc/pkg/globusapi"
)

type TransferItemsLoader struct {
}

func NewTransferItemsLoader() *TransferItemsLoader {
	return &TransferItemsLoader{}
}

func (l *TransferItemsLoader) LoadFiles(transferItems globusapi.TransferItems, c context.Context) error {
	for _, transfer := range transferItems.Transfers {
		if err := l.processTransferItem(transfer); err != nil {
			return err
		}
	}
	return nil
}

// GlobusResponse({'DATA_TYPE': 'successful_transfer',
// 'destination_path': '/__globus_uploads/84aa76e9-0789-c9bd-5eec-49d8cf83a1e4/hello.titan.txt',
// 'source_path': None})

func (l *TransferItemsLoader) processTransferItem(transfer globusapi.Transfer) error {
	return nil
}
