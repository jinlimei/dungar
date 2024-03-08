package markov4

import (
	"gitlab.int.magneato.site/dungar/prototype/internal/cleaner"
)

// LearnStringVariant is LearnString but exposes what cleaner.TokenVariant can be utilized
// for learning (rather than the default cleaner.VariantSlack in LearnString)
func (m *Markov) LearnStringVariant(str string, variant cleaner.TokenVariant) {
	tokenizer := cleaner.Tokenize(str, variant)
	seqs := m.makeSubSequences(tokenizer)

	for _, subseq := range seqs {
		m.LearnSequence(subseq)
	}
}

// LearnString takes the incoming string and builds out its sequence
// (via cleaner.Tokenize with cleaner.VariantSlack) and then calls LearnSequence
func (m *Markov) LearnString(str string) {
	tokenizer := cleaner.Tokenize(str, cleaner.VariantSlack)
	seqs := m.makeSubSequences(tokenizer)

	for _, subseq := range seqs {
		m.LearnSequence(subseq)
	}
}

func (m *Markov) pickRandomToken(t cleaner.TokenType) string {
	return ""
}

func (m *Markov) makeSubSequences(tknr cleaner.TokenList) [][]cleaner.Token {
	var (
		tokens = tknr.GetTokens()
		subseq = make([][]cleaner.Token, 0)
		buff   = make([]cleaner.Token, 0)
	)

	for _, token := range tokens {
		switch token.Type {
		case cleaner.TokenEOL, cleaner.TokenSpace:
		// do nothing, for now?
		case cleaner.TokenPeriod, cleaner.TokenSentEnd, cleaner.TokenNewLine:
			buff = append(buff, token)
			subseq = append(subseq, buff)
			buff = make([]cleaner.Token, 0)

		case cleaner.TokenURL, cleaner.TokenMentionRole,
			cleaner.TokenMentionUser, cleaner.TokenMentionChannel,
			cleaner.TokenEmoticon, cleaner.TokenHashTag:
			m.AddFiller(token.Type, string(token.Value))

			buff = append(buff, token)
		default:
			buff = append(buff, token)
		}
	}

	if len(buff) != 0 {
		subseq = append(subseq, buff)
	}

	return subseq
}

// LearnSequence is our primary means for learning
// new chains and input them into the Markov's various
// data entries
func (m *Markov) LearnSequence(tokens []cleaner.Token) {
	tokenIDs := make([]TokenID, len(tokens))

	for idx, token := range tokens {
		tokenID, ok := m.FindToken(token)

		if !ok {
			tokenID = m.AddToken(token)
		}

		tokenIDs[idx] = tokenID
	}

	chain := m.CreateFragments(tokenIDs)

	for _, link := range chain {
		m.AddFragment(link)
	}
}

// CreateFragments will take the list of wordIDs incoming and
// turn them into a set of weird tuples of things
func (m *Markov) CreateFragments(tokenIDs []TokenID) []Fragment {
	maxLen := len(tokenIDs)
	output := make([]Fragment, 0, maxLen)

	for pos := 0; pos <= maxLen; pos++ {
		output = append(output, Fragment{
			LWord: getTokenID(tokenIDs, maxLen, pos-1),
			CWord: getTokenID(tokenIDs, maxLen, pos),
			RWord: getTokenID(tokenIDs, maxLen, pos+1),
		})
	}

	return output
}
