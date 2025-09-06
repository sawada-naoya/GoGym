/** @type {import('next').NextConfig} */
const nextConfig = {
  // Enable standalone output for Docker deployment
  output: 'standalone',
  
  // Configure images domain
  images: {
    domains: ['localhost'],
  },
}

module.exports = nextConfig