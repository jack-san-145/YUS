include .env
export

#to build and run the yus application
build :
	cd cmd && cd server && go build -o yus . && ./yus

run :
	cd cmd && cd server && ./yus



#to run the application even after close the terminal and store logs in yus.log file
run-forever :
	cd cmd && cd server && nohup ./yus > yus.log 2>&1 &

run-forever-log:
	pgrep yus || (cd cmd/server && nohup ./yus > /var/log/yus/yus.log 2>&1 &)
 
show-log:
	tail -f /var/log/yus/yus.log
	
#to show the cuurently running process of yus application
stop : 
	@pid=$$(ps aux | grep '[y]us' | awk '{print $$2}'); \
	if [ -n "$$pid" ]; then \
		kill $$pid; \
		echo "✅ Process stopped."; \
	else \
		echo "⚠️ No YUS process running."; \
	fi


#to automate the sql migrations
migrate-up:
	migrate -path ./internal/storage/postgres/migrations -database "${POSTGRES_DATABASE_CONNECTION}" up

migrate-down:
	migrate -path ./internal/storage/postgres/migrations -database "${POSTGRES_DATABASE_CONNECTION}" down
drop-schema:
	migrate -path ./internal/storage/postgres/migrations -database "${POSTGRES_DATABASE_CONNECTION}" drop


#to automate the git push
push:
	git add .
	git commit -m "$(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))"
	git push

%:
	@:
