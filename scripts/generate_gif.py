import os
import re
from PIL import Image, ImageDraw, ImageFont

# Set up canvas dimensions
WIDTH, HEIGHT = 960, 600
TERM_WIDTH, TERM_HEIGHT = 840, 480
TERM_X = (WIDTH - TERM_WIDTH) // 2
TERM_Y = (HEIGHT - TERM_HEIGHT) // 2

# Font setup
MONO_FONT_PATH = "/usr/share/fonts/truetype/freefont/FreeMono.ttf"
MONO_BOLD_FONT_PATH = "/usr/share/fonts/truetype/freefont/FreeMonoBold.ttf"
OUTFIT_FONT_PATH = "/usr/share/fonts/truetype/roboto-slab/RobotoSlab-Regular.ttf"
OUTFIT_BOLD_FONT_PATH = "/usr/share/fonts/truetype/roboto-slab/RobotoSlab-Bold.ttf"

font_size = 14
font = ImageFont.truetype(MONO_FONT_PATH, font_size)
font_bold = ImageFont.truetype(MONO_BOLD_FONT_PATH, font_size)

# Calculate character dimensions
bbox = font.getbbox("A")
CHAR_WIDTH = bbox[2] - bbox[0]
CHAR_HEIGHT = bbox[3] - bbox[1]
LINE_HEIGHT = CHAR_HEIGHT + 6

# Color scheme definitions
COLORS = {
    "default": (226, 232, 240), # slate-200
    "prompt": (56, 189, 248),  # sky-400
    "cmd": (255, 255, 255),     # white
    "success": (16, 185, 129), # emerald-500
    "warn": (251, 191, 36),    # amber-400
    "error": (239, 68, 68),    # red-500
    "info": (129, 140, 248),   # indigo-400
    "gray": (100, 116, 139),   # slate-500
    "logo": (129, 140, 248),   # indigo-400
    "magenta": (244, 114, 182) # pink-400
}

# Parse a line containing formatting tags like [green], [bold], [magenta bold], etc.
# Returns a list of (text, color, is_bold) tuple tokens
def parse_formatted_line(line):
    # Regex to find tags like [green], [bold magenta], etc.
    pattern = re.compile(r'\[([^\]]+)\]')
    tokens = []
    
    parts = pattern.split(line)
    
    current_color = COLORS["default"]
    current_bold = False
    
    # Slicing parts: even indices are text, odd indices are tags
    for i, part in enumerate(parts):
        if i % 2 == 1:
            # It's a tag
            tag_parts = part.lower().split()
            current_bold = "bold" in tag_parts
            color_name = next((tp for tp in tag_parts if tp != "bold"), "default")
            current_color = COLORS.get(color_name, COLORS["default"])
        else:
            # It's text
            if part:
                tokens.append((part, current_color, current_bold))
                
    return tokens

# Draw the background with a beautiful radial gradient
def draw_gradient_background():
    img = Image.new("RGB", (WIDTH, HEIGHT))
    draw = ImageDraw.Draw(img)
    
    # Draw radial gradient simulation using expanding rectangles
    for r in range(max(WIDTH, HEIGHT), 0, -4):
        # Calculate color interpolation
        factor = r / max(WIDTH, HEIGHT)
        # Gradient from deep dark blue/indigo (10, 10, 28) to almost black (5, 5, 8)
        color = (
            int(10 * factor + 5 * (1 - factor)),
            int(10 * factor + 5 * (1 - factor)),
            int(28 * factor + 8 * (1 - factor))
        )
        
        left = (WIDTH - r) // 2
        top = (HEIGHT - r) // 2
        right = left + r
        bottom = top + r
        draw.rectangle([left, top, right, bottom], fill=color)
        
    return img

