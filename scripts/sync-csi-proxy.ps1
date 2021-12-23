# helper script, should be called by sync-csi-proxy.sh

# # debug
# Set-PSDebug -Trace 1

# stop the csiproxy service
sc.exe stop csiproxy
Start-Sleep -Seconds 1;
sc.exe delete csiproxy

# copy the binary from the user directory
Copy-Item -Path "C:\Users\$env:UserName\csi-proxy.exe" -Destination "C:\etc\kubernetes\node\bin\csi-proxy.exe"

# restart the csiproxy service
$flags = "-v=5 -windows-service -log_file=C:\etc\kubernetes\logs\csi-proxy.log -logtostderr=false"
sc.exe create csiproxy binPath= "C:\etc\kubernetes\node\bin\csi-proxy.exe $flags"
sc.exe failure csiproxy reset= 0 actions= restart/10000
sc.exe start csiproxy

# start the CSI Proxy before running tests on windows
Start-Sleep -Seconds 5;

Write-Output "Checking the status of csi-proxy"
sc.exe query csiproxy
[System.IO.Directory]::GetFiles("\\.\\pipe\\")

Write-Output "Get logs"
Get-Content C:\etc\kubernetes\logs\csi-proxy.log -Tail 20
