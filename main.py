import daemon
import server

with daemon.DaemonContext():
    server.main()
