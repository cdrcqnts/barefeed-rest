package cnt

const (
	ErrInvalidSID     = "Invalid SID."
	ErrUrlExists      = "This URL already exists for this SID."
	ErrSIDCIDNotFound = "No entry found for the given SID and CID."
	ErrParseURL       = "Failed parsing URL. Make sure the URL provides a valid feed."
	ErrNoAudioFile    = "It seems that this feed does not provide audio files."
)
