include .env
export

#to build and run the yus application
build :
	cd cmd && cd server && go build -o yus . && ./yus

run :
	cd cmd && cd server && ./yus



#to run the application even after close the terminal 
run_forever :
	cd cmd && cd server && nohup ./yus > yus.log 2>&1 &
 
#to show the cuurently running process of yus application
stop : 
	ps aux | grep yus



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
	git commit -m "$(filter-out $@,$(MAKECMDGOALS))"
	git push

%:
	@:
