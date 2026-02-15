package puzzle

// Score calculates the star rating based on keystrokes and par.
func Score(keystrokes, par int) StarRating {
	if keystrokes <= par {
		return ThreeStar
	}
	if keystrokes <= par*3/2 { // floor(par * 1.5)
		return TwoStar
	}
	return OneStar
}
