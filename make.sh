go build -o release/dBacker \
   -ldflags "-w -s \
   -X main.Version=0.1.1 \
   -X main.BuildNumber=29092021 \
   -X main.Commit=latest \
   -X main.BuildTime=${date}" \
  ./cmd/dBacker