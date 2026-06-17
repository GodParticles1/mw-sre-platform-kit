# Python environment plugin validation fixtures

- `PROFILE_SOURCE_EXIT=1` must emit `python_env.profile_source_nonzero`.
- non-empty `PROFILE_SOURCE_OUTPUT_BEGIN/END` must emit `python_env.profile_source_output`.
- `IMPORT_FAIL=` must emit `python_env.import_fail`.
