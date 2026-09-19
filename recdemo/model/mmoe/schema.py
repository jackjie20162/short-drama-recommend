MODEL_VERSION = "mmoe.v1"

SPARSE_FEATURES = [
    "user_id",
    "drama_id",
    "region_id",
    "language_id",
    "genre_id",
]

DENSE_FEATURES = [
    "user_watch_seconds",
    "user_completion",
    "user_pay_rate",
    "user_sessions",
    "drama_popularity",
    "drama_completion",
    "drama_pay_rate",
    "drama_age_days",
]

WATCH_BUCKETS = 7
WATCH_BUCKET_MIDPOINTS = [5.0, 20.0, 45.0, 120.0, 240.0, 450.0, 900.0]

# Cardinalities are intentionally configurable at training/export time.
DEFAULT_CARDINALITIES = {
    "user_id": 1_000_000,
    "drama_id": 1_000_000,
    "region_id": 256,
    "language_id": 256,
    "genre_id": 512,
}
