module Amigo
( startAmigo
) where

import Network
import System.IO
import Data.List
import System.Exit

import IRC

--
-- Bot init. Connects to an IRC server and starts listening.
--
startAmigo :: String -> Int -> String -> String -> IO ()
startAmigo server port nick channel = do
    socket <- ircConnect server port
    ircLogin socket nick
    ircJoin socket channel
    ircListen socket processMsg


--
-- Dispatches messages from the IRC server
--
processMsg :: Handle -> String -> IO ()
processMsg socket msg
    | isPrefixOf "PING :" msg   = ircWrite socket ( "PONG :" ++ ( drop 6 msg ) )
    | isPrefixOf "@amigo " msg  = runCommand socket (drop 7 msg)
    | otherwise                 = talk socket msg


--
-- Amigo command protocol handler.
--
runCommand :: Handle -> String -> IO()
runCommand socket cmd
    | isPrefixOf "go away!" cmd = ircWrite socket "QUIT :Exiting" >> exitWith ExitSuccess
    | otherwise = return () -- ignore everything else


--
-- Amigo tries to talk here. Nothing yet
--
talk :: Handle -> String -> IO()
talk irc msg = return () -- ignore everything
