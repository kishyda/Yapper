import type { NextConfig } from "next";

const nextConfig = {
    //    experimental: {
    //    turbopack: true, // Enable Turbopack
    //},
    async headers() {
        return [
            {
                source: "/api/:path*",
                headers: [
                    { key: "Access-Control-Allow-Credentials", value: "true" },
                    { key: "Access-Control-Allow-Origin", value: "*" }, // replace this your actual origin
                    { key: "Access-Control-Allow-Methods", value: "GET,DELETE,PATCH,POST,PUT" },
                    { key: "Access-Control-Allow-Headers", value: "X-CSRF-Token, X-Requested-With, Accept, Accept-Version, Content-Length, Content-MD5, Content-Type, Date, X-Api-Version" },
                ]
            }
        ]
    },
    //// @ts-ignore
    //webpack(config, { isServer }) {
    //    if (!isServer) {
    //        // Set the devtoolModuleFilenameTemplate for client-side builds
    //        config.devtool = 'source-map'; // Ensure you're using source maps
    //        config.module.rules.push({
    //          test: /\.css$/,
    //          use: ['style-loader', 'css-loader'],
    //        });
    //
    //        // Customize the source map output for modules
    //        config.output.devtoolModuleFilenameTemplate = 'http://yourdomain.com/[namespace]/[resource-path]';
    //    }
    //
    //    return config;
    //},
}

module.exports = nextConfig
