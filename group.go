package cards

import (
	"sort"
	"strings"
)

// Group represents a group of cards. Each card in a group has a direction,
// which can be either face-up or face-down. When indexing cards, they are
// positioned bottom-to-top, with index 0 being the bottom, and the last index
// being the top.
//
// If an index is less than 0, then the index starts from the top (-1 becomes
// group.Len() - 1).
//
// Note that when drawing a number of cards, they are drawn as a group, rather
// than individually.
//
// Group implements sort.Interface, which includes the methods Len, Less, and
// Swap.
type Group interface {
	sort.Interface

	// Flipped returns true if the card at position i is face-up, and false if
	// it is face-down.
	Flipped(i int) bool
	// SetFlipped sets the direction of the card at position i.
	SetFlipped(i int, faceup bool)
	// FlippedArray returns the directions of all cards in the group.
	FlippedArray() []bool
	// FlipEach sets the direction of each card in the group. Returns the
	// group itself.
	FlipEach(faceup bool) Group
	// Flip takes cards in the range i:j, and "flips" them around as a group.
	// That is, the direction of each card is toggled, and the order of the
	// cards within the range is reversed.
	Flip(i, j int)
	// FlipAll flips the entire group. Returns the group itself.
	FlipAll() Group

	// Card returns the card at position i in the group, or nil if there is no
	// card at position i.
	Card(i int) Card
	// Cards returns a list of the cards in the group.
	Cards() []Card

	// Draw removes n cards from the top of the group, returning them in a new
	// group. The cards in the new group will have the same direction.
	Draw(n int) Group
	// DrawBottom removes n cards from the bottom of the group, returning them
	// in a new group. The cards in the new group will have the same
	// direction.
	DrawBottom(n int) Group
	// DrawAt removes cards in the range i:j from the group, returning them
	// in a new group. The cards in the new group will have the same
	// direction.
	DrawAt(i, j int) Group

	// Insert adds cards to the top of the group.
	Insert(g Group)
	// InsertBottom adds cards to the bottom of the group.
	InsertBottom(g Group)
	// InsertAt adds cards to position i in the group.
	InsertAt(i int, g Group)

	// Returns the group as a string.
	String() string
}

type group struct {
	cards []Card
	flipd []bool
}

// NewGroup returns a list of cards as a Group.
func NewGroup(cards ...Card) Group {
	return &group{
		cards: cards,
		flipd: make([]bool, len(cards)),
	}
}

// NewStandardDeck returns a deck as a group, containing the standard 52
// cards.
func NewStandardDeck() Group {
	deck := &group{
		cards: make([]Card, 52),
		flipd: make([]bool, 52),
	}
	for suit := 1; suit <= 4; suit++ {
		for rank := 1; rank <= 13; rank++ {
			deck.cards[(suit-1)*13+(rank-1)] = NewCard(Rank(rank), Suit(suit))
		}
	}
	return deck
}

func (g *group) index(i *int) {
	if *i < 0 {
		*i += len(g.cards)
	}
}

func (g *group) String() string {
	s := make([]string, len(g.cards))
	for i := 0; i < len(g.cards); i++ {
		s[i] = g.cards[i].Short()
	}
	return "[ " + strings.Join(s, " ") + " ]"
}

func (g *group) Len() int {
	return len(g.cards)
}

func (g *group) Less(i, j int) bool {
	g.index(&i)
	g.index(&j)
	ci, cj := g.cards[i], g.cards[j]
	if ci.Joker() && !cj.Joker() {
		return true
	} else if !ci.Joker() && cj.Joker() {
		return false
	}
	if ci.Suit() == cj.Suit() {
		return ci.Rank() < cj.Rank()
	}
	return ci.Suit() < cj.Suit()
}

func (g *group) Swap(i, j int) {
	g.index(&i)
	g.index(&j)
	g.cards[i], g.cards[j] = g.cards[j], g.cards[i]
	g.flipd[i], g.flipd[j] = g.flipd[j], g.flipd[i]
}

