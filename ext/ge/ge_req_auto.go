package ge

import x "github.com/linuxdeepin/go-x11-client"

func QueryVersion(conn *x.Conn, majorVersion, minorVersion uint16) QueryVersionCookie {
	body := encodeQueryVersion(majorVersion, minorVersion)
	req := &x.ProtocolRequest{
		Ext: _ext,
		Header: x.RequestHeader{
			Data: QueryVersionOpcode,
		},
		Body: body,
	}
	seq := conn.SendRequest(x.RequestChecked, req)
	return QueryVersionCookie(seq)
}

func (cookie QueryVersionCookie) Reply(conn *x.Conn) (*QueryVersionReply, error) {
	replyBuf, err := conn.WaitForReply(uint64(cookie))
	if err != nil {
		return nil, err
	}
	r := x.NewReaderFromData(replyBuf)
	var reply QueryVersionReply
	err = readQueryVersionReply(r, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}
