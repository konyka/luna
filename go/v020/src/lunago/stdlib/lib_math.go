/*
* @Author: konyka
* @Date:   2019-05-05 12:22:18
* @Last Modified by:   konyka
* @Last Modified time: 2019-05-05 12:29:47
*/

package stdlib

import "math"
import "math/rand"
import . "luago/api"
import "luago/number"

var mathLib = map[string]GoFunction{
    "random":     mathRandom,
    "randomseed": mathRandomSeed,
    "max":        mathMax,
    "min":        mathMin,
    "exp":        mathExp,
    "log":        mathLog,
    "deg":        mathDeg,
    "rad":        mathRad,
    "sin":        mathSin,
    "cos":        mathCos,
    "tan":        mathTan,
    "asin":       mathAsin,
    "acos":       mathAcos,
    "atan":       mathAtan,
    "ceil":       mathCeil,
    "floor":      mathFloor,
    "fmod":       mathFmod,
    "modf":       mathModf,
    "abs":        mathAbs,
    "sqrt":       mathSqrt,
    "ult":        mathUlt,
    "tointeger":  mathToInt,
    "type":       mathType,
    /* placeholders */
    "pi":         nil,
    "huge":       nil,
    "maxinteger": nil,
    "mininteger": nil,
}
func OpenMathLib(ls LuaState) int {
    ls.NewLib(mathLib)
    ls.PushNumber(math.Pi)
    ls.SetField(-2, "pi")
    ls.PushNumber(math.Inf(1))
    ls.SetField(-2, "huge")
    ls.PushInteger(math.MaxInt64)
    ls.SetField(-2, "maxinteger")
    ls.PushInteger(math.MinInt64)
    ls.SetField(-2, "mininteger")
    return 1
}

/* pseudo-random numbers */

/**
 * math.random ([m [, n]])
 */
func mathRandom(ls LuaState) int {
    var low, up int64
    switch ls.GetTop() { /* check number of arguments */
    case 0: /* no arguments */
        ls.PushNumber(rand.Float64()) /* Number between 0 and 1 */
        return 1
    case 1: /* only upper limit */
        low = 1
        up = ls.CheckInteger(1)
    case 2: /* lower and upper limits */
        low = ls.CheckInteger(1)
        up = ls.CheckInteger(2)
    default:
        return ls.Error2("wrong number of arguments")
    }

    /* random integer in the interval [low, up] */
    ls.ArgCheck(low <= up, 1, "interval is empty")
    ls.ArgCheck(low >= 0 || up <= math.MaxInt64+low, 1,
        "interval too large")
    if up-low == math.MaxInt64 {
        ls.PushInteger(low + rand.Int63())
    } else {
        ls.PushInteger(low + rand.Int63n(up-low+1))
    }
    return 1
}

/**
 * math.randomseed (x)
 */
func mathRandomSeed(ls LuaState) int {
    x := ls.CheckNumber(1)
    rand.Seed(int64(x))
    return 0
}

/* max & min */

/**
 * math.max (x, ···)
 */
func mathMax(ls LuaState) int {
    n := ls.GetTop() /* number of arguments */
    imax := 1        /* index of current maximum value */
    ls.ArgCheck(n >= 1, 1, "value expected")
    for i := 2; i <= n; i++ {
        if ls.Compare(imax, i, LUA_OPLT) {
            imax = i
        }
    }
    ls.PushValue(imax)
    return 1
}

/**
 * math.min (x, ···)
 */
func mathMin(ls LuaState) int {
    n := ls.GetTop() /* number of arguments */
    imin := 1        /* index of current minimum value */
    ls.ArgCheck(n >= 1, 1, "value expected")
    for i := 2; i <= n; i++ {
        if ls.Compare(i, imin, LUA_OPLT) {
            imin = i
        }
    }
    ls.PushValue(imin)
    return 1
}

/* exponentiation and logarithms */

/**
 * math.exp (x)
 */
func mathExp(ls LuaState) int {
    x := ls.CheckNumber(1)
    ls.PushNumber(math.Exp(x))
    return 1
}

/**
 * math.log (x [, base])
 */
func mathLog(ls LuaState) int {
    x := ls.CheckNumber(1)
    var res float64

    if ls.IsNoneOrNil(2) {
        res = math.Log(x)
    } else {
        base := ls.ToNumber(2)
        if base == 2 {
            res = math.Log2(x)
        } else if base == 10 {
            res = math.Log10(x)
        } else {
            res = math.Log(x) / math.Log(base)
        }
    }

    ls.PushNumber(res)
    return 1
}

/* trigonometric functions */

/**
 * math.deg (x)
 */
func mathDeg(ls LuaState) int {
    x := ls.CheckNumber(1)
    ls.PushNumber(x * 180 / math.Pi)
    return 1
}

/**
 * math.rad (x)
 */
