package colscanner

import (
	"fmt"

	"github.com/labstack/gommon/log"
)

type csInterface interface {
	getCols() []string
	genResult() error
	getDest() []interface{}
}

// Only takes string and pointer to colScanner. Will panic if otherwise.
// This Function will join all the args for squirrel's Select query
func CreateSelect(args ...any) []string {
	dest := make([]string, len(args))
	for _, arg := range args {
		cs, isCs := arg.(csInterface)
		if isCs {
			dest = append(dest, cs.getCols()...)
		} else {
			str, ok := arg.(string)
			if ok {
				dest = append(dest, str)
			} else {
				log.Fatal("colScanner::join cols args only takes string and pointer")
			}
		}
	}
	return dest
}

type scannable interface {
	Scan(dest ...any) error
}

// Make sure the order of args is like the order of sql's select.
// Args can only take pointer, whether its colScanner or other type.
// Args then will be used as scan dest
func ScanRow(row scannable, args ...any) error {
	colScanners := make([]csInterface, len(args))
	dest := make([]interface{}, len(args))
	for _, arg := range args {
		cs, isCs := arg.(csInterface)
		if isCs {
			colScanners = append(colScanners, cs)
			dest = append(dest, cs.getDest()...)
		} else {
			dest = append(dest, arg)
		}
	}

	err := row.Scan(dest...)
	if err != nil {
		return err
	}

	for _, cs := range colScanners {
		err := cs.genResult()
		if err != nil {
			return fmt.Errorf("ColScanner::ScanRow : %v", err)
		}
	}
	return nil
}
