// Slice object

package py

// A python Slice object
type Slice struct {
	Start Object
	Stop  Object
	Step  Object
}

var SliceType = NewType("slice", `slice(stop) -> slice object
"slice(stop)
slice(start, stop[, step])

Create a slice object.  This is used for extended slicing (e.g. a[0:10:2]).`)

// Type of this object
func (o *Slice) Type() *Type {
	return SliceType
}

// Make a new slice object
func NewSlice(start, stop, step Object) *Slice {
	return &Slice{
		Start: start,
		Stop:  stop,
		Step:  step,
	}
}

// SliceNew
func SliceNew(metatype *Type, args Tuple, kwargs StringDict) Object {
	var start Object = None
	var stop Object = None
	var step Object = None
	UnpackTuple(args, kwargs, "slice", 1, 3, &start, &stop, &step)
	if len(args) == 1 {
		return NewSlice(None, start, None)
	}
	return NewSlice(start, stop, step)
}

// GetIndices
//
// Retrieve the start, stop, and step indices from the slice object
// slice assuming a sequence of length length, and store the length of
// the slice in slicelength. Out of bounds indices are clipped in a
// manner consistent with the handling of normal slices.
func (r *Slice) GetIndices(length int) (start, stop, step, slicelength int) {
	var defstart, defstop int

	if r.Step == None {
		step = 1
	} else {
		step = IndexInt(r.Step)
		if step == 0 {
			panic(ExceptionNewf(ValueError, "slice step cannot be zero"))
		}
		const PY_SSIZE_T_MAX = int(^uint(0) >> 1)
		/* Here *step might be -PY_SSIZE_T_MAX-1; in this case we replace it
		 * with -PY_SSIZE_T_MAX.  This doesn't affect the semantics, and it
		 * guards against later undefined behaviour resulting from code that
		 * does "step = -step" as part of a slice reversal.
		 */
		if step < -PY_SSIZE_T_MAX {
			step = -PY_SSIZE_T_MAX
		}
	}

	if step < 0 {
		defstart = length - 1
		defstop = -1
	} else {
		defstart = 0
		defstop = length
	}

	if r.Start == None {
		start = defstart
	} else {
		start = IndexInt(r.Start)
		if start < 0 {
			start += length
		}
		if start < 0 {
			if step < 0 {

				start = -1
			} else {
				start = 0
			}
		}
		if start >= length {
			if step < 0 {
				start = length - 1
			} else {
				start = length
			}
		}
	}

	if r.Stop == None {
		stop = defstop
	} else {
		stop = IndexInt(r.Stop)
		if stop < 0 {
			stop += length
		}
		if stop < 0 {
			if step < 0 {
				stop = -1
			} else {
				stop = 0
			}
		}
		if stop >= length {
			if step < 0 {
				stop = length - 1
			} else {
				stop = length
			}
		}
	}

	if (step < 0 && stop >= start) || (step > 0 && start >= stop) {
		slicelength = 0
	} else if step < 0 {
		slicelength = (stop-start+1)/(step) + 1
	} else {
		slicelength = (stop-start-1)/(step) + 1
	}

	return
}

// Check interface is satisfied
