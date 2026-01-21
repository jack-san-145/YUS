package postgres

import (
	"context"
	"log"
	"yus/internal/models"
)

func (pg *PgStore) GetBackupRoutes(ctx context.Context) ([]models.BackupRoute, error) {

	var Routes []models.BackupRoute
	var path_map = make(map[int]bool)
	query := "select route_id,path,direction,route_json,created_at from backup_routes order by direction desc"

	rows, err := pg.Pool.Query(ctx, query)
	if err != nil {
		log.Println("error while get backup routes - ", err)
	}

	defer rows.Close()

	for rows.Next() {
		var common_route models.Route
		var route models.BackupRoute
		err := rows.Scan(&route.ID, &route.Path, &route.Direction, &common_route, &route.CreatedAt)
		if err != nil {
			log.Println("error while scanning backuproutes - ", err)
			return nil, err
		}
		if route.Path == "SAME" {
			route.UpRoute = common_route
			Routes = append(Routes, route)
		} else if route.Path == "DIFFERENT" {
			if path_map[route.ID] {
				for i, r := range Routes {
					if r.ID == route.ID {
						Routes[i].DownRoute = common_route
						delete(path_map, route.ID)
						break
					}
				}
			} else {
				route.UpRoute = common_route
				path_map[route.ID] = true
				Routes = append(Routes, route)
			}

		}

	}
	return Routes, err
}

func (pg *PgStore) StoreToBackupRoute(ctx context.Context, path string, route *models.Route) error {
	var err error
	query := `insert into backup_routes(route_id,route_name,path,src,dest,direction,route_json) 
					values($1,$2,$3,$4,$5,$6,$7)`

	if path == "SAME" {
		_, err = pg.Pool.Exec(ctx, query,
			route.Id,
			route.UpRouteName,
			path,
			route.Src,
			route.Dest,
			route.Direction,
			*route)
	} else if path == "DIFFERENT" {
		if route.Direction == "UP" {
			_, err = pg.Pool.Exec(ctx, query,
				route.Id,
				route.UpRouteName,
				path,
				route.Src,
				route.Dest,
				route.Direction,
				*route)
		} else if route.Direction == "DOWN" {
			_, err = pg.Pool.Exec(ctx, query,
				route.Id,
				route.DownRouteName,
				path,
				route.Src,
				route.Dest,
				route.Direction,
				*route)
		}
	}
	return err

}
