# YUS (Yelloh Bus) 🚍  
*A real-time college bus tracking system*

---

## 📖 Description
**YUS (Yelloh Bus)** is a scalable real-time bus tracking application built to simplify college transportation.  
It allows **drivers** to share live GPS location, **passengers** to track buses on their phones, and **admins** to manage routes and schedules.  

The system is designed using **Go** for high-performance backend services, **PostgreSQL** for persistent data storage, and **Redis** for caching and handling high-concurrency real-time updates.  
Frontend apps (Driver, Passenger, Admin) are built with **React Native**, providing a seamless cross-platform experience.  

---

## ✨ Features
- **Driver App**
  - Shares live GPS coordinates with the server in real-time.  

- **Passenger App**
  - Tracks bus location live using WebSockets.  
  - Displays bus movement on a map (similar to "Where is my train").  

- **Admin App**
  - Define and manage bus routes and stops.  
  - Monitor active buses and schedules.  

- **Backend (Go)**
  - High-performance APIs and WebSocket communication.  
  - PostgreSQL for reliable data storage.  
  - Redis for caching live locations and reducing server load.  

---
