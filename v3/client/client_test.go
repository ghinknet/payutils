package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	stderrors "errors"

	"go.gh.ink/payutils/v3/errors"
	"go.gh.ink/payutils/v3/model"
)

// fakeParam is a driver-specific creation param used only in tests.
type fakeParam struct{ id string }

// otherParam is a param no driver claims.
type otherParam struct{}

// fakePayClient is a controllable model.PayClient stub.
type fakePayClient struct {
	claims       bool // whether Create claims fakeParam
	createResult any
	createErr    error
	createCalled bool

	callbackCalled bool
	callbackBody   string
}

func (f *fakePayClient) Create(_ context.Context, param any) (bool, any, error) {
	if _, ok := param.(fakeParam); !ok || !f.claims {
		return false, nil, nil
	}
	f.createCalled = true
	return true, f.createResult, f.createErr
}

func (f *fakePayClient) Callback(w http.ResponseWriter, _ *http.Request) {
	f.callbackCalled = true
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(f.callbackBody))
}

func (f *fakePayClient) Status(string) (model.ReturnStatus, error) { return model.ReturnStatus{}, nil }
func (f *fakePayClient) Close(string) error                        { return nil }
func (f *fakePayClient) Refund(string, string, string, int64, int64, string) error {
	return nil
}

func TestCreate_DispatchesToClaimingDriver(t *testing.T) {
	claiming := &fakePayClient{claims: true, createResult: map[string]string{"payUrl": "u"}}
	notClaiming := &fakePayClient{claims: false}
	c := &Client{PayClient: map[string]model.PayClient{
		"a": notClaiming,
		"b": claiming,
	}}

	res, err := c.Create(context.Background(), fakeParam{id: "x"})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if !claiming.createCalled {
		t.Error("claiming driver was not invoked")
	}
	m, ok := res.(map[string]string)
	if !ok || m["payUrl"] != "u" {
		t.Errorf("result = %v, want map with payUrl=u", res)
	}
}

func TestCreate_NoDriverClaims(t *testing.T) {
	c := &Client{PayClient: map[string]model.PayClient{
		"a": &fakePayClient{claims: true}, // claims fakeParam, but we pass otherParam
	}}
	_, err := c.Create(context.Background(), otherParam{})
	if !stderrors.Is(err, errors.ErrNoDriverClaimed) {
		t.Errorf("err = %v, want ErrNoDriverClaimed", err)
	}
}

func TestCreate_PropagatesDriverError(t *testing.T) {
	driverErr := errors.New("driver failed")
	c := &Client{PayClient: map[string]model.PayClient{
		"a": &fakePayClient{claims: true, createErr: driverErr},
	}}
	_, err := c.Create(context.Background(), fakeParam{})
	if !stderrors.Is(err, driverErr) {
		t.Errorf("err = %v, want driverErr", err)
	}
}

func TestCallback_ForwardsToDriver(t *testing.T) {
	fake := &fakePayClient{callbackBody: "success"}
	c := &Client{PayClient: map[string]model.PayClient{"alipay": fake}}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/alipay/callback", nil)

	if err := c.Callback("alipay", rec, req); err != nil {
		t.Fatalf("Callback returned error: %v", err)
	}
	if !fake.callbackCalled {
		t.Error("driver Callback was not invoked")
	}
	if rec.Code != http.StatusOK || rec.Body.String() != "success" {
		t.Errorf("response = %d %q, want 200 success", rec.Code, rec.Body.String())
	}
}

func TestCallback_UnknownUpstream(t *testing.T) {
	c := &Client{PayClient: map[string]model.PayClient{}}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/x/callback", nil)

	err := c.Callback("nope", rec, req)
	if !stderrors.Is(err, errors.ErrUpstreamNotFound) {
		t.Errorf("err = %v, want ErrUpstreamNotFound", err)
	}
}
