package colscanner

import (
	"log"
	"reflect"
)

type binderItem struct {
	//used in scan()'s dest arguments
	PtrToHolder interface{}
	//called on colScanner::GetDest()
	ConvToResult func() error
}

type Binder = map[string]binderItem

type Factory[T any, R any] struct {
	genHolder    func() T
	genResHolder func() R
	binder       func(holder *T, result *R) Binder
}

// T and V must be a struct, not a pointer.
// binder's pointer must be a pointer to holder
func CreateFactory[T any, R any](
	binder func(holder *T, result *R) Binder,
) Factory[T, R] {
	{ //check if T and V are struct
		var t T
		var v R
		if reflect.TypeOf(t).Kind() != reflect.Struct ||
			reflect.TypeOf(v).Kind() != reflect.Struct {
			log.Fatal("ColScanner::Factory : all generics must be struct")
		}
	}
	csf := Factory[T, R]{
		genHolder:    func() T { var t T; return t },
		genResHolder: func() R { var v R; return v },
		binder:       binder,
	}
	return csf
}

func (csf Factory[T, R]) buildCS(alias *string, selectionMap map[string]bool) colScanner[R] {
	holder := csf.genHolder()
	resHolder := csf.genResHolder()
	//TODO: allow multiple col selection for binder that will added if all relevant cols selected e.g. (id, start_at), by making binder's key []string
	colMaps := csf.binder(&holder, &resHolder)

	sCols := make([]string, len(selectionMap))
	sPointerDest := make([]interface{}, len(selectionMap))
	sGenRes := make([]func() error, len(selectionMap))

	for k, item := range colMaps {
		isSelected, ok := selectionMap[k]
		if !(ok && isSelected) {
			continue
		}
		sCols = append(sCols, k)
		sPointerDest = append(sPointerDest, item.PtrToHolder)
		sGenRes = append(sGenRes, item.ConvToResult)
	}

	Cols := sCols
	if alias != nil {
		Cols = make([]string, len(sCols))
		for _, v := range sCols {
			Cols = append(Cols, (*alias)+"."+v)
		}
	}

	return colScanner[R]{
		cols:           sCols,
		Cols:           Cols,
		alias:          alias,
		dest:           sPointerDest,
		genResFuncs:    sGenRes,
		result:         &resHolder,
		selectionMap:   selectionMap,
		isResGenerated: false,
	}
}

func (csf Factory[T, R]) genSelectionMap(selections ...string) map[string]bool {
	selectionMap := make(map[string]bool)
	for _, s := range selections {
		selectionMap[s] = true
	}
	return selectionMap
}

func (csf Factory[T, R]) SelectCols(alias *string, selections ...string) colScanner[R] {
	return csf.buildCS(alias, csf.genSelectionMap(selections...))
}

func (csf Factory[T, R]) From(cs colScanner[R]) colScanner[R] {
	return csf.buildCS(cs.alias, cs.selectionMap)
}
