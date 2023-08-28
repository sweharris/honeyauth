## honeyauth

This is a simple GoLang program to capture passwords.  It pretends to be
an SSH daemon and a telnet daemon and when people try to connect it will
present what looks like a normal login session (which will always fail)
and log the username/password used.

It doesn't perfectly emulate the real daemons, so a human attacker could
quickly learn to skip this, but someone using a bot may easily try a
brute force attack.

To try and limit resource consumption we only process a single connection
of each time at a time; the OS TCP backlog may allow multiple connections
to be established, but only one SSH session (for example) will be active.
This attempts to mitigate DoS type activity (don't create too many
threads!).

On Linux to listen to a low "privileged" port you will either need to
run as root (not recommended, since we don't drop privileges) or add
the CAP_NET_BIND_SERVICE capability to the binary

    sudo setcap cap_net_bind_service=+ep ./honeyauth

To use SSH you need to generate host keys.  This can be done with

    make keys

The code will alert you if the keys don't exist

    2023/08/25 10:23:00 SSH: Failed to load private key ssh_host_ecdsa_key. Did you run `make keys`?

To enable the telnet server use the `-t` option, with a port number.
Similarly for ssh, use the `-s` option with port number.

e.g.

    ./honeyauth -s 22 -t 23

or

    ./honeyauth -s 2222 -t 2223

Logs are sent to stdout and are in the following format:

    <DATE> <TIME> <addr>,<protocol>,<method>,<username>,<password>

* The addr is either IP:port or just :port
* The protocol is either TELNET or SSH
* The methid is either Password or Key
* If username or password have a comma in them then they're replaced with %2C

For example:

    % ./honeyauth -t 2223
    2023/08/26 21:03:03 :2223,TELNET,Listening,,
    2023/08/26 21:03:13 127.0.0.1:35912,TELNET,Connection Established,,
    2023/08/26 21:03:16 127.0.0.1:35912,TELNET,Password,hello,there
    2023/08/26 21:03:23 127.0.0.1:35912,TELNET,Password,a%2Cb,c%2Cd
    2023/08/26 21:03:26 127.0.0.1:35912,TELNET,Password,foo,bar

The first line shows we are listening on port 2223 for the TELNET protocol.
The next line shows we've established a TELNET connect from localhost:35912
and the user then tried to login 3 times with hello/there, a,b/c,d and foo/bar.

An SSH session looks similarly:

    2023/08/26 21:05:20 :2222,SSH,Listening,,
    2023/08/26 21:05:31 127.0.0.1:53612,SSH,Connection Established,,
    2023/08/26 21:05:32 127.0.0.1:53612,SSH,Key,user,ssh-rsa
    2023/08/26 21:05:35 127.0.0.1:53612,SSH,Password,user,foo
    2023/08/26 21:05:37 127.0.0.1:53612,SSH,Password,user,a%2Cb
    2023/08/26 21:05:39 127.0.0.1:53612,SSH,Password,user,baz

(Interesting to see `Key` and `ssh-rsa` there; my ssh client attempts an
RSA login because I have ssh-agent configured with a key loaded.)

Why did I write this?  Partially as a learning exercise, partially out of
interest to see what details are being attempted.

Is this safe to use?  Probably not... it depends on if you trust the
golang ssh libraries, I guess!

A potential method to minimize the risk would be to run it in a container
without privileges and use iptable port forwarding to expose the ports.
