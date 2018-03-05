package main

import "strings"

type castlings uint

const (
	shortW = uint(0x1) // white can castle short
	longW  = uint(0x2) // white can castle long
	shortB = uint(0x4) // black can castle short
	longB  = uint(0x8) // black can castle short
)

func (c *castlings) on(val uint){
	(*c) |= castlings(val)
}

func (c *castlings) off(val uint){
	(*c) &= castlings(^val)
}

// parse castling rights in fenstring
func parseCastlings(fenCastl string) castlings {
	c := uint(0)
	if fenCastl == "-" {
		return castlings(0)
	}

	if strings.Index(fenCastl, "K") >= 0 {
		c |= shortW
	}
	if strings.Index(fenCastl, "Q") >= 0 {
		c |= longW
	}
	if strings.Index(fenCastl, "k") >= 0 {
		c |= shortB
	}
	if strings.Index(fenCastl, "q") >= 0 {
		c |= longB
	}

	return castlings(c)
}
