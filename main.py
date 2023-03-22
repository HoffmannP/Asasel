import daemon  # type: ignore
import server

with daemon.DaemonContext():
    server.main()
