// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package anthropic_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/internal"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type closureTransport struct {
	fn func(req *http.Request) (*http.Response, error)
}

func (t *closureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.fn(req)
}

func TestUserAgentHeader(t *testing.T) {
	var userAgent string
	client := anthropic.NewClient(
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					userAgent = req.Header.Get("User-Agent")
					return &http.Response{
						StatusCode: http.StatusOK,
					}, nil
				},
			},
		}),
	)
	client.Messages.New(context.Background(), anthropic.MessageNewParams{
		MaxTokens: anthropic.F(int64(1024)),
		Messages: anthropic.F([]anthropic.MessageParam{{
			Role:    anthropic.F(anthropic.MessageParamRoleUser),
			Content: anthropic.F([]anthropic.MessageParamContentUnion{anthropic.TextBlockParam{Text: anthropic.F("What is a quaternion?"), Type: anthropic.F(anthropic.TextBlockParamTypeText)}}),
		}}),
		Model: anthropic.F(anthropic.ModelClaude_3_5_Sonnet_20240620),
	})
	if userAgent != fmt.Sprintf("Anthropic/Go %s", internal.PackageVersion) {
		t.Errorf("Expected User-Agent to be correct, but got: %#v", userAgent)
	}
}

func TestRetryAfter(t *testing.T) {
	attempts := 0
	client := anthropic.NewClient(
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					attempts++
					return &http.Response{
						StatusCode: http.StatusTooManyRequests,
						Header: http.Header{
							http.CanonicalHeaderKey("Retry-After"): []string{"0.1"},
						},
					}, nil
				},
			},
		}),
	)
	res, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		MaxTokens: anthropic.F(int64(1024)),
		Messages: anthropic.F([]anthropic.MessageParam{{
			Role:    anthropic.F(anthropic.MessageParamRoleUser),
			Content: anthropic.F([]anthropic.MessageParamContentUnion{anthropic.TextBlockParam{Text: anthropic.F("What is a quaternion?"), Type: anthropic.F(anthropic.TextBlockParamTypeText)}}),
		}}),
		Model: anthropic.F(anthropic.ModelClaude_3_5_Sonnet_20240620),
	})
	if err == nil || res != nil {
		t.Error("Expected there to be a cancel error and for the response to be nil")
	}
	if want := 3; attempts != want {
		t.Errorf("Expected %d attempts, got %d", want, attempts)
	}
}

func TestRetryAfterMs(t *testing.T) {
	attempts := 0
	client := anthropic.NewClient(
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					attempts++
					return &http.Response{
						StatusCode: http.StatusTooManyRequests,
						Header: http.Header{
							http.CanonicalHeaderKey("Retry-After-Ms"): []string{"100"},
						},
					}, nil
				},
			},
		}),
	)
	res, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		MaxTokens: anthropic.F(int64(1024)),
		Messages: anthropic.F([]anthropic.MessageParam{{
			Role:    anthropic.F(anthropic.MessageParamRoleUser),
			Content: anthropic.F([]anthropic.MessageParamContentUnion{anthropic.TextBlockParam{Text: anthropic.F("What is a quaternion?"), Type: anthropic.F(anthropic.TextBlockParamTypeText)}}),
		}}),
		Model: anthropic.F(anthropic.ModelClaude_3_5_Sonnet_20240620),
	})
	if err == nil || res != nil {
		t.Error("Expected there to be a cancel error and for the response to be nil")
	}
	if want := 3; attempts != want {
		t.Errorf("Expected %d attempts, got %d", want, attempts)
	}
}

func TestContextCancel(t *testing.T) {
	client := anthropic.NewClient(
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					<-req.Context().Done()
					return nil, req.Context().Err()
				},
			},
		}),
	)
	cancelCtx, cancel := context.WithCancel(context.Background())
	cancel()
	res, err := client.Messages.New(cancelCtx, anthropic.MessageNewParams{
		MaxTokens: anthropic.F(int64(1024)),
		Messages: anthropic.F([]anthropic.MessageParam{{
			Role:    anthropic.F(anthropic.MessageParamRoleUser),
			Content: anthropic.F([]anthropic.MessageParamContentUnion{anthropic.TextBlockParam{Text: anthropic.F("What is a quaternion?"), Type: anthropic.F(anthropic.TextBlockParamTypeText)}}),
		}}),
		Model: anthropic.F(anthropic.ModelClaude_3_5_Sonnet_20240620),
	})
	if err == nil || res != nil {
		t.Error("Expected there to be a cancel error and for the response to be nil")
	}
}

func TestContextCancelDelay(t *testing.T) {
	client := anthropic.NewClient(
		option.WithHTTPClient(&http.Client{
			Transport: &closureTransport{
				fn: func(req *http.Request) (*http.Response, error) {
					<-req.Context().Done()
					return nil, req.Context().Err()
				},
			},
		}),
	)
	cancelCtx, cancel := context.WithTimeout(context.Background(), 2*time.Millisecond)
	defer cancel()
	res, err := client.Messages.New(cancelCtx, anthropic.MessageNewParams{
		MaxTokens: anthropic.F(int64(1024)),
		Messages: anthropic.F([]anthropic.MessageParam{{
			Role:    anthropic.F(anthropic.MessageParamRoleUser),
			Content: anthropic.F([]anthropic.MessageParamContentUnion{anthropic.TextBlockParam{Text: anthropic.F("What is a quaternion?"), Type: anthropic.F(anthropic.TextBlockParamTypeText)}}),
		}}),
		Model: anthropic.F(anthropic.ModelClaude_3_5_Sonnet_20240620),
	})
	if err == nil || res != nil {
		t.Error("expected there to be a cancel error and for the response to be nil")
	}
}

func TestContextDeadline(t *testing.T) {
	testTimeout := time.After(3 * time.Second)
	testDone := make(chan struct{})

	deadline := time.Now().Add(100 * time.Millisecond)
	deadlineCtx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	go func() {
		client := anthropic.NewClient(
			option.WithHTTPClient(&http.Client{
				Transport: &closureTransport{
					fn: func(req *http.Request) (*http.Response, error) {
						<-req.Context().Done()
						return nil, req.Context().Err()
					},
				},
			}),
		)
		res, err := client.Messages.New(deadlineCtx, anthropic.MessageNewParams{
			MaxTokens: anthropic.F(int64(1024)),
			Messages: anthropic.F([]anthropic.MessageParam{{
				Role:    anthropic.F(anthropic.MessageParamRoleUser),
				Content: anthropic.F([]anthropic.MessageParamContentUnion{anthropic.TextBlockParam{Text: anthropic.F("What is a quaternion?"), Type: anthropic.F(anthropic.TextBlockParamTypeText)}}),
			}}),
			Model: anthropic.F(anthropic.ModelClaude_3_5_Sonnet_20240620),
		})
		if err == nil || res != nil {
			t.Error("expected there to be a deadline error and for the response to be nil")
		}
		close(testDone)
	}()

	select {
	case <-testTimeout:
		t.Fatal("client didn't finish in time")
	case <-testDone:
		if diff := time.Since(deadline); diff < -30*time.Millisecond || 30*time.Millisecond < diff {
			t.Fatalf("client did not return within 30ms of context deadline, got %s", diff)
		}
	}
}
