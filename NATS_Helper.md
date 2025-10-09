# NATS NSC tool

You can create credentials using NSC like the below.

```bash
nsc add account claude
nsc edit account claude --sk generate
nsc add user claude -a claude

nsc edit account claude --js-mem-storage 64G --js-disk-storage 100G --js-streams 1000 --js-consumer 10000

nsc push --all
nsc pull --all

nsc generate creds -a claude -n claude -o claude.creds

nats context save claude --nsc nsc://socketglobal/claude/claude 
```
