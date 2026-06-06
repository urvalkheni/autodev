import os
import re
from PIL import Image, ImageDraw, ImageFont

# Dimensions
TERM_WIDTH = 840
HEADER_HEIGHT = 40
PADDING = 20

# Font setup
MONO_FONT_PATH = "/usr/share/fonts/truetype/freefont/FreeMono.ttf"
MONO_BOLD_FONT_PATH = "/usr/share/fonts/truetype/freefont/FreeMonoBold.ttf"

font_size = 14
font = ImageFont.truetype(MONO_FONT_PATH, font_size)
font_bold = ImageFont.truetype(MONO_BOLD_FONT_PATH, font_size)

# Calculate character dimensions
bbox = font.getbbox("A")
CHAR_WIDTH = bbox[2] - bbox[0]
CHAR_HEIGHT = bbox[3] - bbox[1]
LINE_HEIGHT = CHAR_HEIGHT + 6

# Color scheme
COLORS = {
    "default": (226, 232, 240), # slate-200
    "prompt": (56, 189, 248),  # sky-400
    "cmd": (255, 255, 255),     # white
    "success": (16, 185, 129), # emerald-500
    "green": (16, 185, 129),
    "warn": (251, 191, 36),    # amber-400
    "yellow": (251, 191, 36),
    "error": (239, 68, 68),    # red-500
    "red": (239, 68, 68),
    "info": (129, 140, 248),   # indigo-400
    "blue": (129, 140, 248),
    "gray": (100, 116, 139),   # slate-500
    "magenta": (244, 114, 182), # pink-400
    "cyan": (34, 211, 238)     # cyan-400
}

def parse_formatted_line(line):
    pattern = re.compile(r'\[([^\]]+)\]')
    tokens = []
    parts = pattern.split(line)
    current_color = COLORS["default"]
    current_bold = False
    
    for i, part in enumerate(parts):
        if i % 2 == 1:
            tag_parts = part.lower().split()
            current_bold = "bold" in tag_parts
            color_name = next((tp for tp in tag_parts if tp != "bold"), "default")
            current_color = COLORS.get(color_name, COLORS["default"])
        else:
            if part:
                tokens.append((part, current_color, current_bold))
    return tokens

