# YUS – Yelloh Bus 🚍
**Real-Time College Bus Tracking Ecosystem**

---

## Introduction

**YUS (Yelloh Bus)** is a real-time college bus tracking ecosystem designed to modernize and digitize traditional campus transportation systems. While the passenger-facing application is known as **Yelloh Bus**, **YUS** represents the complete backend-driven ecosystem that powers all tracking, scheduling, and operational workflows.

The name *Yelloh Bus* is inspired by the iconic yellow-colored college buses, making the product identity intuitive and directly relatable. Unlike conventional single-application solutions, YUS is architected as a **multi-client, real-time distributed system**.

The platform is built with a strong focus on:
- Reliability
- Scalability
- Low latency
- Real-time data consistency

These qualities are critical for transportation systems operating under real-world constraints.

---

## Problem Statement

Traditional college bus systems often suffer from:

- No real-time visibility of bus locations
- Inaccurate or static schedules
- Poor communication between drivers, students, and administrators
- Manual route planning and driver assignment
- High dependency on phone calls and human intermediaries

These issues result in wasted time, uncertainty, overcrowding at bus stops, and inefficient fleet management. **YUS addresses these challenges** by introducing a centralized, real-time tracking and management ecosystem.

---

## Vision and Goals

The core vision of YUS is to:

- Provide real-time bus tracking for students, staff, parents, and the public
- Enable dynamic route, driver, and bus management for administrators
- Reduce operational friction for drivers through driver-friendly design
- Build a scalable, production-grade system capable of handling high concurrency

YUS is designed not merely as an academic project, but as a **deployment-ready platform** for real college environments.

---

## Ecosystem Overview

The YUS ecosystem consists of multiple specialized applications, all connected through a centralized backend server:

1. **YUS Route** – Mobile App (Admin)
2. **YUS Driver** – Mobile App (Driver)
3. **Yelloh Bus** – Mobile App (Passenger)
4. **YUS Admin** – Web Portal (Admin)
5. **YUS Public Website**

All clients communicate exclusively with the backend server, ensuring centralized control, security, and data consistency.

---

## YUS Route – Route Management Application

### Purpose

**YUS Route** is a mobile application designed for administrators to create and define bus routes accurately using real-world data. Instead of drawing routes manually or relying on third-party APIs, administrators physically travel along the route and mark stops using the device’s native GPS.

This approach captures routes exactly as they exist on the road, minimizing human error and assumptions.

### Key Features

- Live GPS-based route creation
- Manual stop marking using real-time physical location
- Accurate latitude and longitude capture
- Automatic duplicate route detection at server level
- Secure submission of route data to the backend

### Workflow

1. Admin opens the YUS Route mobile app
2. Admin physically travels along the bus route
3. Stops are marked manually using current GPS location
4. Route data is sent to the backend server
5. Server validates and stores the route
6. Duplicate routes are detected and ignored safely

---

## YUS Driver – Driver Application

### Purpose

**YUS Driver** is a mobile application exclusively designed for college bus drivers. Its primary responsibility is to stream live GPS data to the backend server in real time.

### Authentication Model

- Drivers cannot self-register
- Drivers are pre-registered by admins
- Identity verification via email and OTP
- Drivers set their own passwords after verification

This controlled onboarding ensures system security.

### Key Features

- Long-lived sessions (no frequent re-login)
- Automatic schedule retrieval on app launch
- WebSocket-based real-time communication
- Start / Stop location sharing
- Background GPS tracking support

### Workflow

1. Driver logs in using driver ID and password
2. Assigned bus and route schedule is fetched
3. WebSocket connection is established
4. Driver starts location sharing
5. GPS coordinates are sent every 5 seconds
6. Tracking continues in background
7. Driver stops sharing when needed

---

## Yelloh Bus – Passenger Application

### Purpose

**Yelloh Bus** is the passenger-facing mobile application used by students, staff, parents, and the public to track buses in real time.

### Key Features

- Track buses by bus number
- Track buses by source and destination
- Support for up-route and down-route
- Live bus movement visualization
- Dynamic ETA calculation
- Custom route maps similar to *Where Is My Train*

### Real-Time Tracking Flow

1. User opens Yelloh Bus
2. WebSocket connection is established
3. Active route schedules are fetched
4. Routes are cached locally
5. User searches by bus number or source/destination
6. Matching route is rendered on the map
7. Live driver location updates move the bus icon

### Dynamic ETA Calculation

- Each stop has a default arrival time
- Actual arrival time is recorded when driver reaches a stop
- ETA for upcoming stops is recalculated dynamically
- Passengers receive realistic arrival predictions

---

## YUS Admin – Administrative Web Portal

### Purpose

**YUS Admin** is a web-based control panel that serves as the command center of the YUS ecosystem.

### Key Responsibilities

- Register and manage drivers
- Add and manage buses
- Create and modify route schedules
- Assign buses to routes
- Assign drivers to buses
- Dynamically update daily schedules
- Remove outdated routes, drivers, and buses

### Scheduling Capabilities

- Multiple routes per bus
- One route selected dynamically per day
- Driver-to-bus allocation can change daily
- Supports both up-route and down-route scheduling

---

## YUS Public Website

### Purpose

The public-facing YUS website provides:

- Overview of the YUS ecosystem
- Application download links
- Feature explanations
- Contact and onboarding information

---

## Conclusion

YUS is a fully integrated, real-time transportation ecosystem built with production-grade architecture. By combining mobile applications, real-time communication, and a scalable backend, YUS delivers an efficient, reliable, and user-friendly solution for modern college transportation systems.
