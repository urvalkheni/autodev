import os
from PIL import Image, ImageDraw

def generate_favicon():
    workspace = '/media/heet18/Futuristic/Heet/Github/Autodev-Independent'
    
    # 512x512 canvas for high quality generation
    size = 512
    img = Image.new('RGBA', (size, size), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    # Yellow: #FFD700 (255, 215, 0)
    # Black outline for contrast: #0A0A0A (10, 10, 10)
    yellow_color = (255, 215, 0, 255)
    black_color = (10, 10, 10, 255)
    
    # Mathematically centered lightning bolt coordinates
    vertices = [
        (316, 30),   # Top tip
        (166, 270),  # Middle-left tip
        (256, 270),  # Inner-left bend
        (196, 480),  # Bottom tip
        (346, 240),  # Middle-right tip
        (256, 240)   # Inner-right bend
    ]
    
    # 1. Draw a thick black outline by drawing outline lines and filling
    outline_width = 24
    for i in range(len(vertices)):
        p1 = vertices[i]
        p2 = vertices[(i + 1) % len(vertices)]
        draw.line([p1, p2], fill=black_color, width=outline_width, joint="round")
        
    draw.polygon(vertices, fill=black_color)
    
    # 2. Draw the yellow inner bolt on top
    draw.polygon(vertices, fill=yellow_color)
    
    # Save the output files
    website_dir = os.path.join(workspace, 'apps/website')
    app_ico_path = os.path.join(website_dir, 'app/favicon.ico')
    public_ico_path = os.path.join(website_dir, 'public/favicon.ico')
    public_png_path = os.path.join(website_dir, 'public/favicon.png')
    apple_png_path = os.path.join(website_dir, 'public/apple-touch-icon.png')
    
    # Create target directories if they don't exist
    os.makedirs(os.path.dirname(app_ico_path), exist_ok=True)
    os.makedirs(os.path.dirname(public_ico_path), exist_ok=True)
    
    # Generate the .ico with multiple sizes (16, 32, 48, 256) directly from high-res img
    sizes = [(16, 16), (32, 32), (48, 48), (256, 256)]
    img.save(app_ico_path, format='ICO', sizes=sizes)
    img.save(public_ico_path, format='ICO', sizes=sizes)
    
    # Save standard PNGs
    img.resize((32, 32), Image.Resampling.LANCZOS).save(public_png_path, format='PNG')
    img.resize((180, 180), Image.Resampling.LANCZOS).save(apple_png_path, format='PNG')
    
    print("Favicon generation successful!")
    print("Files written to:")
    print(" -", app_ico_path)
    print(" -", public_ico_path)
    print(" -", public_png_path)
    print(" -", apple_png_path)

if __name__ == '__main__':
    generate_favicon()
