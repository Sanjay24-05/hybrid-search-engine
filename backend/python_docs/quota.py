import os

MAX_QUOTA = 50 * 1024 * 1024  # 50MB

def check_quota(user_dir, new_size):
    total = sum(
        os.path.getsize(os.path.join(user_dir, f))
        for f in os.listdir(user_dir)
        if os.path.isfile(os.path.join(user_dir, f))
    )
    return total + new_size <= MAX_QUOTA
