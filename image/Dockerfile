ARG  BASE

FROM --platform=linux/amd64 ubuntu:20.04 as linuxbase

FROM ${BASE}

Add /bin/csi-proxy.exe /csi-proxy.exe

ENV PATH="C:\Windows\system32;C:\Windows;C:\WINDOWS\System32\WindowsPowerShell\v1.0\;"
ENTRYPOINT ["csi-proxy.exe", "-v", "4"],