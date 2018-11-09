package globus

import (
	"context"
	"fmt"
	"strings"

	"github.com/materials-commons/mc/pkg/globusapi"
)

type TransferLoader struct {
}

func NewTransferLoader() *TransferLoader {
	return &TransferLoader{}
}

func (l *TransferLoader) LoadFiles(transferItems globusapi.TransferItems, c context.Context) error {
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

func (l *TransferLoader) processTransferItem(transfer globusapi.Transfer) error {
	mcdir := ""
	fmt.Println(mcdir)
	pathItems := strings.Split(transfer.DestinationPath, "/")
	for _, pathItem := range pathItems {
		fmt.Println(pathItem)
	}

	return nil
}
