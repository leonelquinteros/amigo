--
-- IRC friend
-- Based on http://www.haskell.org/haskellwiki/Roll_your_own_IRC_bot
--

import Network
import System.IO
import Text.Printf
import Data.List
import System.Exit

--
-- Config section. To be replaced by a config file and learning abilities.
--
server = "irc.freenode.org"
port   = 6667
chan   = "#peiiion-bot"
nick   = "peiiionbot"

--
-- Start everything
--
main = do
    irc <- connectTo server (PortNumber (fromIntegral port))
    hSetBuffering irc NoBuffering
    writeTo irc "NICK" nick
    writeTo irc "USER" (nick++" 0 * :Amigo Bot")
    writeTo irc "JOIN" chan
    listenTo irc


--
-- Writes anything to the IRC socket.
-- Takes the command and the params separately, for better formatting.
--
writeTo :: Handle -> String -> String -> IO ()
writeTo irc command params = do
    hPrintf irc "%s %s\r\n" command params
    printf    "> %s %s\n" command params


--
-- Main listenToer loop.
-- Reads over the stream and send everything to the message processor.
-- Handles PING/PONG by itself
--
listenTo :: Handle -> IO ()
listenTo irc = forever $ do
    msg <- hGetLine irc
    let s = init msg
    if ping s then pong s else processMsg irc (clean s)
    putStrLn s
  where
    forever this = this >> forever this

    clean   = drop 1 . dropWhile (/= ':') . drop 1

    ping x  = isPrefixOf "PING :" x
    pong x  = writeTo irc "PONG" (':' : drop 6 x)


--
-- IRC message processor.
-- All IRC messages will go through this to know what to do with them.
-- Checks for commands sent using '@Amigo' prefix.
-- Then tries to talk.
--
processMsg :: Handle -> String -> IO ()
processMsg irc msg
    | isPrefixOf "@Amigo " msg  = runCommand irc (drop 7 msg)
    | otherwise                 = talk irc msg


--
-- Amigo command protocol handler.
--
runCommand :: Handle -> String -> IO()
runCommand irc msg
    | isPrefixOf "go away!" msg = writeTo irc "QUIT" ":Exiting" >> exitWith ExitSuccess
    | isPrefixOf "echo:" msg    = sendMsg irc (drop 5 msg)

--
-- Amigo tries to talk here. Not much yet.
--
talk :: Handle -> String -> IO()
talk irc msg = return () -- ignore everything


--
-- IRC message sender wrapper.
--
sendMsg :: Handle -> String -> IO ()
sendMsg irc msg = writeTo irc "PRIVMSG" (chan ++ " :" ++ msg)
