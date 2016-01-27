package limiter

import (
	"errors"
	"runtime"
)

/*
AdjustableTokenChanLimiter extends TokenChanLimiter with methods enabling the addition and removal of tokens, and satisfies the TokenLimiter and InvocationLimiter interfaces.
*/
type AdjustableTokenChanLimiter struct {
	TokenChanLimiter
	maxTokenCount uint
	tokenCount    uint
}

func NewAdjustableTokenChanLimiter(initialTokens uint, maxTokens uint) (tl *AdjustableTokenChanLimiter) {
	tl = &AdjustableTokenChanLimiter{
		maxTokenCount: maxTokens,
	}
	tl.tokens = make(chan *[16]byte, int(maxTokens))
	tl.AddTokens(initialTokens)
	return
}

func NewAdjustableCpuTokenChanLimiter(initialTokens uint) (tl *AdjustableTokenChanLimiter) {
	maxTokens := uint(runtime.NumCPU())
	if maxTokens < initialTokens {
		maxTokens = initialTokens
	}
	tl = NewAdjustableTokenChanLimiter(initialTokens, maxTokens)
	return
}

/*
GetTokenCount returns the current token count.
*/
func (l *AdjustableTokenChanLimiter) GetTokenCount() uint {
	return l.tokenCount
}

/*
AddTokens creates the specified number of new tokens and adds them to the limiter's supply channel.

If the maximum number of tokens are reached, an error will be returned.
*/
func (l *AdjustableTokenChanLimiter) AddTokens(count uint) (err error) {
	if count == 0 {
		return
	}
	if l.maxTokenCount <= l.tokenCount {
		err = errors.New("Token count maximum has been reached.")
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.tokenCount+count > l.maxTokenCount {
		err = errors.New("Token count maximum has been reached.")
	}
	for i := uint(0); i < count && l.tokenCount < l.maxTokenCount; i += 1 {
		l.tokens <- new([16]byte)
		l.tokenCount += 1
	}
	return
}

/*
RemoveTokens removes up to the specified number of tokens from the limiter's supply channel.

If the number of tokens reaches 0, no action will be taken and no error will be returned. Passing math.MaxUint64 instead of a known token count will safely remove all tokens.
*/
func (l *AdjustableTokenChanLimiter) RemoveTokens(count uint) (err error) {
	if count == 0 || l.tokenCount == 0 {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	for i := uint(0); i < count && l.tokenCount > 0; i += 1 {
		<-l.tokens
		l.tokenCount -= 1
	}
	return
}
