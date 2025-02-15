package elevengo

import (
	"context"
	"fmt"

	"github.com/deadblue/elevengo/internal/util"
	"github.com/deadblue/elevengo/lowlevel/api"
	"github.com/deadblue/elevengo/lowlevel/client"
	"github.com/deadblue/elevengo/lowlevel/errors"
)

// DownloadTicket contains all required information to download a file.
type DownloadTicket struct {
	// Download URL.
	Url string
	// Request headers which SHOULD be sent with download URL.
	Headers map[string]string
	// File name.
	FileName string
	// File size in bytes.
	FileSize int64
}

// DownloadCreateTicket creates ticket which contains all required information
// to download a file. Caller can use third-party tools/libraries to download
// file, such as wget/curl/aria2.
func (a *Agent) DownloadCreateTicket(pickcode string, ticket *DownloadTicket) (err error) {
	// Prepare API spec.
	spec := (&api.DownloadSpec{}).Init(pickcode)
	if err = a.llc.CallApi(spec, context.Background()); err != nil {
		return
	}
	// Convert result.
	if len(spec.Result) == 0 {
		return errors.ErrDownloadEmpty
	}
	for _, info := range spec.Result {
		ticket.FileSize, _ = info.FileSize.Int64()
		if ticket.FileSize == 0 {
			return errors.ErrDownloadDirectory
		}
		ticket.FileName = info.FileName
		ticket.Url = info.Url.Url
		ticket.Headers = map[string]string{
			// User-Agent header
			"User-Agent": a.llc.GetUserAgent(),
			// Cookie header
			"Cookie": util.MarshalCookies(a.llc.ExportCookies(ticket.Url)),
		}
		break
	}
	return
}

// Fetch gets content from url using agent underlying HTTP client.
func (a *Agent) Fetch(url string) (body client.Body, err error) {
	return a.llc.Get(url, nil, context.Background())
}

// Range is used in Agent.GetRange().
type Range struct {
	start, end int64
}

func (r *Range) headerValue() string {
	// Generate Range header.
	// Reference: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Range#syntax
	if r.start < 0 {
		return fmt.Sprintf("bytes=%d", r.start)
	} else {
		if r.end < 0 {
			return fmt.Sprintf("bytes=%d-", r.start)
		} else if r.end > r.start {
			return fmt.Sprintf("bytes=%d-%d", r.start, r.end)
		}
	}
	// (r.start >= 0 && r.end <= r.start) is an invalid range
	return ""
}

// RangeFirst makes a Range parameter to request the first `length` bytes.
func RangeFirst(length int64) Range {
	return Range{
		start: 0,
		end:   length - 1,
	}
}

// RangeLast makes a Range parameter to request the last `length` bytes.
func RangeLast(length int64) Range {
	return Range{
		start: 0 - length,
		end:   0,
	}
}

// RangeMiddle makes a Range parameter to request content starts from `offset`,
// and has `length` bytes (at most).
//
// You can pass a negative number in `length`, to request content starts from
// `offset` to the end.
func RangeMiddle(offset, length int64) Range {
	end := offset + length - 1
	if length < 0 {
		end = -1
	}
	return Range{
		start: offset,
		end:   end,
	}
}

// FetchRange gets partial content from |url|, which is located by |rng|.
func (a *Agent) FetchRange(url string, rng Range) (body client.Body, err error) {
	headers := make(map[string]string)
	if value := rng.headerValue(); value != "" {
		headers["Range"] = value
	}
	return a.llc.Get(url, headers, context.Background())
}
