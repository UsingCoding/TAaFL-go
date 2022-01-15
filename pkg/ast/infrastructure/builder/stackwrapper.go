package builder

import "compiler/pkg/ast/app"

type stackWrapper struct {
	app.Stack
	handledElements []app.StackElementPointer
}

func (w *stackWrapper) get(pointer app.StackElementPointer) app.Node {
	node := w.Stack[pointer]
	w.handledElements = append(w.handledElements, pointer)
	return node
}

func (w stackWrapper) handled(pointer app.StackElementPointer) bool {
	for _, element := range w.handledElements {
		if element == pointer {
			return true
		}
	}
	return false
}
