package cards

// Suit
type Suit byte

const (
	Spades Suit = iota + 1
	Hearts
	Diamonds
	Clubs
)

func (s Suit) String() string {
	switch s {
	case Spades:
		return "Spades"
	case Hearts:
		return "Hearts"
	case Diamonds:
		return "Diamonds"
	case Clubs:
		return "Clubs"
	}
	return "Nulls"
}

func (s Suit) Short() rune {
	switch s {
	case Clubs:
		return 'C' // '♣'
	case Diamonds:
		return 'D' // '♦'
	case Hearts:
		return 'H' // '♥'
	case Spades:
		return 'S' // '♠'
	}
	return '?'
}

// Rank
type Rank byte

const (
	Ace Rank = iota + 1
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

func (r Rank) String() string {
	switch r {
	case Ace:
		return "Ace"
	case Two:
		return "Two"
	case Three:
		return "Three"
	case Four:
		return "Four"
	case Five:
		return "Five"
	case Six:
		return "Six"
	case Seven:
		return "Seven"
	case Eight:
		return "Eight"
	case Nine:
		return "Nine"
	case Ten:
		return "Ten"
	case Jack:
		return "Jack"
	case Queen:
		return "Queen"
	case King:
		return "King"
	}
	return "Null"
}

func (r Rank) Short() rune {
	switch r {
	case Ace:
		return 'A'
	case Two:
		return '2'
	case Three:
		return '3'
	case Four:
		return '4'
	case Five:
		return '5'
	case Six:
		return '6'
	case Seven:
		return '7'
	case Eight:
		return '8'
	case Nine:
		return '9'
	case Ten:
		return 'T'
	case Jack:
		return 'J'
	case Queen:
		return 'Q'
	case King:
		return 'K'
	}
	return '?'
}

// Card
type Card interface {
	Suit() Suit
	Rank() Rank
	Joker() bool
	String() string
	Short() string
}

type card struct {
	rank Rank
	suit Suit
}

func (c card) Suit() Suit {
	return c.suit
}

func (c card) Rank() Rank {
	return c.rank
}

func (card) Joker() bool {
	return false
}

func (c card) String() string {
	return c.rank.String() + " of " + c.suit.String()
}

func (c card) Short() string {
	return string(c.rank.Short()) + string(c.suit.Short())
}

func NewCard(rank Rank, suit Suit) Card {
	return card{
		rank: rank,
		suit: suit,
	}
}

type joker struct{}

func (joker) Suit() Suit {
	return Suit(0)
}

func (joker) Rank() Rank {
	return Rank(0)
}

func (joker) Joker() bool {
	return true
}

func (j joker) String() string {
	return "Joker"
}

func (joker) Short() string {
	return "JO"
}

func Joker() Card {
	return joker{}
}
