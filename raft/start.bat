start raft.exe localhost:8001 localhost:8002 localhost:8003
start raft.exe localhost:8002 localhost:8003 localhost:8001
raft.exe localhost:8003 localhost:8002 localhost:8001 
pause