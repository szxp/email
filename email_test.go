package email_test

import (
	"github.com/szxp/email"

	"bytes"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmailBuilder(t *testing.T) {
	cases := []struct {
		Name     string
		From     string
		To       []string
		Subject  string
		Boundary string
		Headers  http.Header

		PlainCharset  string
		PlainEncoding string
		PlainText     string

		HTMLCharset  string
		HTMLEncoding string
		HTMLText     string

		ExpectedMimeVersion string

		ExpectedPlainCharset  string
		ExpectedPlainEncoding string
		ExpectedPlain         string

		ExpectedHTMLCharset  string
		ExpectedHTMLEncoding string
		ExpectedHTML         string

		ExpectedBoundary string
		NotExpected      []string
	}{
		{
			Name:                  "text only base64",
			From:                  "Google Alerts <googlealerts-noreply@example.com>",
			To:                    []string{"alice@example.com", "Bob <bob@eample.com>"},
			Subject:               "Hello",
			Boundary:              "",
			PlainCharset:          "",
			PlainEncoding:         "base64",
			PlainText:             "Hello world",
			HTMLCharset:           "",
			HTMLEncoding:          "",
			HTMLText:              "",
			ExpectedMimeVersion:   "Mime-Version: 1.0",
			ExpectedPlainCharset:  "utf-8",
			ExpectedPlainEncoding: "base64",
			ExpectedPlain:         "SGVsbG8gd29ybGQ=",
			ExpectedHTMLCharset:   "",
			ExpectedHTMLEncoding:  "",
			ExpectedHTML:          "",
			ExpectedBoundary:      "110000000000863a1705ddeb4f86",
			NotExpected:           []string{"text/html"},
		},
		{
			Name:                  "text only quoted printable",
			From:                  "Google Alerts <googlealerts-noreply@example.com>",
			To:                    []string{"alice@example.com", "Bob <bob@eample.com>"},
			Subject:               "Hello",
			Boundary:              "",
			PlainCharset:          "",
			PlainEncoding:         "quoted-printable",
			PlainText:             "Hélló world",
			HTMLCharset:           "",
			HTMLEncoding:          "",
			HTMLText:              "",
			ExpectedMimeVersion:   "Mime-Version: 1.0",
			ExpectedPlainCharset:  "utf-8",
			ExpectedPlainEncoding: "quoted-printable",
			ExpectedPlain:         "H=C3=A9ll=C3=B3 world",
			ExpectedHTMLCharset:   "",
			ExpectedHTMLEncoding:  "",
			ExpectedHTML:          "",
			ExpectedBoundary:      "110000000000863a1705ddeb4f86",
			NotExpected:           []string{"text/html"},
		},
		{
			Name:                  "html only base64",
			From:                  "Google Alerts <googlealerts-noreply@example.com>",
			To:                    []string{"alice@example.com", "Bob <bob@eample.com>"},
			Subject:               "30+ new jobs for 'software engineer'",
			Boundary:              "",
			PlainCharset:          "",
			PlainEncoding:         "",
			PlainText:             "",
			HTMLCharset:           "",
			HTMLEncoding:          "base64",
			HTMLText:              "<p>Hello world</p>",
			ExpectedMimeVersion:   "Mime-Version: 1.0",
			ExpectedPlainCharset:  "",
			ExpectedPlainEncoding: "",
			ExpectedPlain:         "",
			ExpectedHTMLCharset:   "utf-8",
			ExpectedHTMLEncoding:  "base64",
			ExpectedHTML:          "PHA+SGVsbG8gd29ybGQ8L3A+",
			ExpectedBoundary:      "110000000000863a1705ddeb4f86",
			NotExpected:           []string{"text/plain"},
		},
		{
			Name:                  "html only quoted printable",
			From:                  "Google Alerts <googlealerts-noreply@example.com>",
			To:                    []string{"alice@example.com", "Bob <bob@eample.com>"},
			Subject:               "30+ new jobs for 'software engineer'",
			Boundary:              "",
			PlainCharset:          "",
			PlainEncoding:         "",
			PlainText:             "",
			HTMLCharset:           "",
			HTMLEncoding:          "quoted-printable",
			HTMLText:              "<p>Hélló world</p>",
			ExpectedMimeVersion:   "Mime-Version: 1.0",
			ExpectedPlainCharset:  "",
			ExpectedPlainEncoding: "",
			ExpectedPlain:         "",
			ExpectedHTMLCharset:   "utf-8",
			ExpectedHTMLEncoding:  "quoted-printable",
			ExpectedHTML:          "<p>H=C3=A9ll=C3=B3 world</p>",
			ExpectedBoundary:      "110000000000863a1705ddeb4f86",
			NotExpected:           []string{"text/plain"},
		},
		{
			Name:                  "text and html base64",
			From:                  "Google Alerts <googlealerts-noreply@example.com>",
			To:                    []string{"alice@example.com", "Bob <bob@eample.com>"},
			Subject:               "30+ new jobs for 'software engineer'",
			Boundary:              "",
			PlainCharset:          "",
			PlainEncoding:         "base64",
			PlainText:             "plain text message",
			HTMLCharset:           "",
			HTMLEncoding:          "base64",
			HTMLText:              "<p>HTML message</p>",
			ExpectedMimeVersion:   "Mime-Version: 1.0",
			ExpectedPlainCharset:  "utf-8",
			ExpectedPlainEncoding: "base64",
			ExpectedPlain:         "cGxhaW4gdGV4dCBtZXNzYWdl",
			ExpectedHTMLCharset:   "utf-8",
			ExpectedHTMLEncoding:  "base64",
			ExpectedHTML:          "PHA+SFRNTCBtZXNzYWdlPC9wPg==",
			ExpectedBoundary:      "110000000000863a1705ddeb4f86",
		},
		{
			Name:                  "text and html quoted printable",
			From:                  "Google Alerts <googlealerts-noreply@example.com>",
			To:                    []string{"alice@example.com", "Bob <bob@eample.com>"},
			Subject:               "30+ new jobs for 'software engineer'",
			Boundary:              "",
			PlainCharset:          "",
			PlainEncoding:         "quoted-printable",
			PlainText:             "Hélló world",
			HTMLCharset:           "",
			HTMLEncoding:          "quoted-printable",
			HTMLText:              "<p>Hélló world</p>",
			ExpectedMimeVersion:   "Mime-Version: 1.0",
			ExpectedPlainCharset:  "utf-8",
			ExpectedPlainEncoding: "quoted-printable",
			ExpectedPlain:         "H=C3=A9ll=C3=B3 world",
			ExpectedHTMLCharset:   "utf-8",
			ExpectedHTMLEncoding:  "quoted-printable",
			ExpectedHTML:          "<p>H=C3=A9ll=C3=B3 world</p>",
			ExpectedBoundary:      "110000000000863a1705ddeb4f86",
		},
		{
			Name: "custom charset",
			From: "Google Alerts <googlealerts-noreply@example.com>",
			To: []string{
				"alice@example.com",
				"Bob <bob@eample.com>",
				`"Charlie" <charlie@example.com>`,
			},
			Subject:               "30+ new jobs for 'software engineer'",
			Boundary:              "",
			PlainCharset:          "iso-8859-1",
			PlainEncoding:         "base64",
			PlainText:             "plain text message",
			HTMLCharset:           "iso-8859-2",
			HTMLEncoding:          "base64",
			HTMLText:              "<p>HTML message</p>",
			ExpectedMimeVersion:   "Mime-Version: 1.0",
			ExpectedPlainCharset:  "iso-8859-1",
			ExpectedPlainEncoding: "base64",
			ExpectedPlain:         "cGxhaW4gdGV4dCBtZXNzYWdl",
			ExpectedHTMLCharset:   "iso-8859-2",
			ExpectedHTMLEncoding:  "base64",
			ExpectedHTML:          "PHA+SFRNTCBtZXNzYWdlPC9wPg==",
			ExpectedBoundary:      "110000000000863a1705ddeb4f86",
			NotExpected:           []string{"utf-8"},
		},
		{
			Name: "custom content type",
			From: "Google Alerts <googlealerts-noreply@example.com>",
			To: []string{
				"alice@example.com",
				"Bob <bob@eample.com>",
				`"Charlie" <charlie@example.com>`,
			},
			Subject:  "30+ new jobs for 'software engineer'",
			Boundary: "abc123",
			Headers: http.Header{
				"Content-Type": []string{`multipart/alternative; boundary="efg000"`},
			},
			PlainCharset:          "",
			PlainEncoding:         "base64",
			PlainText:             "plain text message",
			HTMLCharset:           "",
			HTMLEncoding:          "base64",
			HTMLText:              "<p>HTML message</p>",
			ExpectedMimeVersion:   "Mime-Version: 1.0",
			ExpectedPlainCharset:  "utf-8",
			ExpectedPlainEncoding: "base64",
			ExpectedPlain:         "cGxhaW4gdGV4dCBtZXNzYWdl",
			ExpectedHTMLCharset:   "utf-8",
			ExpectedHTMLEncoding:  "base64",
			ExpectedHTML:          "PHA+SFRNTCBtZXNzYWdlPC9wPg==",
			ExpectedBoundary:      "efg000",
		},
		{
			Name: "custom boundary",
			From: "Google Alerts <googlealerts-noreply@example.com>",
			To: []string{
				"alice@example.com",
				"Bob <bob@eample.com>",
				`"Charlie" <charlie@example.com>`,
			},
			Subject:               "30+ new jobs for 'software engineer'",
			Boundary:              "abc123",
			PlainCharset:          "",
			PlainEncoding:         "base64",
			PlainText:             "plain text message",
			HTMLCharset:           "",
			HTMLEncoding:          "base64",
			HTMLText:              "<p>HTML message</p>",
			ExpectedMimeVersion:   "Mime-Version: 1.0",
			ExpectedPlainCharset:  "utf-8",
			ExpectedPlainEncoding: "base64",
			ExpectedPlain:         "cGxhaW4gdGV4dCBtZXNzYWdl",
			ExpectedHTMLCharset:   "utf-8",
			ExpectedHTMLEncoding:  "base64",
			ExpectedHTML:          "PHA+SFRNTCBtZXNzYWdlPC9wPg==",
			ExpectedBoundary:      "abc123",
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			b := email.NewEmailBuilder()
			b.SetFrom(c.From)
			b.SetTo(c.To)
			b.SetSubject(c.Subject)
			b.Boundary = c.Boundary

			for k, v := range c.Headers {
				b.Headers.Set(k, v[0])
			}

			if c.PlainCharset != "" {
				b.SetPlainCharset(c.PlainCharset)
			}
			switch c.PlainEncoding {
			case "base64":
				b.EncodeBase64Plain([]byte(c.PlainText))
			case "quoted-printable":
				b.EncodeQuotedPlain([]byte(c.PlainText))
			}

			if c.HTMLCharset != "" {
				b.SetHTMLCharset(c.HTMLCharset)
			}
			switch c.HTMLEncoding {
			case "base64":
				b.EncodeBase64HTML([]byte(c.HTMLText))
			case "quoted-printable":
				b.EncodeQuotedHTML([]byte(c.HTMLText))
			}

			w := &bytes.Buffer{}
			err := b.Write(w)
			assert.NoError(t, err)

			msg := w.String()
			assert.Contains(t, msg, c.ExpectedMimeVersion+"\r\n")
			assert.Contains(t, msg, "Date: ")
			assert.Contains(t, msg, c.From+"\r\n")
			assert.Contains(t, msg, strings.Join(c.To, ", ")+"\r\n")
			assert.Contains(t, msg, "Subject: "+c.Subject+"\r\n")

			if c.ExpectedPlainCharset != "" {
				assert.Contains(t, msg, "Content-Type: text/plain; charset="+c.ExpectedPlainCharset+"\r\n")
			}
			if c.ExpectedPlainEncoding != "" {
				assert.Contains(t, msg, "Content-Transfer-Encoding: "+c.ExpectedPlainEncoding+"\r\n")
			}
			if c.ExpectedPlain != "" {
				assert.Contains(t, msg, c.ExpectedPlain+"\r\n")
			}

			if c.ExpectedHTMLCharset != "" {
				assert.Contains(t, msg, "Content-Type: text/html; charset="+c.ExpectedHTMLCharset+"\r\n")
			}
			if c.ExpectedHTMLEncoding != "" {
				assert.Contains(t, msg, "Content-Transfer-Encoding: "+c.ExpectedHTMLEncoding+"\r\n")
			}
			if c.ExpectedHTML != "" {
				assert.Contains(t, msg, c.ExpectedHTML+"\r\n")
			}

			boundary, err := b.BoundaryString()
			assert.NoError(t, err)
			assert.Equal(t, c.ExpectedBoundary, boundary)

			assert.Contains(t, msg, `boundary="`+c.ExpectedBoundary+`"`)
			assert.Contains(t, msg, "--"+c.ExpectedBoundary)
			assert.Contains(t, msg, "--"+c.ExpectedBoundary+"--")

			if len(c.NotExpected) > 0 {
				for _, x := range c.NotExpected {
					assert.NotContains(t, msg, x)
				}
			}
		})
	}
}