func mathRad(ls LuaState) int {
    x := ls.CheckNumber(1)
    ls.PushNumber(x * math.Pi / 180)
    return 1
}

/**
 * math.sin (x)
 */
func mathSin(ls LuaState) int {
    x := ls.CheckNumber(1)
    ls.PushNumber(math.Sin(x))
    return 1
}

/**
 * math.cos (x)
 */
func mathCos(ls LuaState) int {
    x := ls.CheckNumber(1)
    ls.PushNumber(math.Cos(x))
    return 1
}

/**
 * math.tan (x)
 */
func mathTan(ls LuaState) int {
    x := ls.CheckNumber(1)
    ls.PushNumber(math.Tan(x))
    return 1
}

/**
 * math.asin (x)
 */
func mathAsin(ls LuaState) int {
    x := ls.CheckNumber(1)
    ls.PushNumber(math.Asin(x))
    return 1
}

/**
 * math.acos (x)
 */
func mathAcos(ls LuaState) int {
    x := ls.CheckNumber(1)
    ls.PushNumber(math.Acos(x))
    return 1
}

/**
 * math.atan (y [, x])
 */
func mathAtan(ls LuaState) int {
    y := ls.CheckNumber(1)
    x := ls.OptNumber(2, 1.0)
    ls.PushNumber(math.Atan2(y, x))
    return 1
}

/* rounding functions */

/**
 * math.ceil (x)
 */
func mathCeil(ls LuaState) int {
    if ls.IsInteger(1) {
        ls.SetTop(1) /* integer is its own ceil */
    } else {
        x := ls.CheckNumber(1)
        _pushNumInt(ls, math.Ceil(x))
    }
    return 1
}

/**
 * math.floor (x)
 */
func mathFloor(ls LuaState) int {
    if ls.IsInteger(1) {
        ls.SetTop(1) /* integer is its own floor */
    } else {
        x := ls.CheckNumber(1)
        _pushNumInt(ls, math.Floor(x))
    }
    return 1
}

/* others */

/**
 * math.fmod (x, y)
 */
func mathFmod(ls LuaState) int {
    if ls.IsInteger(1) && ls.IsInteger(2) {
        d := ls.ToInteger(2)
        if uint64(d)+1 <= 1 { /* special cases: -1 or 0 */
            ls.ArgCheck(d != 0, 2, "zero")
            ls.PushInteger(0) /* avoid overflow with 0x80000... / -1 */
        } else {
            ls.PushInteger(ls.ToInteger(1) % d)
        }
    } else {
        x := ls.CheckNumber(1)
        y := ls.CheckNumber(2)
        ls.PushNumber(x - math.Trunc(x/y)*y)
    }

    return 1
}

/**
 * math.modf (x)
 */
func mathModf(ls LuaState) int {
    if ls.IsInteger(1) {
        ls.SetTop(1)     /* number is its own integer part */
        ls.PushNumber(0) /* no fractional part */
    } else {
        x := ls.CheckNumber(1)
        i, f := math.Modf(x)
        _pushNumInt(ls, i)
        if math.IsInf(x, 0) {
            ls.PushNumber(0)
        } else {
            ls.PushNumber(f)
        }
    }

    return 2
}

/**
 * math.abs (x)
 */
func mathAbs(ls LuaState) int {
    if ls.IsInteger(1) {
        x := ls.ToInteger(1)
        if x < 0 {
            ls.PushInteger(-x)
        }
    } else {
        x := ls.CheckNumber(1)
        ls.PushNumber(math.Abs(x))
    }
    return 1
}

/**
 * math.sqrt (x)
 */
func mathSqrt(ls LuaState) int {
    x := ls.CheckNumber(1)
    ls.PushNumber(math.Sqrt(x))
    return 1
}

/**
 * math.ult (m, n)
 */
func mathUlt(ls LuaState) int {
    m := ls.CheckInteger(1)
    n := ls.CheckInteger(2)
    ls.PushBoolean(uint64(m) < uint64(n))
    return 1
}

/**
 * math.tointeger (x)
 */
func mathToInt(ls LuaState) int {
    if i, ok := ls.ToIntegerX(1); ok {
        ls.PushInteger(i)
    } else {
        ls.CheckAny(1)
        ls.PushNil() /* value is not convertible to integer */
    }
    return 1
}

/**
 * math.type (x)
 */
func mathType(ls LuaState) int {
    if ls.Type(1) == LUA_TNUMBER {
        if ls.IsInteger(1) {
            ls.PushString("integer")
        } else {
            ls.PushString("float")
        }
    } else {
        ls.CheckAny(1)
        ls.PushNil()
    }
    return 1
}

func _pushNumInt(ls LuaState, d float64) {
    if i, ok := number.FloatToInteger(d); ok { /* does 'd' fit in an integer? */
        ls.PushInteger(i) /* result is integer */
    } else {
        ls.PushNumber(d) /* result is float */
    }
}











