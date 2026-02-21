# MCP Postgres server: install at build time to avoid npm deprecation warning on every run.
# The @modelcontextprotocol/server-postgres package shows "no longer supported" from npm
# but is still the official read-only Postgres MCP server and works correctly.
FROM node:20-alpine

WORKDIR /app

RUN npm install @modelcontextprotocol/server-postgres@0.6.2

ENTRYPOINT ["node", "node_modules/@modelcontextprotocol/server-postgres/dist/index.js"]
