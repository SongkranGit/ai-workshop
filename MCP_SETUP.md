# MCP Setup Guide - Browser Automation & GitHub Integration

## What is MCP?

MCP (Model Context Protocol) allows GitHub Copilot to interact with external tools and services. In this project, we've configured:
- 🎭 **Playwright MCP**: Browser automation and testing
- 🐙 **GitHub MCP**: Repository analysis and code style guidelines

## Configuration

File: `.vscode/mcp.json`

```json
{
  "mcpServers": {
    "playwright": {
      "command": "npx",
      "args": ["-y", "@executeautomation/playwright-mcp-server"]
    },
    "github": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-github"],
      "env": {
        "GITHUB_PERSONAL_ACCESS_TOKEN": ""
      }
    }
  }
}
```

## Prerequisites

- ✅ VS Code with GitHub Copilot extension
- ✅ Node.js and npm installed
- ✅ npx available (comes with npm)
- ✅ GitHub Personal Access Token (for GitHub MCP)

## How to Use

### Playwright MCP Examples

#### Example 1: Check Website Contact Location

**Prompt:**
```
ใช้ Playwright MCP ตรวจสอบ https://www.kbtg.tech/th/home ว่าเว็บไซต์ใส่ Contact อยู่ที่ Nonthaburi ใช่ไหม
```

**What it does:**
1. Opens browser via Playwright
2. Navigates to KBTG website
3. Finds contact/location information
4. Verifies if "Nonthaburi" is mentioned

#### Example 2: Test Transfer API Swagger UI

**Prompt:**
```
ใช้ Playwright MCP เปิด http://localhost:3000/swagger และทดสอบ POST /api/transfers
```

**What it does:**
1. Opens Swagger UI in browser
2. Interacts with API documentation
3. Executes test requests
4. Validates responses

### GitHub MCP Examples

#### Example 1: Follow Repository Style Guidelines

**Prompt:**
```
จง follow style guideline ตาม repo kbtg-ai-workshop-oct Docs
```

**What it does:**
1. Reads code patterns from https://github.com/mikelopster/kbtg-ai-workshop-oct
2. Analyzes workshop documentation and examples
3. Applies style guidelines to your code
4. Suggests improvements based on workshop patterns

#### Example 2: Compare Implementation

**Prompt:**
```
ใช้ GitHub MCP เปรียบเทียบ code structure ของเรากับ repo kbtg-ai-workshop-oct
```

**What it does:**
1. Fetches repository structure from GitHub
2. Compares with current implementation
3. Highlights differences
4. Suggests alignment improvements

## Available MCP Commands

### Playwright MCP Tools

- 🌐 **Browser Navigation**: `playwright_navigate`
- 🖱️ **Element Interaction**: `playwright_click`, `playwright_fill`
- 📸 **Screenshots**: `playwright_screenshot`
- ✅ **Assertions**: Check text, attributes, visibility
- 🔍 **Element Selection**: By CSS, XPath, text content

### GitHub MCP Tools

- 📁 **Repository Access**: `github_list_repos`, `github_get_repo`
- 📄 **File Operations**: `github_get_file_contents`, `github_search_code`
- 🔍 **Search**: `github_search_repositories`, `github_search_issues`
- 📊 **Analysis**: Read code patterns, documentation, examples
- 🔄 **Compare**: Analyze differences between implementations

## Troubleshooting

### Issue: MCP server not loading

**Solution:**
```bash
# Verify npx is available
npx --version

# Test Playwright MCP manually
npx -y @executeautomation/playwright-mcp-server

# Test GitHub MCP manually
npx -y @modelcontextprotocol/server-github
```

### Issue: GitHub MCP authentication failed

**Solution:**
1. Create Personal Access Token at https://github.com/settings/tokens
2. Add token to `.vscode/mcp.json` or export as environment variable:
   ```bash
   export GITHUB_PERSONAL_ACCESS_TOKEN="your_token_here"
   ```
3. Restart VS Code

### Issue: VS Code doesn't recognize MCP

**Solution:**
1. Restart VS Code
2. Check GitHub Copilot is enabled
3. Verify `.vscode/mcp.json` exists in workspace root

### Issue: Permission denied

**Solution:**
```bash
# Make sure npm global packages are accessible
npm config get prefix

# Update npm if needed
npm install -g npm@latest
```

## Verification Steps

1. ✅ `.vscode/mcp.json` file exists
2. ✅ npx command is available: `which npx`
3. ✅ GitHub token configured (for GitHub MCP)
4. ✅ VS Code workspace opened at repository root
5. ✅ GitHub Copilot extension active
6. ✅ Try a simple prompt to test both MCPs

## Example Prompts to Test

### Test Playwright MCP
```
ใช้ Playwright MCP เปิด https://www.kbtg.tech และบอกว่ามีหัวข้ออะไรบ้าง
```

### Test GitHub MCP
```
ใช้ GitHub MCP อ่าน README จาก repo mikelopster/kbtg-ai-workshop-oct
```

### Combined Test
```
จง follow style guideline ตาม repo kbtg-ai-workshop-oct และแสดงตัวอย่าง code ที่ควรปรับปรุง
```

## Example Output

When you ask Copilot to use Playwright MCP:

```
🤖 Using Playwright MCP to navigate to https://www.kbtg.tech/th/home...

✅ Page loaded successfully
📍 Found contact section
🔍 Searching for "Nonthaburi"...
✅ Confirmed: Contact address is in Nonthaburi
```

## Resources

- 📚 [VS Code MCP Documentation](https://code.visualstudio.com/docs/copilot/customization/mcp-servers)
- 🎭 [Playwright MCP Server](https://mcp.so/server/playwright-mcp/microsoft)
- � [GitHub MCP Server](https://github.com/modelcontextprotocol/servers/tree/main/src/github)
- �🛠️ [Workshop Guide](https://github.com/mikelopster/kbtg-ai-workshop-oct/blob/main/workshop-5/prompt.md)
- 📖 [Workshop Repository](https://github.com/mikelopster/kbtg-ai-workshop-oct)

## Notes

- MCP servers run locally via npx
- First run may take longer (npm package download)
- Browser sessions are temporary (closed after task)
- Playwright runs in headless mode by default
- GitHub MCP requires Personal Access Token
- Token should have `repo` and `read:org` scopes
- Never commit tokens to git repository

---

**Setup Date:** October 17, 2025  
**Workshop:** KBTG AI Workshop - Workshop 5  
**MCPs Configured:** Playwright + GitHub
