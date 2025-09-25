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
show_yus_process : 
	ps aux | grep yus



#to automate the sql migrations
migrate_up:
	migrate -path ./internal/storage/postgres/migrations -database "${POSTGRES_DATABASE_CONNECTION}" up

migrate_down:
	migrate -path ./internal/storage/postgres/migrations -database "${POSTGRES_DATABASE_CONNECTION}" down

