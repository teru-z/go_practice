package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

type errorReadCloser struct{}

func (errorReadCloser) Read(_ []byte) (int, error) {
	return 0, errors.New("read failed")
}

func (errorReadCloser) Close() error {
	return nil
}

func TestDataProcessor_ProcessesValidInputsAndSkipsInvalid(t *testing.T) {
	in := make(chan []byte, 6)
	out := make(chan Result, 6)

	go DataProcessor(in, out)

	in <- []byte("id-add\n+\n2\n3")
	in <- []byte("id-sub\n-\n10\n3")
	in <- []byte("id-mul\n*\n4\n5")
	in <- []byte("id-div\n/\n20\n4")
	in <- []byte("id-bad-op\n^\n1\n1")
	in <- []byte("id-bad-num\n+\nnot-number\n1")
	close(in)

	got := make([]Result, 0, 4)
	for r := range out {
		got = append(got, r)
	}

	want := []Result{
		{Id: "id-add", Value: 5},
		{Id: "id-sub", Value: 7},
		{Id: "id-mul", Value: 20},
		{Id: "id-div", Value: 5},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected results\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestWriteData_WritesExpectedFormat(t *testing.T) {
	in := make(chan Result, 2)
	in <- Result{Id: "a", Value: 1}
	in <- Result{Id: "b", Value: 42}
	close(in)

	var sb strings.Builder
	WriteData(in, &sb)

	if sb.String() != "a:1\nb:42\n" {
		t.Fatalf("unexpected output: %q", sb.String())
	}
}

func TestNewController_AcceptsAndQueuesData(t *testing.T) {
	out := make(chan []byte, 1)
	h := NewController(out)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("payload"))
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Fatalf("unexpected status code: %d", rr.Code)
	}
	if rr.Body.String() != "OK: 1" {
		t.Fatalf("unexpected response body: %q", rr.Body.String())
	}

	select {
	case got := <-out:
		if string(got) != "payload" {
			t.Fatalf("unexpected queued payload: %q", string(got))
		}
	default:
		t.Fatal("expected payload to be queued, but channel was empty")
	}
}

func TestNewController_Returns503WhenQueueIsFull(t *testing.T) {
	out := make(chan []byte, 1)
	out <- []byte("already-full")
	h := NewController(out)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("new-data"))
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("unexpected status code: %d", rr.Code)
	}
	if rr.Body.String() != "Too Busy: 1" {
		t.Fatalf("unexpected response body: %q", rr.Body.String())
	}
}

func TestNewController_Returns400WhenBodyReadFails(t *testing.T) {
	out := make(chan []byte, 1)
	h := NewController(out)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Body = errorReadCloser{}
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status code: %d", rr.Code)
	}
	if rr.Body.String() != "Bad Input" {
		t.Fatalf("unexpected response body: %q", rr.Body.String())
	}

	select {
	case <-out:
		t.Fatal("did not expect data to be queued on bad input")
	default:
	}
}

// drainWithTimeout は out チャネルが閉じられるまで全結果を収集する。
// タイムアウトした場合はテストを失敗させる。
func drainWithTimeout(t *testing.T, out <-chan Result) []Result {
	t.Helper()
	ch := make(chan []Result, 1)
	go func() {
		var results []Result
		for r := range out {
			results = append(results, r)
		}
		ch <- results
	}()
	select {
	case results := <-ch:
		return results
	case <-time.After(2 * time.Second):
		t.Fatal("DataProcessor が output チャネルを閉じなかった (timeout)")
		return nil
	}
}

// --- DataProcessor の追加テスト ---

func TestDataProcessor_SkipsZeroDivision(t *testing.T) {
	in := make(chan []byte, 2)
	out := make(chan Result, 2)

	go DataProcessor(in, out)

	in <- []byte("id-zero-div\n/\n10\n0")
	in <- []byte("id-add\n+\n1\n2")
	close(in)

	got := drainWithTimeout(t, out)
	want := []Result{{Id: "id-add", Value: 3}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ゼロ除算はスキップされるべき\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestDataProcessor_SkipsMalformedInput(t *testing.T) {
	in := make(chan []byte, 2)
	out := make(chan Result, 2)

	go DataProcessor(in, out)

	in <- []byte("id-short\n+\n1") // 4行必要なのに3行しかない
	in <- []byte("id-add\n+\n1\n2")
	close(in)

	got := drainWithTimeout(t, out)
	want := []Result{{Id: "id-add", Value: 3}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("行数不足の入力はスキップされるべき\nwant: %#v\ngot:  %#v", want, got)
	}
}

// --- WriteData の追加テスト ---

func TestWriteData_EmptyInput(t *testing.T) {
	in := make(chan Result)
	close(in)

	var sb strings.Builder
	WriteData(in, &sb)

	if sb.String() != "" {
		t.Fatalf("入力が空のとき出力も空であるべき、got: %q", sb.String())
	}
}

// --- NewController の追加テスト ---

func TestNewController_CounterNotIncrementedOnBadInput(t *testing.T) {
	out := make(chan []byte, 2)
	h := NewController(out)

	// 1回目: ボディ読み込み失敗 (400)
	req1 := httptest.NewRequest(http.MethodPost, "/", nil)
	req1.Body = errorReadCloser{}
	rr1 := httptest.NewRecorder()
	h.ServeHTTP(rr1, req1)

	// 2回目: 正常
	req2 := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("data"))
	rr2 := httptest.NewRecorder()
	h.ServeHTTP(rr2, req2)

	// 400 はカウントされないので、正常リクエストは "OK: 1" のはず
	if rr2.Body.String() != "OK: 1" {
		t.Fatalf("400 エラーはカウントに含まれるべきでない、got %q", rr2.Body.String())
	}
}

func TestNewController_CounterNotIncrementedOnQueueFull(t *testing.T) {
	out := make(chan []byte, 1)
	h := NewController(out)

	// 1回目: 成功 (キューが埋まる)
	req1 := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("first"))
	rr1 := httptest.NewRecorder()
	h.ServeHTTP(rr1, req1)
	_ = rr1

	// 2回目: キュー満杯 (503)
	req2 := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("second"))
	rr2 := httptest.NewRecorder()
	h.ServeHTTP(rr2, req2)
	_ = rr2

	// キューを空にする
	<-out

	// 3回目: 成功
	req3 := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("third"))
	rr3 := httptest.NewRecorder()
	h.ServeHTTP(rr3, req3)

	// 503 はカウントされないので、2回目の成功は "OK: 2" のはず
	if rr3.Body.String() != "OK: 2" {
		t.Fatalf("503 エラーはカウントに含まれるべきでない、got %q", rr3.Body.String())
	}
}

func TestNewController_RejectsNonPostMethod(t *testing.T) {
	out := make(chan []byte, 1)
	h := NewController(out)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST 以外は 405 を返すべき、got %d", rr.Code)
	}
	select {
	case <-out:
		t.Fatal("GET リクエストでキューに積まれるべきでない")
	default:
	}
}

func TestNewController_CountsAcceptedRequests(t *testing.T) {
	out := make(chan []byte, 2)
	h := NewController(out)

	req1 := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("x"))
	rr1 := httptest.NewRecorder()
	h.ServeHTTP(rr1, req1)

	req2 := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("y"))
	rr2 := httptest.NewRecorder()
	h.ServeHTTP(rr2, req2)

	if rr1.Body.String() != "OK: 1" {
		t.Fatalf("first response body mismatch: %q", rr1.Body.String())
	}
	if rr2.Body.String() != "OK: 2" {
		t.Fatalf("second response body mismatch: %q", rr2.Body.String())
	}

}