func (g *group) Flipped(i int) bool {
	g.index(&i)
	return g.flipd[i]
}

func (g *group) SetFlipped(i int, faceup bool) {
	g.index(&i)
	g.flipd[i] = faceup
}

func (g *group) FlippedArray() []bool {
	c := make([]bool, len(g.flipd))
	copy(c, g.flipd)
	return c
}

func (g *group) FlipEach(faceup bool) Group {
	for i := range g.flipd {
		g.flipd[i] = faceup
	}
	return g
}

func (g *group) flip(i, j int) {
	for n := i; n < j; n++ {
		g.flipd[n] = !g.flipd[n]
	}
	c := g.cards[i:j]
	d := g.flipd[i:j]
	for p := len(c)/2 - 1; p >= 0; p-- {
		q := len(c) - 1 - p
		c[p], c[q] = c[q], c[p]
		d[p], d[q] = d[q], d[p]
	}
}

func (g *group) Flip(i, j int) {
	g.index(&i)
	g.index(&j)
	g.flip(i, j)
}

func (g *group) FlipAll() Group {
	g.flip(0, len(g.cards))
	return g
}

func (g *group) Card(i int) Card {
	g.index(&i)
	if i < 0 || i >= len(g.cards) {
		return nil
	}
	return g.cards[i]
}

func (g *group) Cards() []Card {
	c := make([]Card, len(g.cards))
	copy(c, g.cards)
	return c
}

func (g *group) Draw(n int) Group {
	if n < 0 {
		n = 0
	}
	if n > len(g.cards) {
		n = len(g.cards)
	}
	return g.drawAt(len(g.cards)-n, len(g.cards))
}

func (g *group) DrawBottom(n int) Group {
	if n < 0 {
		n = 0
	}
	if n > len(g.cards) {
		n = len(g.cards)
	}
	return g.drawAt(0, n)
}

func (g *group) DrawAt(i, j int) Group {
	g.index(&i)
	g.index(&j)
	return g.drawAt(i, j)
}

func (g *group) drawAt(i, j int) Group {
	n := j - i
	ng := &group{cards: make([]Card, n), flipd: make([]bool, n)}
	copy(ng.cards, g.cards[i:j])
	copy(ng.flipd, g.flipd[i:j])

	copy(g.cards[i:], g.cards[j:])
	copy(g.cards[len(g.cards)-n:], make([]Card, n))
	g.cards = g.cards[:len(g.cards)-n]

	copy(g.flipd[i:], g.flipd[j:])
	g.flipd = g.flipd[:len(g.flipd)-n]

	return ng
}

func (g *group) Insert(ng Group) {
	g.cards = append(g.cards, ng.Cards()...)
	g.flipd = append(g.flipd, ng.FlippedArray()...)
}

func (g *group) InsertBottom(ng Group) {
	c := make([]Card, len(g.cards)+ng.Len())
	copy(c, ng.Cards())
	copy(c[ng.Len():], g.cards)
	g.cards = c

	d := make([]bool, len(g.flipd)+ng.Len())
	copy(d, ng.FlippedArray())
	copy(d[ng.Len():], g.flipd)
	g.flipd = d
}

func (g *group) InsertAt(i int, ng Group) {
	g.index(&i)
	if i < 0 {
		i = 0
	} else if i > len(g.cards) {
		i = len(g.cards)
	}

	c := make([]Card, len(g.cards)+ng.Len())
	copy(c, g.cards[:i])
	copy(c[i:], ng.Cards())
	copy(c[i+ng.Len():], g.cards[i:])
	g.cards = c

	d := make([]bool, len(g.flipd)+ng.Len())
	copy(d, g.flipd[:i])
	copy(d[i:], ng.FlippedArray())
	copy(d[i+ng.Len():], g.flipd[i:])
	g.flipd = d
}
