package server

import "strconv"

func getCommand(client *Client) {
	if len(client.args) != 2 {
		client.ReplyAritryError()
		return
	}
	key, ok := client.server.GetKey(client.args[1])
	if !ok {
		client.ReplyNull()
		return
	}
	switch s := key.(type) {
	default:
		client.ReplyTypeError()
	case string:
		client.ReplyBulk(s)
	}
}

func setCommand(client *Client) {
	if len(client.args) != 3 {
		client.ReplyAritryError()
		return
	}
	client.server.SetKey(client.args[1], client.args[2])
	client.ReplyString("OK")
	client.dirty++
}
func appendCommand(client *Client) {
	if len(client.args) != 3 {
		client.ReplyAritryError()
		return
	}
	key, ok := client.server.GetKey(client.args[1])
	if !ok {
		client.server.SetKey(client.args[1], client.args[2])
		client.ReplyInt(len(client.args[2]))
		client.dirty++
	} else {
		switch s := key.(type) {
		default:
			client.ReplyTypeError()
		case string:
			s += client.args[2]
			client.server.SetKey(client.args[1], s)
			client.ReplyInt(len(s))
			client.dirty++
		}
	}
}

func bitcountCommand(client *Client) {
	var start, end int
	var all bool
	switch len(client.args) {
	default:
		client.ReplyAritryError()
	case 2:
		all = true
	case 4:
		n1, err1 := strconv.ParseInt(client.args[2], 10, 64)
		n2, err2 := strconv.ParseInt(client.args[3], 10, 64)
		if err1 != nil || err2 != nil {
			client.ReplyError("value is not an integer or out of range")
			return
		}
		start, end = int(n1), int(n2)
	}
	key, ok := client.server.GetKey(client.args[1])
	if !ok {
		client.ReplyInt(0)
		return
	}
	switch s := key.(type) {
	default:
		client.ReplyTypeError()
	case string:
		var count int
		if all {
			start, end = 0, len(s)
		} else {
			if start < 0 {
				start = len(s) + start
				if start < 0 {
					start = 0
				}
			}
			if end < 0 {
				end = len(s) + end
				if end < 0 {
					end = 0
				}
			}
		}
		for i := start; i <= end && i < len(s); i++ {
			c := s[i]
			for j := 0; j < 8; j++ {
				count += int((c >> uint(j)) & 0x01)
			}
		}
		client.ReplyInt(count)
	}
}