func ExampleEmailBuilder() {
	b := email.NewEmailBuilder()
	b.SetFrom("hello@example.com")
	b.SetTo([]string{
		"alice@example.com",
		"Bob <bob@example.com>",
	})
	b.SetSubject("See you tomorrow")

	// Add custom headers
	b.Headers.Set("Reply-To", "hello@example.com")
	b.Headers.Set("Return-Path", "bounces@example.com")
	b.Headers.Set("Message-ID", "myid")

	b.SetPlainCharset("utf-8")
	b.EncodeBase64Plain([]byte("See you tomorrow"))
	b.SetHTMLCharset("utf-8")
	b.EncodeQuotedHTML([]byte("<p>See you tomorrow</p>"))

	w := &bytes.Buffer{}
	err := b.Write(w)
	if err != nil {
		// handle error
	}
	msg := w.String()
	fmt.Println(msg)

	// This is the result:

    // From: hello@example.com
    // Message-Id: myid
    // Reply-To: hello@example.com
    // Return-Path: bounces@example.com
    // Subject: See you tomorrow
    // To: alice@example.com, Bob <bob@example.com>
    // Content-Type: multipart/alternative; boundary="110000000000863a1705ddeb4f86"
    // Date: Mon, 02 May 2022 16:38:28 +0200
    // Mime-Version: 1.0
    // 
    // --110000000000863a1705ddeb4f86
    // Content-Transfer-Encoding: base64
    // Content-Type: text/plain; charset=utf-8
    // 
    // U2VlIHlvdSB0b21vcnJvdw==
    // --110000000000863a1705ddeb4f86
    // Content-Transfer-Encoding: quoted-printable
    // Content-Type: text/html; charset=utf-8
    // 
    // <p>See you tomorrow</p>
    // --110000000000863a1705ddeb4f86--
}
