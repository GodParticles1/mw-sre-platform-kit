# Redis plugin validation fixtures

- `role:slave + master_link_status:down` must emit `redis.slave_link_down`.
- connection refused evidence must emit `redis.port_down`.
- more than one `role:master` evidence block without slave evidence must emit `redis.both_master`.
