package main

import (
	"context"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/exp/rand"
)

type AppRuntime struct {
	conversationID int
	app            *App
}

func (r *AppRuntime) WaitForApproval(ctx context.Context, message string) error {
	approvalChan := make(chan struct{}, 1)
	approvalID := make([]byte, 8)
	rand.Read(approvalID)
	r.app.m.Lock()
	r.app.pendingApprovalRequests[r.conversationID] = approvalRequest{
		approvalID:   fmt.Sprintf("%x", approvalID),
		approvalChan: approvalChan,
		message:      message,
	}
	r.app.m.Unlock()
	runtime.EventsEmit(ctx, fmt.Sprintf("conversation-%d-approvals-updated", r.conversationID))
	defer func() {
		r.app.m.Lock()
		delete(r.app.pendingApprovalRequests, r.conversationID)
		r.app.m.Unlock()
		runtime.EventsEmit(ctx, fmt.Sprintf("conversation-%d-approvals-updated", r.conversationID))
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-approvalChan:
		return nil
	}
}
