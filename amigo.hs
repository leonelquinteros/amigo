--
-- IRC friend
-- Based on http://www.haskell.org/haskellwiki/Roll_your_own_IRC_bot
--

import Network
import System.IO
import Text.Printf

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
    write irc "NICK" nick
    write irc "USER" (nick++" 0 * :Amigo")
    write irc "JOIN" chan
    listen irc


--
-- Writes anything to the IRC socket.
-- Takes the command and the params separately, for better formatting.
--
write :: Handle -> String -> String -> IO ()
write irc command params = do
    hPrintf irc "%s %s\r\n" command params
    printf    "> %s %s\n" command params


--
-- Main listener loop.
-- Reads over the stream and send everything to the message processor.
--
listen :: Handle -> IO ()
listen irc = forever $ do
    s <- hGetLine irc
    putStrLn s
  where
    forever this = do
        this
        forever this
