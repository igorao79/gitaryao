import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  images: {
    unoptimized: true,
  },
  // API URL will be set via NEXT_PUBLIC_API_URL env var on Vercel
  // Default: http://localhost:8080 for local dev
};

export default nextConfig;
