# GitHub MCP Setup

## Setup GitHub Personal Access Token

1. Go to https://github.com/settings/tokens
2. Click "Generate new token (classic)"
3. Select scopes:
   - ✅ `repo` (Full control of private repositories)
   - ✅ `read:org` (Read org and team membership)
4. Copy the token

## Configure MCP

Edit `.vscode/mcp.json` and add your token:

```json
{
  "mcpServers": {
    "github": {
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-github"],
      "env": {
        "GITHUB_PERSONAL_ACCESS_TOKEN": "your_token_here"
      }
    }
  }
}
```

⚠️ **Important:** Never commit your token to git! The token field is left empty in the repository for security.

## Alternative: Use Environment Variable

Instead of putting token in mcp.json, set it in your shell:

```bash
# Add to ~/.zshrc or ~/.bash_profile
export GITHUB_PERSONAL_ACCESS_TOKEN="your_token_here"
```

Then in mcp.json, it will automatically read from environment.

## Usage

After setup, you can use GitHub MCP in Copilot Chat:

```
จง follow style guideline ตาม repo kbtg-ai-workshop-oct Docs
```

This will:
1. Read style guidelines from the workshop repository
2. Apply those guidelines to your code
3. Suggest improvements based on the workshop patterns
