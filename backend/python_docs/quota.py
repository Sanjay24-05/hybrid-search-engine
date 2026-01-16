import os

def check_quota(user_dir, new_size, max_quota_bytes):
    if not os.path.exists(user_dir):
        return True

    total = sum(
        os.path.getsize(os.path.join(user_dir, f))
        for f in os.listdir(user_dir)
        if os.path.isfile(os.path.join(user_dir, f))
    )

    return total + new_size <= max_quota_bytes
