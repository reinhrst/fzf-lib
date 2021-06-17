package fzf

// Case denotes case-sensitivity of search
type Case int

const (
	CaseSmart Case = iota
	CaseIgnore
	CaseRespect
)

// Sort criteria
type Criterion int

const (
	ByScore Criterion = iota
	ByLength
	ByBegin
	ByEnd
)

func isAlphabet(char uint8) bool {
	return char >= 'a' && char <= 'z'
}

func isNumeric(char uint8) bool {
	return char >= '0' && char <= '9'
}
