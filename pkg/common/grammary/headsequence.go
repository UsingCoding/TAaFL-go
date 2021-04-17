package grammary

type GrammarWithHeadSequences struct {
	Grammar   Grammar
	Sequences map[Symbol][][]string
}
