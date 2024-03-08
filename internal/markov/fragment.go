package markov

type fragment struct {
	ID         int64
	RWordID    int64
	LWordID    int64
	SentenceID int64
	WordID     int64
	Word       string
}

func (frag *fragment) getWordIds() []int64 {
	return []int64{
		frag.LWordID,
		frag.WordID,
		frag.RWordID,
	}
}

func (frag *fragment) strAttr(name string) string {
	switch name {
	case "Word":
		return frag.Word
	default:
		panic("Unknown string column name " + name)
	}
}

func (frag *fragment) column(name string) string {
	switch name {
	case "ID":
		return "id"
	case "LWordID":
		return "l_word_id"
	case "RWordID":
		return "r_word_id"
	case "SentenceID":
		return "sentence_id"
	case "WordID":
		return "word_id"
	case "Word":
		return "word"
	default:
		panic("Unknown column " + name)
	}
}

func (frag *fragment) intAttr(name string) int64 {
	switch name {
	case "ID":
		return frag.ID
	case "LWordID":
		return frag.LWordID
	case "RWordID":
		return frag.RWordID
	case "SentenceID":
		return frag.SentenceID
	case "WordID":
		return frag.WordID
	default:
		panic("Unknown int column name " + name)
	}
}
