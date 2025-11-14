package dispatcher

import (
	"context"

	"github.com/martencassel/opsmsg/message"
)

type Dispatcher interface {
	Dispatch(ctx context.Context, msg message.Message) error
}
