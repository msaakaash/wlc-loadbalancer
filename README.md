
<h1 align="center">Scalable Loadbalancer</h1>

This project is a lightweight HTTP load balancer built using Go, implementing the Weighted Least Connections (WLC) algorithm.
It intelligently distributes incoming client requests across multiple backend servers based on the number of active connections and server-assigned weights.


## 🛠 Tech Stack

- **Go (Golang)**  
- **Standard Libraries**:
  - `net/http`
  - `net/http/httputil` (for reverse proxy)
- **Concurrency Management**:
  - `sync` primitives (Mutex for safe concurrent operations)
- **(Optional)**:
  - **Docker** (for containerization and deployment)


## 👥Development Team
- [`Aakaash M S`](https://github.com/msaakaash)
- [`Venkatesa M G`](https://github.com/VenkatesaMG)
