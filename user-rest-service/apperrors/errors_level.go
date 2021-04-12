package apperrors

type level string

const (
	levelInfo     level = "info"
	levelError    level = "error"
	levelCritical level = "critical"
)

func (e *appError) LevelInfo() *appError {
	e.level = levelInfo

	return e
}

func (e *appError) LevelError() *appError {
	e.level = levelError

	return e
}

func (e *appError) LevelCritical() *appError {
	e.level = levelCritical

	return e
}

func (e *appError) IsLevelInfo() bool     { return e.checkLevel(levelInfo) }
func (e *appError) IsLevelError() bool    { return e.checkLevel(levelError) }
func (e *appError) IsLevelCritical() bool { return e.checkLevel(levelCritical) }

func (e *appError) checkLevel(lv level) bool {
	if e.level != "" {
		return e.level == lv
	}

	next := AsAppError(e.next)
	if next != nil {
		return next.checkLevel(lv)
	}

	// Default level is critical
	return lv == levelCritical
}
