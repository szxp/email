package email

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime/quotedprintable"
	"net/http"
	"strings"
	"time"
)

const DefaultBoundary = "110000000000863a1705ddeb4f86"

func NewEmailBuilder() *EmailBuilder {
	return &EmailBuilder{
		Headers:      make(http.Header),
		PlainHeaders: make(http.Header),
		HTMLHeaders:  make(http.Header),
	}
}

// EmailBuilder helps build a multipart email message
// in MIME (Multipurpose Internet Mail Extensions) encoding.
type EmailBuilder struct {
	// Headers stores the custom key-value pairs of a MIME message.
	Headers http.Header

	// Boundary is the custom boundary.
	// Used only if the Headers does not contain a Content-Type header.
	// If the Headers contains a Content-Type header, the boundary
	// will be parsed from the header value.
	Boundary string

	// PlainHeaders stores the custom key-value pairs for the plain body part.
	PlainHeaders http.Header

	// Plain is the encoded plain text body part in wire format without the trailing \r\n.
	Plain bytes.Buffer

	// HTMLHeaders stores the custom key-value pairs for the HTML body part.
	HTMLHeaders http.Header

	// HTML is the encoded HTML text body part in wire format without the trailing \r\n.
	HTML bytes.Buffer
}

// SetFrom creates the From header.
func (b *EmailBuilder) SetFrom(from string) {
	b.Headers.Set("From", from)
}

// SetTo creates the To header.
func (b *EmailBuilder) SetTo(to []string) {
	buf := &bytes.Buffer{}
	for i, s := range to {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(s)
	}
	b.Headers.Set("To", buf.String())
}

// SetSubject creates the Subject header with the specified s value.
func (b *EmailBuilder) SetSubject(s string) {
	b.Headers.Set("Subject", s)
}

// SetPlainCharset creates the plain text Content-Type header
// with the specified s charset.
func (b *EmailBuilder) SetPlainCharset(s string) {
	b.PlainHeaders.Set("Content-Type", "text/plain; charset="+s)
}

// SetHTMLCharset creates the HTML text Content-Type header
// with the specified s charset.
func (b *EmailBuilder) SetHTMLCharset(s string) {
	b.HTMLHeaders.Set("Content-Type", "text/html; charset="+s)
}

// EncodeBase64Plain encodes s using base64 encoding
// and writes it to Plain buffer.
// Plain buffer will be reset to be empty before encoding,
// but the underlying storage will be retained.
func (b *EmailBuilder) EncodeBase64Plain(s []byte) {
	b.Plain.Reset()
	b.PlainHeaders.Set("Content-Transfer-Encoding", "base64")
	encoder := base64.NewEncoder(base64.StdEncoding, &b.Plain)
	encoder.Write(s)
	encoder.Close()
}

// EncodeBase64HTML encodes s using base64 encoding
// and writes it to HTML buffer.
// HTML buffer will be reset to be empty before encoding,
// but the underlying storage will be retained.
func (b *EmailBuilder) EncodeBase64HTML(s []byte) {
	b.HTML.Reset()
	b.HTMLHeaders.Set("Content-Transfer-Encoding", "base64")
	encoder := base64.NewEncoder(base64.StdEncoding, &b.HTML)
	encoder.Write(s)
	encoder.Close()
}

// EncodeQuotedPlain encodes s using quoted-printable encoding
// and writes it to Plain buffer.
// It limits line length to 76 characters.
// Plain buffer will be reset to be empty before encoding,
// but the underlying storage will be retained.
func (b *EmailBuilder) EncodeQuotedPlain(s []byte) error {
	b.Plain.Reset()
	b.PlainHeaders.Set("Content-Transfer-Encoding", "quoted-printable")
	w := quotedprintable.NewWriter(&b.Plain)
	_, err := w.Write(s)
	if err != nil {
		return err
	}
	return w.Close()
}

// EncodeQuotedHTML encodes s using quoted-printable encoding
// and writes it to HTML buffer.
// It limits line length to 76 characters.
// HTML buffer will be reset to be empty before encoding,
// but the underlying storage will be retained.
func (b *EmailBuilder) EncodeQuotedHTML(s []byte) error {
	b.HTML.Reset()
	b.HTMLHeaders.Set("Content-Transfer-Encoding", "quoted-printable")
	w := quotedprintable.NewWriter(&b.HTML)
	_, err := w.Write(s)
	if err != nil {
		return err
	}
	return w.Close()
}

