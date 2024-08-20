.PHONY: run_user run_balance run_order run_zookeeper_server run_kafka_server run_channel run_user_kafka run_channel_kafka run_donation_kafka run_balance_kafka run_all stop_all

run_user:
	cmd /c start "User Service" cmd /k "cd user_service && go run cmd/main.go"

run_balance:
	cmd /c start "Balance Service" cmd /k "cd balance_service && go run cmd/main.go"

run_order:
	cmd /c start "Order Service" cmd /k "cd order_service && go run cmd/main.go"

run_donation:
	cmd /c start "Donation Provider Service" cmd /k "cd donation_provider_service && go run cmd/main.go"

run_zookeeper_server:
	cmd /c start "Zookeeper Server" cmd /k "zookeeper-server-start.bat %KAFKA_CONFIG%/zookeeper.properties"

run_kafka_server:
	cmd /c start "Kafka Server" cmd /k "kafka-server-start.bat %KAFKA_CONFIG%/server.properties"

run_channel:
	cmd /c start "Channel Service" cmd /k "cd channel_service && go run cmd/main.go"

run_user_kafka:
	cmd /c start "User Kafka Service" cmd /k "cd user_kafka_service && go run cmd/main.go"

run_channel_kafka:
	cmd /c start "Channel Kafka Service" cmd /k "cd channel_kafka_service && go run cmd/main.go"

run_donation_kafka:
	cmd /c start "Donation Kafka Service" cmd /k "cd donation_kafka_service && go run cmd/main.go"

run_balance_kafka:
	cmd /c start "Balance Kafka Service" cmd /k "cd balance_kafka_service && go run cmd/main.go"

run_orchestration:
	cmd /c start "Orchestration Kafka Service" cmd /k "cd orchestration_service && go run cmd/main.go"

run_all: run_zookeeper_server run_kafka_server run_user run_donation run_balance run_order run_channel run_user_kafka run_channel_kafka run_donation_kafka run_balance_kafka run_orchestration
	@echo "All services are running."

stop_all:
	@taskkill /F /IM "go.exe"
	@taskkill /F /IM "zookeeper-server-start.bat"
	@taskkill /F /IM "kafka-server-start.bat"
	@echo "All services have been stopped."
