package cmddispatch

import (
    // "log"
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

// func main() {
//     var currentCommand StoredCommand
//     var msg string
//     currentCommand, msg = Command([]byte("MAIL FROM"), []byte("feakuru"), currentCommand)
//     log.Println(msg)
//     log.Println(currentCommand)
//     currentCommand, msg = Command([]byte("RCPT TO"), []byte("feakuru"), currentCommand)
//     log.Println(msg)
//     log.Println(currentCommand)
//     currentCommand, msg = Command([]byte("RCPT TO"), []byte("feakuru"), currentCommand)
//     log.Println(msg)
//     log.Println(currentCommand)
//     currentCommand, msg = Command([]byte("DATA"), []byte(""), currentCommand)
//     log.Println(msg)
//     log.Println(currentCommand)
// }

// MAIL FROM
// RCPT TO (x N)
// DATA
func Command(cmd []byte, arg []byte, previousCommand StoredCommand) (StoredCommand, string) {
    if bytes.Equal(cmd, []byte("MAIL FROM")) {
        if bytes.Equal(previousCommand.cmd, []byte("")) {
            previousCommand.strdSender = string(arg)
            previousCommand.cmd = cmd
            previousCommand.arg = arg
            return previousCommand, "250 OK"
        } else {
            return previousCommand, "500 Error"
        }
    } else if bytes.Equal(cmd, []byte("RCPT TO")) {
        if bytes.Equal(previousCommand.cmd, []byte("MAIL FROM")) || bytes.Equal(previousCommand.cmd, []byte("RCPT TO")) {
            previousCommand.strdRcpts = append(previousCommand.strdRcpts, string(arg))
            previousCommand.cmd = cmd
            previousCommand.arg = arg
            return previousCommand, "250 OK"
        } else {
            return previousCommand, "500 Error"
        }
    } else if bytes.Equal(cmd, []byte("DATA")) {
        previousCommand.cmd = cmd
        previousCommand.arg = arg
        return previousCommand, "354 Send message content; end with <CRLF>.<CRLF>"
    } else {
        return previousCommand, "500 Unknown command"
    }
}