// Write writes a MIME email in wire format.
func (b *EmailBuilder) Write(w io.Writer) error {
	text := b.Plain.Len() > 0
	html := b.HTML.Len() > 0
	multipart := text && html
	contentType := b.Headers.Get("Content-Type")

	extraHeaders := make(http.Header)
	if b.Headers.Get("MIME-Version") == "" {
		extraHeaders.Set("MIME-Version", "1.0")
	}

	if b.Headers.Get("Date") == "" {
		extraHeaders.Set("Date", time.Now().Format(time.RFC1123Z))
	}

	var boundary string
	if multipart {
		bo, err := b.BoundaryString()
		if err != nil {
			return err
		}
		boundary = bo

		if contentType == "" {
			extraHeaders.Set(
				"Content-Type",
				fmt.Sprintf(`multipart/alternative; boundary="%s"`, boundary),
			)
		}
	}

	err := b.Headers.Write(w)
	if err != nil {
		return err
	}
	err = extraHeaders.Write(w)
	if err != nil {
		return err
	}

	if multipart {
		err = b.writeln(w)
		if err != nil {
			return err
		}

		err = b.writePartPlain(w, boundary, b.PlainHeaders)
		if err != nil {
			return err
		}

		err = b.writePartHTML(w, boundary, b.HTMLHeaders)
		if err != nil {
			return err
		}

		_, err = w.Write([]byte("--" + boundary + "--"))
		return err
	}

	if html {
		return b.writePartHTML(w, "", b.HTMLHeaders)
	}

	return b.writePartPlain(w, "", b.PlainHeaders)
}

func (b *EmailBuilder) writePartPlain(w io.Writer, boundary string, headers http.Header) error {
	extraHeaders := make(http.Header)
	if headers.Get("Content-Type") == "" {
		extraHeaders.Set("Content-Type", "text/plain; charset=utf-8")
	}

	if boundary != "" {
		_, err := w.Write([]byte("--" + boundary))
		if err != nil {
			return err
		}
		err = b.writeln(w)
		if err != nil {
			return err
		}
	}

	err := headers.Write(w)
	if err != nil {
		return err
	}

	err = extraHeaders.Write(w)
	if err != nil {
		return err
	}

	err = b.writeln(w)
	if err != nil {
		return err
	}

	_, err = w.Write(b.Plain.Bytes())
	if err != nil {
		return err
	}

	return b.writeln(w)
}

func (b *EmailBuilder) writePartHTML(w io.Writer, boundary string, headers http.Header) error {
	extraHeaders := make(http.Header)
	if headers.Get("Content-Type") == "" {
		extraHeaders.Set("Content-Type", "text/html; charset=utf-8")
	}

	if boundary != "" {
		_, err := w.Write([]byte("--" + boundary))
		if err != nil {
			return err
		}
		err = b.writeln(w)
		if err != nil {
			return err
		}
	}

	err := headers.Write(w)
	if err != nil {
		return err
	}

	err = extraHeaders.Write(w)
	if err != nil {
		return err
	}

	err = b.writeln(w)
	if err != nil {
		return err
	}

	_, err = w.Write(b.HTML.Bytes())
	if err != nil {
		return err
	}

	return b.writeln(w)
}

func (b *EmailBuilder) writeln(w io.Writer) error {
	_, err := w.Write([]byte("\r\n"))
	return err
}

// BoundaryString returns the boundary string.
// If a custom Content-Type header is specified in the Headers
// the boundary will be parsed from that header value.
// If a custom Boundary field is specified it will return that value.
// Otherwise it will return the DefaultBoundary.
func (b *EmailBuilder) BoundaryString() (string, error) {
	contentType := b.Headers.Get("Content-Type")
	if contentType != "" {
		pattern := `boundary="`
		i := strings.Index(contentType, pattern)
		if i == -1 {
			return "", fmt.Errorf("beginning of the boundary not found in Content-Type header")
		}
		i = i + len(pattern)
		j := strings.Index(contentType[i:], `"`)
		if j == -1 {
			return "", fmt.Errorf("end of the boundary not found in Content-Type header")
		}
		boundary := contentType[i : i+j]
		if boundary == "" {
			return "", fmt.Errorf("empty boundary in Content-Type header")
		}
		return boundary, nil
	}

	boundary := b.Boundary
	if boundary == "" {
		boundary = DefaultBoundary
	}
	return boundary, nil
}
