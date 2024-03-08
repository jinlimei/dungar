package cleaner

// TokenVariant declares what sort of parsing should occur on
// individual tokens when working with the data at hand
type TokenVariant uint8

const (
	// VariantUndefined is when we have not defined the variant (default zero-case0
	VariantUndefined TokenVariant = iota
	// VariantSlack is for Slack incoming data
	VariantSlack
	// VariantTwitter is for Twitter incoming data (tweets)
	VariantTwitter
	// VariantDiscord is for Discord messages
	VariantDiscord
	// VariantMatrix is for Matrix messages
	VariantMatrix
	// VariantBook is for our alice_in_wonderland messaging
	// (and potentially others?)
	VariantBook
	// VariantPlain is for just plain incoming text
	VariantPlain
	// VariantXMPP is for Jabber/XMPP handling of incoming text (ugh)
	VariantXMPP
	// VariantSonnet is for ... sonnets
	VariantSonnet
	// VariantBible is for the bible (King James, so far)
	VariantBible
)

func (tv TokenVariant) isHTMLVariant() bool {
	return tv == VariantSlack || tv == VariantMatrix
}

func (tv TokenVariant) isMarkdownVariant() bool {
	return tv == VariantSlack || tv == VariantTwitter ||
		tv == VariantDiscord || tv == VariantMatrix
}

func (tv TokenVariant) isSimpleVariant() bool {
	return tv == VariantPlain || tv == VariantBook ||
		tv == VariantSonnet || tv == VariantBible
}

func (tv TokenVariant) isEmojiVariant() bool {
	return tv == VariantSlack || tv == VariantDiscord ||
		tv == VariantXMPP
}
