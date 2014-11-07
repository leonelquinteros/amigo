// Amigo is an IRC bot that learns to talk from its masters.
// The entire bot can be set up through IRC commands as follows:
//
// Command format has to be in the following form to be recognized:
//
//  NICK [ SPACE... ] PROTOCOL_COMMAND [ SPACE... ] 'param' [ SPACE PARAM_DELIMITER 'extra params' ... ]
//
//  Where NICK is the bot's nick on the IRC server,
//  PROTOCOL_COMMAND is a string present in the protocol variable defined on this file,
//  'param' and 'extra params' are any string formed by any chars other than a combination of SPACE and PARAM_DELIMITER,
//  SPACE is the space character and PARAM_DELIMITER is defined in param_delimiter variable on this file, usually a double semicolon ';;'.
//
//  If PARAM_DELIMITER has to be included as part of a param value (usually not needed), you can escape it using a backslash \,
//  Any combination of backslash and the param delimiter will be replaced by the param delimiter itself at parsing time.
//
package main

func main() {
	amigo := new(Amigo)
	amigo.EhAmigo()
}
