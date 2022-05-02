
# MIME email builder

## Usage example:
```
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
```

## The result is:
```
From: hello@example.com
Message-Id: myid
Reply-To: hello@example.com
Return-Path: bounces@example.com
Subject: See you tomorrow
To: alice@example.com, Bob <bob@example.com>
Content-Type: multipart/alternative; boundary="110000000000863a1705ddeb4f86"
Date: Mon, 02 May 2022 16:38:28 +0200
Mime-Version: 1.0

--110000000000863a1705ddeb4f86
Content-Transfer-Encoding: base64
Content-Type: text/plain; charset=utf-8

U2VlIHlvdSB0b21vcnJvdw==
--110000000000863a1705ddeb4f86
Content-Transfer-Encoding: quoted-printable
Content-Type: text/html; charset=utf-8

<p>See you tomorrow</p>
--110000000000863a1705ddeb4f86--
```

## License

Copyright 2022 Szakszon Péter

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.


