FROM node:22-alpine AS frontend_builder
WORKDIR /frontend
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable

COPY package.json .
COPY pnpm-lock.yaml .
RUN pnpm install
COPY . .
RUN pnpm run build

FROM node:22-alpine
WORKDIR /frontend
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable

COPY --from=frontend_builder /frontend/package.json .
COPY --from=frontend_builder /frontend/pnpm-lock.yaml .
COPY --from=frontend_builder /frontend/next.config.ts .
COPY --from=frontend_builder /frontend/.next/standalone .
COPY --from=frontend_builder /frontend/.next/static ./.next/static
COPY --from=frontend_builder /frontend/start.sh .
COPY --from=frontend_builder /frontend/wait-for.sh .

EXPOSE 80
