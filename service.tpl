[Unit]
Description=$NAME web remote control
Documentation=https://github.com/HoffmannP/Asasel
After=network.target
StartLimitIntervalSec=0

[Service]
Type=simple
ExecStart=$FULLNAME
RestartSec=1
Restart=always
User=root

[Install]
WantedBy=multi-user.target
