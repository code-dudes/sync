package sync

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	errContextCancelled = fmt.Errorf("context cancelled")
	errTest             = fmt.Errorf("test error")
)

func delay(delay time.Duration) func(context.Context) error {
	return func(ctx context.Context) error {
		t := time.NewTimer(delay)
		select {
		case <-t.C:
			return errTest
		case <-ctx.Done():
			return errContextCancelled
		}
	}
}

func delayWithError(delay time.Duration, err error) func(context.Context) error {
	return func(ctx context.Context) error {
		t := time.NewTimer(delay)
		select {
		case <-t.C:
			return err
		case <-ctx.Done():
			return errContextCancelled
		}
	}
}

func Test_Cancel_WithContext(t *testing.T) {
	assert := assert.New(t)

	// Cancel parent context returns error from Go()
	{
		ctx, cancel := context.WithCancel(context.Background())
		fg, childCtx := WithContext(ctx)

		ts := time.Now()
		fg.Go(delay(time.Second))
		cancel()
		<-childCtx.Done()
		assert.Less(time.Since(ts), time.Second, "")

		err := fg.Wait()
		assert.Equal(errContextCancelled, err)
	}

	// Cancel() FuncGroup returns nil
	{
		fg, ctx := WithContext(context.Background())

		ts := time.Now()
		fg.Go(delay(time.Second))
		fg.Cancel()
		<-ctx.Done()

		assert.Less(time.Since(ts), time.Second, "")

		err := fg.Wait()
		assert.Nil(err)
	}
}

func Test_Cancel_WithNewContext(t *testing.T) {
	t.Skip()
	assert := assert.New(t)

	// Cancel() FuncGroup returns nil
	{
		fg, ctx := WithNewContext()

		ts := time.Now()
		fg.Go(delay(time.Second))
		fg.Cancel()
		<-ctx.Done()

		assert.Less(time.Since(ts), time.Second, "")

		err := fg.Wait()
		assert.Nil(err)
	}
}

func Test_MultipleGoCalls_Return_FirstError(t *testing.T) {
	t.Skip()
	assert := assert.New(t)

	{
		var e1, e2 error
		e1 = errors.New("e1")
		e2 = errors.New("e2")
		fg, _ := WithNewContext()

		fg.Go(delayWithError(time.Second, e1))
		fg.Go(delayWithError(time.Second*2, e2))

		err := fg.Wait()
		assert.Equal(e1, err)
	}
}

func Test_MultipleGoCalls_WithoutError(t *testing.T) {
	t.Skip()
	assert := assert.New(t)

	{
		fg, _ := WithNewContext()

		fg.Go(delayWithError(time.Second, nil))
		fg.Go(delayWithError(time.Second*2, nil))

		err := fg.Wait()
		assert.Nil(err)
	}
}
