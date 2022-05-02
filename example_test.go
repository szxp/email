package email_test

import (
	"github.com/szxp/email"

	"bytes"
	"fmt"
)

func ExampleEmailBuilder_textOnly() {
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
	// Date: Mon, 02 May 2022 19:51:17 +0200
	// Mime-Version: 1.0
	// Content-Transfer-Encoding: base64
	// Content-Type: text/plain; charset=utf-8
	//
	// U2VlIHlvdSB0b21vcnJvdw==
}

func ExampleEmailBuilder_htmlOnly() {
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
	// Date: Mon, 02 May 2022 19:51:17 +0200
	// Mime-Version: 1.0
	// Content-Transfer-Encoding: quoted-printable
	// Content-Type: text/html; charset=utf-8
	//
	// <p>See you tomorrow</p>
}

func ExampleEmailBuilder_textAndHTML() {
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
