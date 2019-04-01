GOOS=linux go build gosim_engine.go
zip handler.zip ./gosim_engine

GOOS=linux go build gosim_worker.go
zip handler.zip ./gosim_worker ./golang_thing.* ./root-CA.crt
