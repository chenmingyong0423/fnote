This is a [Next.js](https://nextjs.org) project bootstrapped with [`create-next-app`](https://nextjs.org/docs/app/api-reference/cli/create-next-app).

## Getting Started

First, run the development server:

```bash
npm run dev
# or
yarn dev
# or
pnpm dev
# or
bun dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

You can start editing the page by modifying `app/page.tsx`. The page auto-updates as you edit the file.

This project uses [`next/font`](https://nextjs.org/docs/app/building-your-application/optimizing/fonts) to automatically optimize and load [Geist](https://vercel.com/font), a new font family for Vercel.

## Learn More

To learn more about Next.js, take a look at the following resources:

- [Next.js Documentation](https://nextjs.org/docs) - learn about Next.js features and API.
- [Learn Next.js](https://nextjs.org/learn) - an interactive Next.js tutorial.

You can check out [the Next.js GitHub repository](https://github.com/vercel/next.js) - your feedback and contributions are welcome!

## Deploy on Vercel

The easiest way to deploy your Next.js app is to use the [Vercel Platform](https://vercel.com/new?utm_medium=default-template&filter=next.js&utm_source=create-next-app&utm_campaign=create-next-app-readme) from the creators of Next.js.

Check out our [Next.js deployment documentation](https://nextjs.org/docs/app/building-your-application/deploying) for more details.

## Deployment (Docker)

### Environment files
- `.env.local`: local dev (already present)
- `.env.production`: production values for running Next.js (example domains)
- `.env.docker`: defaults for container network (`api` service on 8080)
- `.env.nginx`: vars for nginx reverse proxy (`UPSTREAM_WEB`, `SERVER_NAME`)

### Build & run Next.js container
```bash
# build
DOCKER_BUILDKIT=1 docker build -t nextjs-web:latest -f Dockerfile .

# run with production env file
docker run --env-file .env.production -p 3000:3000 nextjs-web:latest
```

### Build & run nginx reverse proxy
```bash
# build
DOCKER_BUILDKIT=1 docker build -t nextjs-web-nginx:latest -f Dockerfile.nginx .

# run and point to a running Next.js container named "web"
docker run --env-file .env.nginx --link web -p 80:80 nextjs-web-nginx:latest
```

> Tip: prefer docker-compose to wire `web` + `nginx` services if you need both running together.
