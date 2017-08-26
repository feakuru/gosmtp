package cmddispatch

import (
    "bytes"
)

type StoredCommand struct {
    cmd []byte
    arg []byte
    strdSender string
    strdRcpts []string
    strdTxt string
    err error
}

// MAIL FROM
// RCPT TO (x N)
// DATA
func Command(cmd []byte, arg []byte, previousCommand StoredCommand) (StoredCommand, string) {
    var currentCommand StoredCommand
    if bytes.Equal(cmd, []byte("MAIL FROM")) {
        if bytes.Equal(previousCommand.cmd, []byte("")) {
            currentCommand.strdSender = string(arg)
            currentCommand.cmd = cmd
            currentCommand.arg = arg
            return currentCommand, "250 OK\r\n"
        } else {
            return currentCommand, "500 Error\r\n"
        }
    } else if bytes.Equal(cmd, []byte("RCPT TO")) {
        if bytes.Equal(previousCommand.cmd, []byte("MAIL FROM")) || bytes.Equal(previousCommand.cmd, []byte("RCPT TO")) {
            if (previousCommand.strdRcpts == nil) {
                currentCommand.strdRcpts = make([]string, 10)
            } else {
                currentCommand.strdRcpts = previousCommand.strdRcpts
            }
            currentCommand.strdRcpts = append(currentCommand.strdRcpts, string(arg))
            currentCommand.cmd = cmd
            currentCommand.arg = arg
            return currentCommand, "250 OK\r\n"
        } else {
            return previousCommand, "500 Error\r\n"
        }
    } else if bytes.Equal(cmd, []byte("DATA")) {
        currentCommand.cmd = cmd
        currentCommand.arg = arg
        return currentCommand, "354 Send message content; end with <CRLF>.<CRLF>\r\n"
    } else {
        return currentCommand, "500 Unknown command\r\n"
    }
}
