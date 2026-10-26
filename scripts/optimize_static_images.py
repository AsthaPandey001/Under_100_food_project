import os
from PIL import Image

STATIC_DIR = "static"

jpg_files = [f for f in os.listdir(STATIC_DIR) if f.lower().endswith(".jpg")]

print(f"Found {len(jpg_files)} JPG files in {STATIC_DIR}")

total_orig_size = 0
total_webp_size = 0

for filename in sorted(jpg_files):
    src_path = os.path.join(STATIC_DIR, filename)
    base_name = os.path.splitext(filename)[0]
    webp_path = os.path.join(STATIC_DIR, f"{base_name}.webp")
    
    orig_size = os.path.getsize(src_path)
    total_orig_size += orig_size
    
    with Image.open(src_path) as im:
        # Convert RGB if needed
        if im.mode != "RGB":
            im = im.convert("RGB")
        # Save as optimized WebP
        im.save(webp_path, "WEBP", quality=82, method=6)
        
    webp_size = os.path.getsize(webp_path)
    total_webp_size += webp_size
    ratio = (1 - webp_size / orig_size) * 100
    print(f"{filename:22} {orig_size/1024:6.1f}KB -> {webp_size/1024:6.1f}KB ({ratio:.1f}% reduction)")

print("-" * 55)
print(f"Total size: {total_orig_size/1024/1024:.2f}MB -> {total_webp_size/1024/1024:.2f}MB ({(1 - total_webp_size/total_orig_size)*100:.1f}% reduction)")
