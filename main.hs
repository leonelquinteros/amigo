--
-- IRC friendly bot =)
-- Based on http://www.haskell.org/haskellwiki/Roll_your_own_IRC_bot
--

import Amigo

--
-- Config section. To be replaced by a config file and learning abilities.
--
server  = "irc.freenode.org"
port    = 6667
nick    = "peiiionbot"
channel = "#peiiion-bot"

--
-- Start everything
--
main :: IO ()
main = startAmigo server port nick channel
