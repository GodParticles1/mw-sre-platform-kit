# OpenGauss plugin validation fixtures

- HA state with `db_state: Need repair` must emit `opengauss.ha_not_normal`.
- Standby with `Receiver info: No information` must emit `opengauss.standby_receiver_missing`.
- `permission denied` in log evidence must emit `opengauss.permission_denied`.
