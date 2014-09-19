module IRC
( ircConnect
, ircLogin
, ircJoin
, ircListen
, ircWrite
) where

import Network
import System.IO
import Text.Printf
import Data.List

--
-- Connects to the IRC server
--
ircConnect :: String -> Int -> IO Handle
ircConnect server port = do
    putStrLn ( "Connecting to " ++ server ++ ":" ++ ( show port ) ++ "..." )
    socket <- connectTo server (PortNumber (fromIntegral port))
    hSetBuffering socket NoBuffering
    return socket


--
-- Writes to the network socket
--
ircWrite :: Handle -> String -> IO ()
ircWrite socket msg = do
    hPrintf socket "%s\r\n" msg
    putStrLn ( "> " ++ msg )


--
-- Sends USER and NICK commands to init session on the server
--
ircLogin :: Handle -> String -> IO ()
ircLogin socket nick = do
    ircWrite socket ( "NICK " ++ nick )
    ircWrite socket ( "USER " ++ nick ++ " 0 * :Amigo Bot" )


--
-- Joins a channel
--
ircJoin :: Handle -> String -> IO ()
ircJoin socket channel = ircWrite socket ( "JOIN " ++ channel )


--
-- Loops reading stream input and handles each one through a dispatcher.
--
ircListen :: Handle -> ( Handle -> String -> IO () ) -> IO ()
ircListen socket dispatcher = forever $ do
    msg <- hGetLine socket
    let s = init msg
    dispatcher socket ( clean s )
    putStrLn s
  where
    forever this    = this >> forever this
    clean           = drop 1 . dropWhile (/= ':') . drop 1

