# MCP Setup Guide - Playwright Browser Automation

## What is MCP?

MCP (Model Context Protocol) allows GitHub Copilot to interact with external tools and services. In this project, we've configured Playwright MCP for browser automation.

## Configuration

File: `.vscode/mcp.json`

```json
{
  "mcpServers": {
    "playwright": {
      "command": "npx",
      "args": [
        "-y",
        "@executeautomation/playwright-mcp-server"
      ]
    }
  }
}
```

## Prerequisites

- ✅ VS Code with GitHub Copilot extension
- ✅ Node.js and npm installed
- ✅ npx available (comes with npm)

## How to Use

### Example 1: Check Website Contact Location

**Prompt:**
```
ใช้ Playwright MCP ตรวจสอบ https://www.kbtg.tech/th/home ว่าเว็บไซต์ใส่ Contact อยู่ที่ Nonthaburi ใช่ไหม
```

**What it does:**
1. Opens browser via Playwright
2. Navigates to KBTG website
3. Finds contact/location information
4. Verifies if "Nonthaburi" is mentioned

### Example 2: Test Transfer API Swagger UI

**Prompt:**
```
ใช้ Playwright MCP เปิด http://localhost:3000/swagger และทดสอบ POST /api/transfers
```

**What it does:**
1. Opens Swagger UI in browser
2. Interacts with API documentation
3. Executes test requests
4. Validates responses

### Example 3: Automated Testing

**Prompt:**
```
ใช้ Playwright MCP ทดสอบว่า homepage โหลดภายใน 3 วินาที และมีคำว่า "Transfer" ปรากฏ
```

## Available MCP Commands

The Playwright MCP server provides tools for:

- 🌐 **Browser Navigation**: `playwright_navigate`
- 🖱️ **Element Interaction**: `playwright_click`, `playwright_fill`
- 📸 **Screenshots**: `playwright_screenshot`
- ✅ **Assertions**: Check text, attributes, visibility
- 🔍 **Element Selection**: By CSS, XPath, text content

## Troubleshooting

### Issue: MCP server not loading

**Solution:**
```bash
# Verify npx is available
npx --version

# Test Playwright MCP manually
npx -y @executeautomation/playwright-mcp-server
```

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
3. ✅ VS Code workspace opened at repository root
4. ✅ GitHub Copilot extension active
5. ✅ Try a simple Playwright prompt to test

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
- 🛠️ [Workshop Guide](https://github.com/mikelopster/kbtg-ai-workshop-oct/blob/main/workshop-5/prompt.md)

## Notes

- MCP servers run locally via npx
- First run may take longer (npm package download)
- Browser sessions are temporary (closed after task)
- Playwright runs in headless mode by default

---

**Setup Date:** October 17, 2025  
**Workshop:** KBTG AI Workshop - Workshop 5
