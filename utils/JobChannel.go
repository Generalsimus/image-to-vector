package utils

func NewJobChannel[T any]() (
	*func(value T) T,
	func(callback func(v T) T),
	// func(callback func(v T) T),
) {
	startPont := func(value T) T {
		return value
	}
	getValue := &startPont
	addModifierAfter := func(callback func(v T) T) {
		current := *getValue
		*getValue = func(value T) T {
			return callback(current(value))
		}
	}
	// addModifierBefore := func(callback func(v T) T) {
	// 	current := *getValue
	// 	*getValue = func(value T) T {
	// 		return current(callback(value))
	// 	}
	// }
	return getValue, addModifierAfter
	// , addModifierBefore
}
