package types

func VoteOptionFromString(str string) VoteOption {
	switch str {
	case "yes", "Yes":
		return VoteOption_VOTE_OPTION_YES
	case "no", "No":
		return VoteOption_VOTE_OPTION_NO
	case "abstain", "Abstain":
		return VoteOption_VOTE_OPTION_ABSTAIN
	}
	return VoteOption_VOTE_OPTION_UNSPECIFIED
}
