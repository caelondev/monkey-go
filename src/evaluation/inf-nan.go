package evaluation

import (
	"math"

	"github.com/caelondev/monkey/src/object"
	"github.com/caelondev/monkey/src/token"
)

// NOTE: This file only contains Infinity and NaN's semantics ---
// I (caelondev) placed it on a seperate file because the NaN and ---
// Inf rule is so massive that it's worth placing it on a new file ---

func infinityWithSign(sign int) object.Object {
	if sign >= 0 {
		return object.INFINITY
	}
	return object.NEG_INFINITY
}

func signFromNumber(n float64) int {
	if math.Signbit(n) {
		return -1
	}
	return 1
}

func evalInfInf(op token.TokenType, l, r *object.Infinity) object.Object {
	switch op {
	case token.PLUS:
		if l.Sign == r.Sign {
			return infinityWithSign(l.Sign)
		}
		return object.NAN

	case token.MINUS:
		if l.Sign == r.Sign {
			return object.NAN
		}
		return infinityWithSign(l.Sign)

	case token.STAR:
		return infinityWithSign(l.Sign * r.Sign)

	case token.SLASH:
		return object.NAN

	case token.CARET:
		if r.Sign < 0 {
			return &object.Number{Value: 0}
		}
		return object.INFINITY

	case token.EQUAL:
		return eBool(l.Sign == r.Sign)
	case token.NOT_EQUAL:
		return eBool(l.Sign != r.Sign)
	case token.LESS:
		return eBool(l.Sign < r.Sign)
	case token.GREATER:
		return eBool(l.Sign > r.Sign)
	case token.LESS_EQUAL:
		return eBool(l.Sign <= r.Sign)
	case token.GREATER_EQUAL:
		return eBool(l.Sign >= r.Sign)
	}

	return object.NIL
}

func evalInfNum(op token.TokenType, inf *object.Infinity, num *object.Number) object.Object {
	switch op {
	case token.PLUS:
		return infinityWithSign(inf.Sign)

	case token.MINUS:
		return infinityWithSign(inf.Sign)

	case token.STAR:
		if num.Value == 0 {
			return object.NAN
		}
		return infinityWithSign(inf.Sign * signFromNumber(num.Value))

	case token.SLASH:
		if num.Value == 0 {
			return infinityWithSign(inf.Sign)
		}
		return infinityWithSign(inf.Sign * signFromNumber(num.Value))

	case token.CARET:
		if num.Value == 0 {
			return &object.Number{Value: 1}
		}
		if num.Value < 0 {
			return &object.Number{Value: 0}
		}
		if inf.Sign < 0 && num.Value != math.Floor(num.Value) {
			return object.NAN
		}
		if inf.Sign < 0 && int(num.Value)%2 == 0 {
			return object.INFINITY
		}
		return infinityWithSign(inf.Sign)

	case token.EQUAL:
		return object.FALSE
	case token.NOT_EQUAL:
		return object.TRUE
	case token.LESS:
		return eBool(inf.Sign < 0)
	case token.GREATER:
		return eBool(inf.Sign > 0)
	case token.LESS_EQUAL:
		return eBool(inf.Sign < 0)
	case token.GREATER_EQUAL:
		return eBool(inf.Sign > 0)
	}

	return object.NAN
}

func evalNumInf(op token.TokenType, num *object.Number, inf *object.Infinity) object.Object {
	switch op {
	case token.PLUS:
		return infinityWithSign(inf.Sign)

	case token.MINUS:
		return infinityWithSign(-inf.Sign)

	case token.STAR:
		if num.Value == 0 {
			return object.NAN
		}
		return infinityWithSign(inf.Sign * signFromNumber(num.Value))

	case token.SLASH:
		return &object.Number{Value: math.Copysign(0, num.Value)}

	case token.CARET:
		absNum := math.Abs(num.Value)

		if num.Value == 0 {
			if inf.Sign > 0 {
				return &object.Number{Value: 0}
			}
			return object.INFINITY
		}

		if absNum == 1 {
			if num.Value == 1 {
				return &object.Number{Value: 1}
			}
			return object.NAN
		}

		if inf.Sign > 0 {
			if absNum < 1 {
				return &object.Number{Value: 0}
			}
			return infinityWithSign(signFromNumber(num.Value))
		} else {
			if absNum < 1 {
				return object.INFINITY
			}
			return &object.Number{Value: 0}
		}

	case token.EQUAL:
		return object.FALSE
	case token.NOT_EQUAL:
		return object.TRUE
	case token.LESS:
		return eBool(inf.Sign > 0)
	case token.GREATER:
		return eBool(inf.Sign < 0)
	case token.LESS_EQUAL:
		return eBool(inf.Sign > 0)
	case token.GREATER_EQUAL:
		return eBool(inf.Sign < 0)
	}

	return object.NAN
}

func eBool(v bool) *object.Boolean {
	if v {
		return object.TRUE
	}
	return object.FALSE
}
