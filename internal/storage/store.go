package storage

import (
	"context"
	"yus/internal/models"
)

type InMemoryStore interface {
	CreateClient(ctx context.Context) error
	GenerateSessionID(ctx context.Context) (string, error)
	DeleteSession(ctx context.Context, sessionID string) error
	AdminExists(ctx context.Context) (bool, error)
	CreateAdminSession(ctx context.Context, adminEmail string) (string, error)
	CheckAdminSession(ctx context.Context, sessionID string) (bool, error)
	AdminLogin(ctx context.Context, email string, password string) (bool, error)
	AddAdmin(ctx context.Context, name string, email string, password string) (string, error)
	CreateDriverSession(ctx context.Context, driverID int) (string, error)
	CheckDriverSession(ctx context.Context, sessionID string) (bool, int, error)
	GetOtp(ctx context.Context, email string) (string, error)
	SetOtp(ctx context.Context, email string, otp string) error
	StoreArrivalStatus(ctx context.Context, driverID int, arrivalStatus map[int]string) error
	GetArrivalStatus(ctx context.Context, driverID int) (map[int]string, error)
	CacheBusRoute(ctx context.Context, route []models.CurrentRoute) error
	GetCachedRoute(ctx context.Context) ([]models.CurrentRoute, error)
	RateLimiter(ctx context.Context, rateLimit *models.RateLimit) (int, error)
}

type DBStore interface {

	// Connection
	Connect(ctx context.Context) error

	// Driver Management
	AddDriver(ctx context.Context, driver *models.Driver) error
	DriverExists(ctx context.Context, driverID int) (bool, error)
	SetDriverPassword(ctx context.Context, driverID int, email string, password string) (bool, error)
	ValidateDriver(ctx context.Context, driverID int, password string) (bool, error)
	GetAvailableDrivers(ctx context.Context) ([]models.AvailableDriver, error)
	DriverExistsInCBR(ctx context.Context, driverID int) (bool, error)
	RemoveDriver(ctx context.Context, driverID int) error
	StoreDriverRemovalRequest(ctx context.Context, driverID int) error
	GetDriverRemovalRequest(ctx context.Context) ([]models.DriverRemovalRequest, error)

	// Bus Management
	AddBus(ctx context.Context, busID int) error
	UpdateBusDriver(ctx context.Context, driverID int, busID int) error
	UpdateBusRoute(ctx context.Context, route *models.BusRoute) error
	RemoveBus(ctx context.Context, busID int) error

	// Route Operation & Validation
	SaveRoute(ctx context.Context, route *models.Route) (string, error)
	InsertRoute(ctx context.Context, route *models.Route) (int, error)
	GetLastRouteID(ctx context.Context) (int, error)
	CheckRouteExists(ctx context.Context, src string, dest string, stops []models.RouteStops) error
	RouteExistsInCBR(ctx context.Context, routeID int) (bool, error)
	ChangeRouteDirection(ctx context.Context, direction string) (bool, error)
	RemoveRoute(ctx context.Context, routeID int) error

	// Driver–Bus–Route Mapping
	AssignDriverToBus(ctx context.Context, driverID int, busID int) error
	AssignRouteToBus(ctx context.Context, routeID int, busID int) error

	// Route Discovery & Querying
	GetAvailableRoutes(ctx context.Context) ([]models.AvilableRoute, error)
	GetAllottedBusForDriver(ctx context.Context, driverID int) (models.AllotedBus, error)
	FindRouteByBusOrDriverID(ctx context.Context, busID int, requestFrom string) (*models.AllRoute, error)
	FindRoutesBySrcDst(ctx context.Context, src string, dest string) ([]models.CurrentRoute, error)
	FindReverseRoutesBySrcDest(ctx context.Context, src string, dest string) ([]models.CurrentRoute, error)
	FindRoutesBySrcDstStop(ctx context.Context, src string, dest string, stop string) ([]models.CurrentRoute, error)
	FindStops(ctx context.Context, route *models.CurrentRoute) error
	GetSrcDestNameByRouteID(ctx context.Context, routeID int) (string, string, string, error)

	// Scheduling Operations
	GetCurrentSchedule(ctx context.Context) ([]models.CurrentSchedule, error)
	GetCurrentBusRoutes(ctx context.Context) ([]models.CurrentRoute, error)
	ScheduleBus(ctx context.Context, schedule *models.CurrentSchedule) error

	// Cache Operations
	CacheRoute(ctx context.Context, route *models.BusRoute) error
	GetCachedRoutesByBusID(ctx context.Context, busID int) ([]models.BusRoute, error)

	// Passenger Operations
	CheckRouteExistsForPassengerWS(ctx context.Context, route models.PassengerWsRequest) (bool, error)
}

type Store struct {
	InMemoryDB InMemoryStore
	DB         DBStore
}