def render_frame(buffer, active_cmd="", show_cursor=True, show_overlay=False):
    # Get base background
    img = draw_gradient_background()
    draw = ImageDraw.Draw(img)
    
    # 1. Draw Terminal Glassmorphic Container Shadow & Border
    # Draw simple drop shadow
    draw.rounded_rectangle(
        [TERM_X + 4, TERM_Y + 4, TERM_X + TERM_WIDTH + 4, TERM_Y + TERM_HEIGHT + 4],
        radius=12, fill=(0, 0, 0, 128)
    )
    # Draw terminal background
    draw.rounded_rectangle(
        [TERM_X, TERM_Y, TERM_X + TERM_WIDTH, TERM_Y + TERM_HEIGHT],
        radius=12, fill=(13, 13, 23)
    )
    # Draw border
    draw.rounded_rectangle(
        [TERM_X, TERM_Y, TERM_X + TERM_WIDTH, TERM_Y + TERM_HEIGHT],
        radius=12, outline=(99, 102, 241), width=1
    )
    
    # 2. Draw Terminal Header
    # Header background
    draw.rounded_rectangle(
        [TERM_X + 1, TERM_Y + 1, TERM_X + TERM_WIDTH - 1, TERM_Y + 40],
        radius=0, fill=(20, 20, 35)
    )
    # Re-draw top corners of header to match main container
    draw.pieslice([TERM_X, TERM_Y, TERM_X + 24, TERM_Y + 24], 180, 270, fill=(20, 20, 35))
    draw.pieslice([TERM_X + TERM_WIDTH - 24, TERM_Y, TERM_X + TERM_WIDTH, TERM_Y + 24], 270, 360, fill=(20, 20, 35))
    draw.rectangle([TERM_X + 12, TERM_Y + 1, TERM_X + TERM_WIDTH - 12, TERM_Y + 20], fill=(20, 20, 35))
    
    # Traffic lights
    draw.ellipse([TERM_X + 16, TERM_Y + 14, TERM_X + 28, TERM_Y + 26], fill=(239, 68, 68)) # red
    draw.ellipse([TERM_X + 36, TERM_Y + 14, TERM_X + 48, TERM_Y + 26], fill=(234, 179, 8)) # yellow
    draw.ellipse([TERM_X + 56, TERM_Y + 14, TERM_X + 68, TERM_Y + 26], fill=(34, 197, 94)) # green
    
    # Header title
    title_text = "heet@pop-os: ~/projects/autodev"
    draw.text((TERM_X + (TERM_WIDTH - len(title_text) * 8) // 2, TERM_Y + 13), title_text, fill=(100, 116, 139), font=font)
    
    # Header badge
    badge_text = "v0.3.2"
    draw.rounded_rectangle([TERM_X + TERM_WIDTH - 70, TERM_Y + 10, TERM_X + TERM_WIDTH - 16, TERM_Y + 30], radius=4, fill=(99, 102, 241, 40), outline=(99, 102, 241), width=1)
    draw.text((TERM_X + TERM_WIDTH - 60, TERM_Y + 13), badge_text, fill=(165, 180, 252), font=font)
    
    # 3. Draw Terminal Body Text
    text_area_y = TERM_Y + 50
    text_area_height = TERM_HEIGHT - 60
    max_visible_lines = text_area_height // LINE_HEIGHT
    
    # Build complete lines buffer including the active prompt
    display_buffer = list(buffer)
    if active_cmd is not None:
        prompt_line = f"[prompt]heet@pop-os:~/projects/autodev$ [cmd]{active_cmd}"
        display_buffer.append(prompt_line)
        
    # Scroll text if it exceeds maximum visible lines
    if len(display_buffer) > max_visible_lines:
        display_buffer = display_buffer[-max_visible_lines:]
        
    current_y = text_area_y
    for raw_line in display_buffer:
        tokens = parse_formatted_line(raw_line)
        current_x = TERM_X + 20
        
        for text, color, is_bold in tokens:
            fnt = font_bold if is_bold else font
            draw.text((current_x, current_y), text, fill=color, font=fnt)
            # Since it's monospaced, X advances by character length * CHAR_WIDTH
            current_x += len(text) * CHAR_WIDTH
            
        # Draw cursor block at the end of the last line if active
        if raw_line == display_buffer[-1] and active_cmd is not None and show_cursor:
            draw.rectangle([current_x + 2, current_y + 2, current_x + 10, current_y + LINE_HEIGHT - 4], fill=(255, 255, 255))
            
        current_y += LINE_HEIGHT
        
    # 4. Draw Overlay Metrics Dashboard if show_overlay is True
    if show_overlay:
        # Semi-transparent overlay background
        overlay = Image.new("RGBA", (TERM_WIDTH - 2, TERM_HEIGHT - 42), (10, 10, 22, 245))
        img.paste(overlay, (TERM_X + 1, TERM_Y + 41), overlay)
        
        # Load fonts for dashboard
        dash_font_title = ImageFont.truetype(OUTFIT_BOLD_FONT_PATH, 28)
        dash_font_sub = ImageFont.truetype(OUTFIT_FONT_PATH, 14)
        dash_font_label = ImageFont.truetype(OUTFIT_FONT_PATH, 12)
        dash_font_val = ImageFont.truetype(OUTFIT_BOLD_FONT_PATH, 20)
        
        # Draw badge icon
        draw.text((TERM_X + TERM_WIDTH // 2 - 16, TERM_Y + 80), "⚡", fill=(129, 140, 248), font=dash_font_title)
        
        # Title
        t_text = "AutoDev Ready"
        draw.text((TERM_X + (TERM_WIDTH - len(t_text) * 15) // 2, TERM_Y + 130), t_text, fill=(255, 255, 255), font=dash_font_title)
        
        # Subtitle description
        desc_lines = [
            "Automatically bootstrapped developer environments, installed missing",
            "tools, audited security, and generated roadmaps in under 2 minutes."
        ]
        sub_y = TERM_Y + 180
        for dl in desc_lines:
            draw.text((TERM_X + (TERM_WIDTH - len(dl) * 7.5) // 2, sub_y), dl, fill=(148, 163, 184), font=dash_font_sub)
            sub_y += 20
            
        # Draw metric grid cards
        metrics = [
            {"val": "30 min -> 2 min", "label": "Setup Time Saved"},
            {"val": "1 Command", "label": "Bootstrapping Cost"},
            {"val": "100% Secure", "label": "Vulnerability Scans Passed"},
            {"val": "Interactive", "label": "Learning Roadmaps"}
        ]
        
        # Grid positions
        card_w, card_h = 190, 70
        grid_x = [TERM_X + 110, TERM_X + 320, TERM_X + 530] # we want 2x2 grid centered
        
        # Let's manually draw cards in 2x2 layout
        card_coords = [
            (TERM_X + 210, TERM_Y + 250), # row 1 col 1
            (TERM_X + 440, TERM_Y + 250), # row 1 col 2
            (TERM_X + 210, TERM_Y + 340), # row 2 col 1
            (TERM_X + 440, TERM_Y + 340), # row 2 col 2
        ]
        
        for idx, m in enumerate(metrics):
            cx, cy = card_coords[idx]
            # Draw card outline
            draw.rounded_rectangle([cx, cy, cx + card_w, cy + card_h], radius=6, fill=(255, 255, 255, 5), outline=(99, 102, 241, 100), width=1)
            # Draw value
            val_w = len(m["val"]) * 11
            draw.text((cx + (card_w - val_w) // 2, cy + 15), m["val"], fill=(129, 140, 248), font=dash_font_val)
            # Draw label
            lbl_w = len(m["label"]) * 7
            draw.text((cx + (card_w - lbl_w) // 2, cy + 42), m["label"], fill=(100, 116, 139), font=dash_font_label)
            
    return img

# Definitions of output segments to append to terminal buffer
OUTPUT_AUTODEV = [
    "  [magenta bold]█████╗ ██╗   ██╗████████╗ ██████╗ ██████╗ ███████╗██╗   ██╗[/magenta]",
    "  [magenta bold]██╔══██╗██║   ██║╚══██╔══╝██╔═══██╗██╔══██╗██╔════╝██║   ██║[/magenta]",
    "  [magenta bold]███████║██║   ██║   ██║   ██║   ██║██║  ██║█████╗  ██║   ██║[/magenta]",
    "  [magenta bold]██╔══██║██║   ██║   ██║   ██║   ██║██║  ██║██╔══╝  ╚██╗ ██╔╝[/magenta]",
    "  [magenta bold]██║  ██║╚██████╔╝   ██║   ╚██████╔╝██████╔╝███████╗ ╚████╔╝ [/magenta]",
    "  [magenta bold]╚═╝  ╚═╝ ╚═════╝    ╚═╝    ╚══════╝ ╚═════╝ ╚══════╝  ╚═══╝ [/magenta]",
    "",
    "  [bold]The App Store for Developers.[/bold]",
    "  Run with no arguments to open the interactive installer.",
    "",
    "[yellow bold]Usage:[/yellow]",
    "  autodev [flags]",
    "  autodev [command]",
    "",
    "[yellow bold]Available Commands:[/yellow]",
    "  [green bold]audit[/green]        Audit repository dependencies for security vulnerabilities",
    "  [green bold]benchmark[/green]    Display AI token and efficiency benchmarks",
    "  [green bold]clean[/green]        Remove AutoDev cache and temp files",
    "  [green bold]clone[/green]        Clone a Git repo, scan it, and install all missing tools",
    "  [green bold]containerize[/green] Generate DevContainer and VSCode workspace configuration",
    "  [green bold]create[/green]       Create a new pre-configured boilerplate project",
    "  [green bold]doctor[/green]       Check the health and security of your codebase",
    "  [green bold]export[/green]       Export environment as a reproducible JSON lockfile",
    "  [green bold]profile[/green]      Install a pre-defined developer profile (role-based tool set)",
    "  [green bold]scan[/green]         Scan a repository for languages, frameworks, and dependencies",
    "  [green bold]setup[/green]        Detect and install all missing runtimes and dependencies",
    "  [green bold]skills[/green]       Generate a personalized learning roadmap",
    "  [green bold]ui[/green]           Start the local AutoDev interactive web dashboard",
    "",
    "Use \"autodev [command] --help\" for more information."
]

OUTPUT_SCAN = [
    "┌───────────────┐",
    "│  [cyan bold]Scanning: .[/cyan]  │",
    "└───────────────┘",
    "",
    "  [yellow bold][LANGUAGES][/yellow]",
    "    - Node.js",
    "    - TypeScript",
    "    - Go",
    "",
    "  [yellow bold][FRAMEWORKS][/yellow]",
    "    - Next.js",
    "    - React",
    "",
    "  [yellow bold][PACKAGE MANAGERS][/yellow]",
    "    - pnpm",
    "",
    "  [yellow bold][INFRASTRUCTURE][/yellow]",
    "    - Docker",
    "    - Kubernetes",
    "",
    "   [green bold]Docker[/green]    Container-ready",
    "   [green bold]Kubernetes[/green]   Found k8s manifests",
    "",
    "  [yellow bold][MONOREPO SUBPROJECTS][/yellow]",
    "    - website (apps/website) -> Next.js, Node.js, TypeScript, React",
    "    - cli (packages/cli) -> Go",
    "    - scanner (packages/scanner) -> Go",
    "    - skills (packages/skills) -> Go",
    "    - core (packages/core) -> Go",
    "",
    "  [yellow bold][SETUP PLAN][/yellow]",
    "    1. autodev install docker",
    "    2. autodev install nodejs",
    "    3. autodev install go",
    "    4. npm install -g pnpm",
    "",
    "  [green]Scanned in 4ms | 8 technologies detected[/green]",
    "  Run '[cyan]autodev setup[/cyan]' to install all missing tools.",
    "  Run '[cyan]autodev audit[/cyan]' to check dependencies for security risks.",
    "",
    "  [info bold]🛡️  PROJECT ENHANCEMENTS DETECTED (Next.js):[/info]",
    "  ──────────────────────────────────────────────────────────",
    "    [green]✓[/green]  Tailwind CSS           Already Configured",
    "    [green]✓[/green]  Dockerfile             Already Configured",
    "    [red]✗[/red]  ESLint & Prettier      Not Configured - Lint rules check",
    "    [green]✓[/green]  GitHub Actions CI/CD   Already Configured",
    "  ──────────────────────────────────────────────────────────",
    "",
    "  Configure missing enhancements? [y/N] "
]

OUTPUT_SCAN_SKIP = [
    "  Enhancements skipped.",
    "",
    "  ──────────────────────────────────────────────────────────",
    "  ⭐ Love this tool? Star the repo to support AutoDev: [info]https://github.com/HEETMEHTA18/autodev[/info]",
    "  ──────────────────────────────────────────────────────────"
]

OUTPUT_DOCTOR = [
    "[cyan bold]⚡ AUTODEV CODEBASE HEALTH & SECURITY DOCTOR[/cyan]",
    "",
    "[bold]SYSTEM SPECIFICATIONS[/bold]",
    "  OS                   Pop!_OS 24.04 LTS",
    "  Architecture         amd64",
    "  Package Manager      apt",
    "",
    "[bold]CODEBASE DIAGNOSTICS SCAN[/bold]",
    "  Scanning for secrets, configuration mismatches, and code errors...",
    "",
    "  [[green bold]OK[/green]]       Git Configuration (.gitignore) Clean & Healthy",
    "  [[green bold]OK[/green]]       Exposed Secrets Scanner   Clean & Healthy",
    "  [[green bold]OK[/green]]       Environment Config (.env) Clean & Healthy",
    "  [[green bold]OK[/green]]       Dependency Lockfiles      Clean & Healthy",
    "  [[green bold]OK[/green]]       Linter & Code Format Status Clean & Healthy",
    "  [[green bold]OK[/green]]       Environment Lockfile Mismatch Clean & Healthy",
    "  [[green bold]OK[/green]]       Supply-Chain Vulnerabilities Clean & Healthy",
    "",
    "  [green bold]✓ Codebase is completely healthy, secure, and ready for production![/green]"
]

OUTPUT_AUDIT = [
    "[cyan bold]🛡️  AutoDev Supply-Chain Safety Audit[/cyan]",
    "  Auditing dependencies against the OSV Vulnerability Database...",
    "",
    "  [green bold]✓ No known security vulnerabilities found! All dependencies are safe.[/green]",
    "",
    "  ──────────────────────────────────────────────────────────",
    "  ⭐ Love this tool? Star the repo to support AutoDev: [info]https://github.com/HEETMEHTA18/autodev[/info]",
    "  ──────────────────────────────────────────────────────────"
]

OUTPUT_SKILLS = [
    "  [magenta bold]⚡ AutoDev Skills Engine v0.3.0[/magenta]",
    "  Powered by skills.sh",
    "",
    "╔══════════════════════════════════════════════════",
    "║  [bold]Personalized Roadmap[/bold] - Node.js, TypeScript, Go, Next.js, React",
    "╚══════════════════════════════════════════════════",
    "",
    "  [yellow bold][CURRENT SKILLS DETECTED][/yellow]",
    "   - Node.js              [gray][beginner][/gray] Runtime",
    "   - TypeScript           [info][intermediate][/info] Language",
    "   - Go                   [info][intermediate][/info] Language",
    "   - Next.js              [info][intermediate][/info] Framework",
    "   - React                [info][intermediate][/info] Framework",
    "",
    "  [yellow bold][RECOMMENDED NEXT STEPS][/yellow]",
    "   - Express              [gray][beginner][/gray]",
    "     Link: [info]https://expressjs.com/en/starter/installing.html[/info]",
    "   - Docker               [info][intermediate][/info]",
    "     Link: [info]https://docs.docker.com/get-started[/info]",
    "   - NestJS               [info][intermediate][/info]",
    "     Link: [info]https://docs.nestjs.com[/info]",
    "   - CI/CD                [info][intermediate][/info]",
    "     Link: [info]https://docs.github.com/actions[/info]",
    "   - Kubernetes           [magenta][advanced][/magenta]",
    "     Link: [info]https://kubernetes.io/docs/tutorials[/info]",
    "   - Terraform            [magenta][advanced][/magenta]",
    "     Link: [info]https://developer.hashicorp.com/terraform/tutorials[/info]",
    "",
    "  [yellow bold][LONG-TERM GOALS][/yellow]",
    "   - PostgreSQL           [info][intermediate][/info]",
    "",
    "  Run '[cyan]autodev skills --save-rules[/cyan]' to download AI instructions locally.",
    "  Visit [info]https://skills.sh[/info] for interactive learning paths."
]

def main():
    buffer = []
    frames = []
    durations = []
    
    # Helper to add typing frames
    def add_typing(cmd_text):
        for i in range(1, len(cmd_text) + 1):
            frames.append(render_frame(buffer, active_cmd=cmd_text[:i], show_cursor=True))
            durations.append(80) # 80ms typing speed
            
        # Settle frame at the end of typing
        frames.append(render_frame(buffer, active_cmd=cmd_text, show_cursor=False))
        durations.append(250)
        
    # Helper to append lines to buffer and render output
    def add_output(lines, display_duration=2000):
        # We append all lines to buffer
        buffer.append(f"[prompt]heet@pop-os:~/projects/autodev$ [cmd]{frames[-1].info.get('cmd', '')}")
        # Update last item in buffer to reflect actual command typed
        # Wait, let's just create prompt line explicitly in the buffer
        pass

    # Let's construct frame-by-frame sequences manually
    
    # 1. Type autodev
    add_typing("autodev")
    
    # Append the command prompt to the buffer
    buffer.append("[prompt]heet@pop-os:~/projects/autodev$ [cmd]autodev")
    # Add help lines
    for line in OUTPUT_AUTODEV:
        buffer.append(line)
    
    # Render static display frame for autodev help
    frames.append(render_frame(buffer, active_cmd=None))
    durations.append(2500) # Show help for 2.5s
    
    # 2. Type autodev scan
    add_typing("autodev scan")
    buffer.append("[prompt]heet@pop-os:~/projects/autodev$ [cmd]autodev scan")
    for line in OUTPUT_SCAN:
        buffer.append(line)
        
    # Render scan output frame
    frames.append(render_frame(buffer, active_cmd=None))
    durations.append(3000) # Show scan results for 3s
    
    # 3. Type N at prompt
    # Update the last line of buffer (which is the prompt) to show N
    buffer[-1] = buffer[-1] + "[cmd]N"
    frames.append(render_frame(buffer, active_cmd=None))
    durations.append(500)
    
    # Show skipped enhancements output
    for line in OUTPUT_SCAN_SKIP:
        buffer.append(line)
    frames.append(render_frame(buffer, active_cmd=None))
    durations.append(1500)
    
    # 4. Type autodev doctor
    add_typing("autodev doctor")
    buffer.append("[prompt]heet@pop-os:~/projects/autodev$ [cmd]autodev doctor")
    for line in OUTPUT_DOCTOR:
        buffer.append(line)
    frames.append(render_frame(buffer, active_cmd=None))
    durations.append(2500)
    
    # 5. Type autodev audit
    add_typing("autodev audit")
    buffer.append("[prompt]heet@pop-os:~/projects/autodev$ [cmd]autodev audit")
    for line in OUTPUT_AUDIT:
        buffer.append(line)
    frames.append(render_frame(buffer, active_cmd=None))
    durations.append(2000)
    
    # 6. Type autodev skills
    add_typing("autodev skills")
    buffer.append("[prompt]heet@pop-os:~/projects/autodev$ [cmd]autodev skills")
    for line in OUTPUT_SKILLS:
        buffer.append(line)
    frames.append(render_frame(buffer, active_cmd=None))
    durations.append(3000)
    
    # 7. Render ready dashboard overlay screen
    frames.append(render_frame(buffer, active_cmd=None, show_overlay=True))
    durations.append(5000) # Show final screen for 5s
    
    # Save animated GIF
    output_path_gif = "/media/heet18/Futuristic/Heet/Github/Autodev/autodev-demo.gif"
    print(f"Saving {len(frames)} frames to {output_path_gif}...")
    frames[0].save(
        output_path_gif,
        save_all=True,
        append_images=frames[1:],
        duration=durations,
        loop=0,
        optimize=True
    )
    print("GIF saved successfully!")

    # Save animated WebP
    output_path_webp = "/media/heet18/Futuristic/Heet/Github/Autodev/autodev-demo.webp"
    print(f"Saving {len(frames)} frames to {output_path_webp}...")
    frames[0].save(
        output_path_webp,
        save_all=True,
        append_images=frames[1:],
        duration=durations,
        loop=0,
        optimize=True
    )
    print("WebP saved successfully!")

if __name__ == "__main__":
    main()