def render_screenshot(output_lines, filename):
    num_lines = len(output_lines)
    content_height = num_lines * LINE_HEIGHT
    term_height = HEADER_HEIGHT + content_height + (PADDING * 2)
    
    # Create image canvas (transparent border for nice drop shadow / outline)
    border_padding = 10
    img_w = TERM_WIDTH + (border_padding * 2)
    img_h = term_height + (border_padding * 2)
    
    img = Image.new("RGBA", (img_w, img_h), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)
    
    term_x = border_padding
    term_y = border_padding
    
    # Draw drop shadow
    draw.rounded_rectangle(
        [term_x + 4, term_y + 4, term_x + TERM_WIDTH + 4, term_y + term_height + 4],
        radius=12, fill=(0, 0, 0, 80)
    )
    # Draw terminal background
    draw.rounded_rectangle(
        [term_x, term_y, term_x + TERM_WIDTH, term_y + term_height],
        radius=12, fill=(13, 13, 23, 255)
    )
    # Draw border
    draw.rounded_rectangle(
        [term_x, term_y, term_x + TERM_WIDTH, term_y + term_height],
        radius=12, outline=(50, 50, 75, 255), width=1
    )
    
    # Header background
    draw.rounded_rectangle(
        [term_x + 1, term_y + 1, term_x + TERM_WIDTH - 1, term_y + HEADER_HEIGHT],
        radius=0, fill=(20, 20, 35, 255)
    )
    # Re-draw top corners of header
    draw.pieslice([term_x, term_y, term_x + 24, term_y + 24], 180, 270, fill=(20, 20, 35, 255))
    draw.pieslice([term_x + TERM_WIDTH - 24, term_y, term_x + TERM_WIDTH, term_y + 24], 270, 360, fill=(20, 20, 35, 255))
    draw.rectangle([term_x + 12, term_y + 1, term_x + TERM_WIDTH - 12, term_y + 20], fill=(20, 20, 35, 255))
    
    # Traffic lights
    draw.ellipse([term_x + 16, term_y + 14, term_x + 28, term_y + 26], fill=(239, 68, 68, 255)) # red
    draw.ellipse([term_x + 36, term_y + 14, term_x + 48, term_y + 26], fill=(234, 179, 8, 255)) # yellow
    draw.ellipse([term_x + 56, term_y + 14, term_x + 68, term_y + 26], fill=(34, 197, 94, 255)) # green
    
    # Header title
    title_text = "heet@pop-os: ~/projects/autodev"
    draw.text((term_x + (TERM_WIDTH - len(title_text) * CHAR_WIDTH) // 2, term_y + 13), title_text, fill=(100, 116, 139, 255), font=font)
    
    # Header badge
    badge_text = "v0.3.2"
    draw.rounded_rectangle([term_x + TERM_WIDTH - 70, term_y + 10, term_x + TERM_WIDTH - 16, term_y + 30], radius=4, fill=(99, 102, 241, 40), outline=(99, 102, 241, 255), width=1)
    draw.text((term_x + TERM_WIDTH - 60, term_y + 13), badge_text, fill=(165, 180, 252, 255), font=font)
    
    # Render body
    current_y = term_y + HEADER_HEIGHT + PADDING
    for raw_line in output_lines:
        tokens = parse_formatted_line(raw_line)
        current_x = term_x + 20
        
        for text, color, is_bold in tokens:
            fnt = font_bold if is_bold else font
            draw.text((current_x, current_y), text, fill=color, font=fnt)
            current_x += len(text) * CHAR_WIDTH
            
        current_y += LINE_HEIGHT
        
    os.makedirs(os.path.dirname(filename), exist_ok=True)
    img.save(filename, "PNG")
    print(f"Saved: {filename}")

# Definitions of commands outputs
outputs = {
    "screenshot-create.png": [
        "[prompt]heet@pop-os:~/projects/autodev$ [cmd]autodev create nextjs my-app",
        "",
        "[cyan bold]Creating a new project from Next.js template...[/cyan]",
        "  [green]✓[/green] Clone template [bold]nextjs[/bold]",
        "  [green]✓[/green] Installed dependencies ([bold]pnpm install[/bold])",
        "  [green]✓[/green] Configured linters and git hooks",
        "  [green]✓[/green] Dockerfile and devcontainer configuration generated",
        "",
        "[success bold]Project 'my-app' successfully created![/success]",
        "To start developing, run the following commands:",
        "  [info]cd my-app && autodev setup[/info]",
        "",
        "  ──────────────────────────────────────────────────────────",
        "  ⭐ Star the repo to support AutoDev: [info]https://github.com/HEETMEHTA18/autodev[/info]",
        "  ──────────────────────────────────────────────────────────"
    ],
    "screenshot-clone.png": [
        "[prompt]heet@pop-os:~/projects/autodev$ [cmd]autodev clone https://github.com/example/repo.git",
        "",
        "[cyan bold]Cloning repository 'https://github.com/example/repo.git'...[/cyan]",
        "  Cloning into 'repo'...",
        "  [green]✓[/green] Cloned successfully",
        "",
        "[cyan bold]Scanning stack and installing dependencies...[/cyan]",
        "  Detected: Node.js, Go, Docker",
        "  [green]✓[/green] Installed Node.js v20.11.0",
        "  [green]✓[/green] Installed Go v1.22.0",
        "  [green]✓[/green] Installed project packages (npm install)",
        "",
        "[success bold]✓ Repository cloned and dev environment setup completed successfully![/success]",
        "",
        "  ──────────────────────────────────────────────────────────",
        "  ⭐ Star the repo to support AutoDev: [info]https://github.com/HEETMEHTA18/autodev[/info]",
        "  ──────────────────────────────────────────────────────────"
    ],
    "screenshot-install.png": [
        "[prompt]heet@pop-os:~/projects/autodev$ [cmd]autodev install nodejs",
        "",
        "[cyan bold]Installing nodejs...[/cyan]",
        "  [green]✓[/green] Downloaded nodejs binary archive (v20.11.0)",
        "  [green]✓[/green] Extracted and linked nodejs to system path",
        "  [green]✓[/green] Verified installation: node --version -> [success]v20.11.0[/success]",
        "",
        "[success bold]✓ nodejs installed successfully![/success]",
        "",
        "  ──────────────────────────────────────────────────────────",
        "  ⭐ Star the repo to support AutoDev: [info]https://github.com/HEETMEHTA18/autodev[/info]",
        "  ──────────────────────────────────────────────────────────"
    ],
    "screenshot-skills.png": [
        "[prompt]heet@pop-os:~/projects/autodev$ [cmd]autodev skills",
        "",
        "  [magenta bold]⚡ AutoDev Skills Engine v0.3.2[/magenta]",
        "  Powered by skills.sh",
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
        "   - Docker               [info][intermediate][/info]",
        "   - Kubernetes           [magenta][advanced][/magenta]",
        "",
        "  [success]✓ Telemetry summary saved to .autodev-skills.md[/success]",
        "  Run '[cyan]autodev skills --save-rules[/cyan]' to download AI instructions locally.",
        "",
        "  ──────────────────────────────────────────────────────────",
        "  ⭐ Star the repo to support AutoDev: [info]https://github.com/HEETMEHTA18/autodev[/info]",
        "  ──────────────────────────────────────────────────────────"
    ],
    "screenshot-mcp.png": [
        "[prompt]heet@pop-os:~/projects/autodev$ [cmd]autodev mcp start",
        "",
        "[cyan bold]Starting AutoDev Model Context Protocol (MCP) server...[/cyan]",
        "  [green]✓[/green] MCP Server listening on stdio interface",
        "  [green]✓[/green] Registered tools with backend parser",
        "",
        "  [bold]EXPOSED TOOLS:[/bold]",
        "    - autodev_scan            Scan codebase technologies",
        "    - autodev_setup           Bootstrap local project toolchains",
        "    - autodev_install         Install managed compiler runtimes",
        "    - autodev_audit           Scan dependency vulnerabilities",
        "    - autodev_doctor          Perform codebase health diagnostics",
        "    - autodev_containerize    Generate DevContainer configurations",
        "",
        "[success bold]● MCP Server is running and waiting for client connections...[/success]",
        "",
        "  ──────────────────────────────────────────────────────────",
        "  ⭐ Star the repo to support AutoDev: [info]https://github.com/HEETMEHTA18/autodev[/info]",
        "  ──────────────────────────────────────────────────────────"
    ],
    "screenshot-benchmark.png": [
        "[prompt]heet@pop-os:~/projects/autodev$ [cmd]autodev benchmark",
        "",
        "  [magenta bold]⚡ AI Token Efficiency Benchmark[/magenta]",
        "",
        "  [bold]METRIC                 TRADITIONAL PROMPTING   AUTODEV RULES[/bold]",
        "  Context Window Size    84,500 tokens           4,100 tokens",
        "  API Response Latency   24.5 seconds            1.8 seconds",
        "  Estimated Cost / Query $0.84                   $0.04",
        "  Efficiency Ratio       1.00x                   [success bold]21.00x (95% savings)[/success]",
        "",
        "  [success]✓ Benchmarks successfully analyzed![/success]",
        "",
        "  ──────────────────────────────────────────────────────────",
        "  ⭐ Star the repo to support AutoDev: [info]https://github.com/HEETMEHTA18/autodev[/info]",
        "  ──────────────────────────────────────────────────────────"
    ],
    "screenshot-report.png": [
        "[prompt]heet@pop-os:~/projects/autodev$ [cmd]autodev report",
        "",
        "[cyan bold]Generating configuration report...[/cyan]",
        "  [green]✓[/green] Scanned codebase technologies",
        "  [green]✓[/green] Checked system compiler status",
        "  [green]✓[/green] Run supply-chain security audit",
        "  [green]✓[/green] Exported report format: html",
        "",
        "[success bold]✓ Report successfully saved to: ./autodev-report.html[/success]",
        "",
        "  ──────────────────────────────────────────────────────────",
        "  ⭐ Star the repo to support AutoDev: [info]https://github.com/HEETMEHTA18/autodev[/info]",
        "  ──────────────────────────────────────────────────────────"
    ],
    "screenshot-github.png": [
        "[prompt]heet@pop-os:~/projects/autodev$ [cmd]autodev github HEETMEHTA18",
        "",
        "[cyan bold]Scanning public repositories for user: HEETMEHTA18...[/cyan]",
        "  [green]✓[/green] Retrieved 14 repository meta specs",
        "",
        "  [bold]AGGREGATE TECHNOLOGY FOOTPRINT:[/bold]",
        "    Go              [cyan]■■■■■■■■■■■■■■■■■■■■[/cyan]  54%",
        "    TypeScript      [cyan]■■■■■■■■■■■■[/cyan]          32%",
        "    Python          [cyan]■■■■[/cyan]                  10%",
        "    Other           [cyan]■■[/cyan]                     4%",
        "",
        "[success bold]✓ Scan completed![/success]",
        "",
        "  ──────────────────────────────────────────────────────────",
        "  ⭐ Star the repo to support AutoDev: [info]https://github.com/HEETMEHTA18/autodev[/info]",
        "  ──────────────────────────────────────────────────────────"
    ],
    "screenshot-exec.png": [
        "[prompt]heet@pop-os:~/projects/autodev$ [cmd]autodev exec go run main.go",
        "",
        "[cyan]Running command inside AutoDev isolated environment...[/cyan]",
        "  [green]✓[/green] Resolved isolated PATH and environment variables",
        "  [green]✓[/green] Initialized Go virtual sandbox",
        "",
        "Hello from AutoDev! Executed command successfully.",
        "",
        "[success]✓ Environment restored.[/success]",
        "",
        "  ──────────────────────────────────────────────────────────",
        "  ⭐ Star the repo to support AutoDev: [info]https://github.com/HEETMEHTA18/autodev[/info]",
        "  ──────────────────────────────────────────────────────────"
    ],
    "screenshot-prompts.png": [
        "[prompt]heet@pop-os:~/projects/autodev$ [cmd]autodev prompts --today",
        "",
        "  [yellow bold]📅 Captured Prompts for Today (2026-06-06)[/yellow]",
        "",
        "  [green bold]●[/green] [bold]session-39a2f1b[/bold] (TypeScript, Node.js)",
        "    [gray]16:15:32[/gray] Add support for SQLite DB schema",
        "    [gray]16:21:40[/gray] Optimize SQL query indexing for prompts",
        "",
        "  [green bold]●[/green] [bold]session-c2f8812[/bold] (Go)",
        "    [gray]14:10:05[/gray] Write unit tests for capture stdin proxy",
        "",
        "  [success]✓ 2 active capture sessions logged today.[/success]",
        "",
        "  ──────────────────────────────────────────────────────────",
        "  ⭐ Star the repo to support AutoDev: [info]https://github.com/HEETMEHTA18/autodev[/info]",
        "  ──────────────────────────────────────────────────────────"
    ]
}

def main():
    public_dir = "apps/website/public"
    for img_name, lines in outputs.items():
        out_path = os.path.join(public_dir, img_name)
        render_screenshot(lines, out_path)

if __name__ == "__main__":
    main()
