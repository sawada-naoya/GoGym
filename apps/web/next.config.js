/** @type {import('next').NextConfig} */
const nextConfig = {
  output: "standalone",

  images: {
    remotePatterns: [
      {
        protocol: "http",
        hostname: "localhost",
      },
      {
        protocol: "https",
        hostname: "**.onrender.com",
      },
      {
        protocol: "https",
        hostname: "**.railway.app",
      },
      {
        protocol: "https",
        hostname: "**.supabase.co",
      },
      {
        protocol: "https",
        hostname: "**.fly.dev",
      },
    ],
  },
};

module.exports = nextConfig;
