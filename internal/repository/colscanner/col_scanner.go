package colscanner

type colScanner[R any] struct {
	result *R
	//store raw unaliased col
	cols []string
	//use this in sql
	Cols []string
	//use this for scan()'s args.
	//if there is multiple col scanner, make sure the order is same as in sql
	dest []interface{}
	//MUST: call this after scan()
	genResFuncs []func() error
	//flag, set true after genResult() called
	isResGenerated bool
	//will be used incase csf using From() on this cs
	selectionMap map[string]bool
	alias        *string
}

func (c *colScanner[R]) getCols() []string {
	return c.cols
}

func (cs *colScanner[R]) genResult() error {
	for _, g := range cs.genResFuncs {
		if err := g(); err != nil {
			return err
		}
	}
	cs.isResGenerated = true
	return nil
}

func (cs *colScanner[R]) getDest() []interface{} {
	return cs.dest
}

// call this ONLY after using ScanRow() with this colScanner as arg.
//
// As long as you sure that this called after ScanRow, then you can
// ignore the hasScanned return value
func (cs *colScanner[R]) GetResult() (result R, hasScanned bool) {
	if !cs.isResGenerated {
		var r R
		return r, false
	}
	return *cs.result, true
}
