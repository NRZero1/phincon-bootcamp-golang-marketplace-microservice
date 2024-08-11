run_user:
	cd user_service && go run cmd/main.go

run_package:
	cd package_service && go run cmd/main.go

run_gateway:
	cd gateway && go run cmd/main.go

run_zookeeper-server:
	zookeeper-server-start.bat %KAFKA_CONFIG%/zookeeper.properties

run_kafka-server:
	kafka-server-start.bat %KAFKA_CONFIG%/server.properties
