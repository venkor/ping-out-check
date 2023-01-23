# ping-out-check


## About:

Pinging the other side of wireguard tunnel that connects from home assistant to the main location. If the other side is not reachable, write something to a file which home assistant is watching. Later automation will pick-up that the file has changed and will restart the wireguard client.
