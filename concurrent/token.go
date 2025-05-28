package concurrent

import "sync"

type Token struct {
	sync.Once
	release func()
}

func NewToken(release func()) *Token {
	return &Token{
		release: release,
	}
}

func (t *Token) Release() {
	t.Do(t.release)
}
