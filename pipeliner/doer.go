package pipeliner

import (
	"context"
	"reflect"
	"sync"
)

type doer func(context.Context, []*pipelinerCmd) error

func getdoer(f interface{}) doer {
	fv := reflect.ValueOf(f)
	ft := fv.Type()

	if ft.Kind() != reflect.Func ||
		(ft.NumIn() < 1 || ft.NumIn() > 2) ||
		(ft.NumIn() == 1 &&
			ft.In(0).Kind() != reflect.Slice) ||
		(ft.NumIn() == 2 &&
			ft.In(0) != reflect.TypeOf((*context.Context)(nil)).Elem() &&
			ft.In(1).Kind() != reflect.Slice) ||
		ft.NumOut() != 1 ||
		ft.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
		panic(`f must be "func([]T) error" or "func(context.Context, []T)" error`)
	}

	var (
		contextEnabled = ft.NumIn() == 2
		dataIdx        = ft.NumIn() - 1
		pool           sync.Pool
	)

	return func(ctx context.Context, cmds []*pipelinerCmd) error {
		// step 1. get slice from pool
		in, _ := pool.Get().(*reflect.Value)
		if in == nil {
			_in := reflect.New(ft.In(dataIdx))
			_in.Elem().Set(reflect.MakeSlice(_in.Elem().Type(), 0, len(cmds)))
			in = &_in
		}

		// step 2. reset the slice before use
		in.Elem().SetLen(0)

		// step 3. set the data to the slice
		for _, cmd := range cmds {
			in.Elem().Set(reflect.Append(in.Elem(), reflect.ValueOf(cmd.v)))
		}

		// step 4. execute the function
		var err error
		if contextEnabled {
			err, _ = fv.Call([]reflect.Value{reflect.ValueOf(ctx), in.Elem()})[0].Interface().(error)
		} else {
			err, _ = fv.Call([]reflect.Value{in.Elem()})[0].Interface().(error)
		}
		pool.Put(in)
		return err
	}

}
