package xmpp

import (
	"encoding/xml"
	"strings"
)

const (
	nsClient  = "jabber:client"
	nsStreams = "http://etherx.jabber.org/streams"
	nsSASL    = "urn:ietf:params:xml:ns:xmpp-sasl"
	nsBind    = "urn:ietf:params:xml:ns:xmpp-bind"
	nsSession = "urn:ietf:params:xml:ns:xmpp-session"
	nsRosterV = "urn:xmpp:features:rosterver"
	nsMUCUser = "http://jabber.org/protocol/muc#user"
)

type node struct {
	XMLName  xml.Name
	Attrs    []xml.Attr `xml:",any,attr"`
	Children []node     `xml:",any"`
	Text     string     `xml:",chardata"`
}

func parse(raw string) (node, error) {
	var n node
	err := xml.Unmarshal([]byte(raw), &n)
	return n, err
}

func (n node) attr(name string) string {
	for _, a := range n.Attrs {
		if a.Name.Local == name {
			return a.Value
		}
	}
	return ""
}

func (n node) child(name string) (node, bool) {
	for _, c := range n.Children {
		if c.XMLName.Local == name {
			return c, true
		}
	}
	return node{}, false
}

type element struct {
	name     string
	attrs    [][2]string
	children []element
	content  string
}

func el(name string, attrs ...string) element {
	e := element{name: name}
	for i := 0; i+1 < len(attrs); i += 2 {
		e.attrs = append(e.attrs, [2]string{attrs[i], attrs[i+1]})
	}
	return e
}

func (e element) with(children ...element) element {
	e.children = append(e.children[:len(e.children):len(e.children)], children...)
	return e
}

func (e element) text(content string) element {
	e.content = content
	return e
}

func (e element) String() string {
	var b strings.Builder
	e.write(&b)
	return b.String()
}

func (e element) write(b *strings.Builder) {
	b.WriteByte('<')
	b.WriteString(e.name)
	for _, a := range e.attrs {
		b.WriteByte(' ')
		b.WriteString(a[0])
		b.WriteString(`="`)
		_ = xml.EscapeText(b, []byte(a[1]))
		b.WriteByte('"')
	}
	if len(e.children) == 0 && e.content == "" {
		b.WriteString("/>")
		return
	}
	b.WriteByte('>')
	_ = xml.EscapeText(b, []byte(e.content))
	for _, c := range e.children {
		c.write(b)
	}
	b.WriteString("</")
	b.WriteString(e.name)
	b.WriteByte('>')
}

func streamHeader(domain, id string) string {
	return `<?xml version="1.0"?><stream:stream from="` + domain + `" id="` + id +
		`" version="1.0" xml:lang="en" xmlns:stream="` + nsStreams + `" xmlns="` + nsClient + `">`
}

const streamFooter = "</stream:stream>"

func features(authenticated bool) element {
	f := el("stream:features")
	if authenticated {
		return f.with(el("bind", "xmlns", nsBind), el("session", "xmlns", nsSession), el("ver", "xmlns", nsRosterV))
	}
	return f.with(el("mechanisms", "xmlns", nsSASL).with(el("mechanism").text("PLAIN")), el("ver", "xmlns", nsRosterV))
}

func success() element {
	return el("success", "xmlns", nsSASL)
}

func authFailure() element {
	return el("failure", "xmlns", nsSASL).with(el("not-authorized"))
}

func iqResult(domain, to, id string) element {
	return el("iq", "to", to, "from", domain, "id", id, "xmlns", nsClient, "type", "result")
}

func bindResult(jid, id string) element {
	return el("iq", "to", jid, "id", id, "xmlns", nsClient, "type", "result").
		with(el("bind", "xmlns", nsBind).with(el("jid").text(jid)))
}

func roomPresence(to, from, role, nick, jid string, unavailable bool, codes ...string) element {
	item := el("item", "nick", nick, "jid", jid, "role", role)
	if role != "none" {
		item = el("item", "nick", nick, "jid", jid, "role", role, "affiliation", "none")
	}
	x := el("x", "xmlns", nsMUCUser).with(item)
	for _, code := range codes {
		x = x.with(el("status", "code", code))
	}
	p := el("presence", "to", to, "from", from, "xmlns", nsClient)
	if unavailable {
		p = el("presence", "to", to, "from", from, "xmlns", nsClient, "type", "unavailable")
	}
	return p.with(x)
}

func message(to, from, id, kind, body string) element {
	m := el("message", "to", to, "from", from, "xmlns", nsClient)
	if id != "" {
		m = el("message", "to", to, "from", from, "id", id, "xmlns", nsClient)
	}
	if kind != "" {
		m.attrs = append(m.attrs, [2]string{"type", kind})
	}
	return m.with(el("body").text(body))
}

func selfPresence(jid, status string, away bool) element {
	p := el("presence", "to", jid, "from", jid, "xmlns", nsClient, "type", "available")
	if away {
		p = p.with(el("show").text("away"))
	}
	return p.with(el("status").text(status))
}
